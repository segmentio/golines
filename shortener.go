package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"go/token"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	log "github.com/sirupsen/logrus"
)

var (
	// Strings to look for to identify generated files
	generatedTerms = []string{
		"do not edit",
		"generated by",
		"automatically regenerated",
	}
	// Go directives (should be ignored)
	goDirectiveLine = regexp.MustCompile(`\s*//\s*go:.*`)
)

// The maximum number of shortening "rounds" that we'll allow. The shortening
// process should converge quickly, but we have this here as a safety mechanism to
// prevent loops that prevent termination.
const maxRounds = 20

// ShortenerConfig stores the configuration options exposed by a Shortener instance.
type ShortenerConfig struct {
	MaxLen                   int    // Max target width for each line
	TabLen                   int    // Width of a tab character
	KeepAnnotations          bool   // Whether to keep annotations in final result (for debugging only)
	ShortenComments          bool   // Whether to shorten comments
	ReformatTags             bool   // Whether to reformat struct tags in addition to shortening long lines
	IgnoreGenerated          bool   // Whether to ignore generated files
	DotFile                  string // Path to write dot-formatted output to (for debugging only)
	ChainSplitDots           bool   // Whether to split chain methods by putting dots at ends of lines
	IgnoreBeforeIndentChange bool   // Whether to ignore line length before indent change

	// Formatter that will be run before and after main shortening process. If empty,
	// defaults to goimports (if found), otherwise gofmt.
	BaseFormatterCmd string
}

// Shortener shortens a single go file according to a small set of user style
// preferences.
type Shortener struct {
	config ShortenerConfig

	// Some extra params around the base formatter generated from the BaseFormatterCmd
	// argument in the config.
	baseFormatter     string
	baseFormatterArgs []string
}

// NewShortener creates a new shortener instance from the provided config.
func NewShortener(config ShortenerConfig) *Shortener {
	var formatterComponents []string

	if config.BaseFormatterCmd == "" {
		_, err := exec.LookPath("goimports")
		if err != nil {
			formatterComponents = []string{"gofmt"}
		} else {
			formatterComponents = []string{"goimports"}
		}
	} else {
		formatterComponents = strings.Split(config.BaseFormatterCmd, " ")
	}

	s := &Shortener{
		config:        config,
		baseFormatter: formatterComponents[0],
	}

	if len(formatterComponents) > 1 {
		s.baseFormatterArgs = formatterComponents[1:]
	} else {
		s.baseFormatterArgs = []string{}
	}

	return s
}

// Shorten shortens the provided golang file content bytes.
func (s *Shortener) Shorten(contents []byte) ([]byte, error) {
	if s.config.IgnoreGenerated && s.isGenerated(contents) {
		return contents, nil
	}

	round := 0
	var err error

	// Do initial, non-line-length-aware formatting
	contents, err = s.formatSrc(contents)
	if err != nil {
		return nil, fmt.Errorf("Error formatting source: %+v", err)
	}

	for {
		log.Debugf("Starting round %d", round)

		// Annotate all long lines
		lines := strings.Split(string(contents), "\n")
		annotatedLines, linesToShorten := s.annotateLongLines(lines)
		var stop bool

		if linesToShorten == 0 {
			if round == 0 {
				if !s.config.ReformatTags {
					stop = true
				} else if !HasMultiKeyTags(lines) {
					stop = true
				}
			} else {
				stop = true
			}
		}

		if stop {
			log.Debug("Nothing more to shorten or reformat, stopping")
			break
		}

		contents = []byte(strings.Join(annotatedLines, "\n"))

		// Generate AST
		result, err := decorator.Parse(contents)
		if err != nil {
			return nil, err
		}

		if s.config.DotFile != "" {
			dotFile, err := os.Create(s.config.DotFile)
			if err != nil {
				return nil, err
			}
			defer dotFile.Close()

			log.Debugf("Writing dot file output to %s", s.config.DotFile)
			err = CreateDot(result, dotFile)
			if err != nil {
				return nil, err
			}
		}

		// Shorten the file starting at the top-level declarations
		for _, decl := range result.Decls {
			s.formatNode(decl)
		}

		// Materialize output
		output := bytes.NewBuffer([]byte{})
		err = decorator.Fprint(output, result)
		if err != nil {
			return nil, fmt.Errorf("Error parsing source: %+v", err)
		}
		contents = output.Bytes()

		round++

		if round > maxRounds {
			log.Debugf("Hit max rounds, stopping")
			break
		}
	}

	if !s.config.KeepAnnotations {
		contents = s.removeAnnotations(contents)
	}
	if s.config.ShortenComments {
		contents = s.shortenCommentsFunc(contents)
	}

	// Do final round of non-line-length-aware formatting after we've fixed up the comments
	contents, err = s.formatSrc(contents)
	if err != nil {
		return nil, fmt.Errorf("Error formatting source: %+v", err)
	}

	return contents, nil
}

// formatSrc formats the provided source bytes using the configured "base" formatter (typically
// goimports or gofmt).
func (s *Shortener) formatSrc(contents []byte) ([]byte, error) {
	if s.baseFormatter == "gofmt" {
		return format.Source(contents)
	}

	cmd := exec.Command(s.baseFormatter, s.baseFormatterArgs...)
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	outBuffer := &bytes.Buffer{}
	cmd.Stdout = outBuffer

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	_, err = stdinPipe.Write(contents)
	if err != nil {
		return nil, err
	}
	stdinPipe.Close()

	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	return outBuffer.Bytes(), nil
}

// annotateLongLines adds specially-formatted comments to all eligible lines that are longer than
// the configured target length. If a line already has one of these comments from a previous
// shortening round, then the comment contents are updated.
func (s *Shortener) annotateLongLines(lines []string) ([]string, int) {
	annotatedLines := []string{}
	linesToShorten := 0
	prevLen := -1

	for _, line := range lines {
		length := s.lineLen(line)

		if prevLen > -1 {
			if length <= s.config.MaxLen {
				// Shortening successful, remove previous annotation
				annotatedLines = annotatedLines[:len(annotatedLines)-1]
			} else if length < prevLen {
				// Replace annotation with new length
				annotatedLines[len(annotatedLines)-1] = CreateAnnotation(length)
				linesToShorten++
			}
		} else if !s.isComment(line) && length > s.config.MaxLen {
			annotatedLines = append(
				annotatedLines,
				CreateAnnotation(length),
			)
			linesToShorten++
		}

		annotatedLines = append(annotatedLines, line)
		prevLen = ParseAnnotation(line)
	}

	return annotatedLines, linesToShorten
}

// removeAnnotations removes all comments that were added by the annotateLongLines
// function above.
func (s *Shortener) removeAnnotations(contents []byte) []byte {
	cleanedLines := []string{}
	lines := strings.Split(string(contents), "\n")

	for _, line := range lines {
		if !IsAnnotation(line) {
			cleanedLines = append(cleanedLines, line)
		}
	}

	return []byte(strings.Join(cleanedLines, "\n"))
}

// shortenCommentsFunc attempts to shorten long comments in the provided source. As noted
// in the repo README, this functionality has some quirks and is disabled by default.
func (s *Shortener) shortenCommentsFunc(contents []byte) []byte {
	cleanedLines := []string{}
	words := []string{} // all words in a contiguous sequence of long comments
	prefix := ""
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		if s.isComment(line) && !IsAnnotation(line) &&
			!s.isGoDirective(line) &&
			s.lineLen(line) > s.config.MaxLen {
			start := strings.Index(line, "//")
			prefix = line[0:(start + 2)]
			trimmedLine := strings.Trim(line[(start+2):], " ")
			currLineWords := strings.Split(trimmedLine, " ")
			words = append(words, currLineWords...)
		} else {
			// Reflow the accumulated `words` before appending the unprocessed `line`.
			currLineLen := 0
			currLineWords := []string{}
			maxCommentLen := s.config.MaxLen - s.lineLen(prefix)
			for _, word := range words {
				if currLineLen > 0 && currLineLen+1+len(word) > maxCommentLen {
					cleanedLines = append(
						cleanedLines,
						fmt.Sprintf(
							"%s %s",
							prefix,
							strings.Join(currLineWords, " "),
						),
					)
					currLineWords = []string{}
					currLineLen = 0
				}
				currLineWords = append(currLineWords, word)
				currLineLen += 1 + len(word)
			}
			if currLineLen > 0 {
				cleanedLines = append(
					cleanedLines,
					fmt.Sprintf(
						"%s %s",
						prefix,
						strings.Join(currLineWords, " "),
					),
				)
			}
			words = []string{}

			cleanedLines = append(cleanedLines, line)
		}
	}
	return []byte(strings.Join(cleanedLines, "\n"))
}

// lineLen gets the width of the provided line after tab expansion.
func (s *Shortener) lineLen(line string) int {
	length := 0

	for _, char := range line {
		if char == '\t' {
			length += s.config.TabLen
		} else {
			length++
		}
	}

	return length
}

// isComment determines whether the provided line is a non-block comment.
func (s *Shortener) isComment(line string) bool {
	return strings.HasPrefix(strings.Trim(line, " \t"), "//")
}

// isGoDirective determines whether the provided line is a go directive, e.g. for go generate.
func (s *Shortener) isGoDirective(line string) bool {
	return goDirectiveLine.MatchString(line)
}

// formatNode formats the provided AST node. The appropriate helper function is called
// based on whether the node is a declaration, expression, statement, or spec.
func (s *Shortener) formatNode(node dst.Node) {
	switch n := node.(type) {
	case dst.Decl:
		log.Debugf("Processing declaration: %+v", n)
		s.formatDecl(n)
	case dst.Expr:
		log.Debugf("Processing expression: %+v", n)
		s.formatExpr(n, false, false)
	case dst.Stmt:
		log.Debugf("Processing statement: %+v", n)
		s.formatStmt(n)
	case dst.Spec:
		log.Debugf("Processing spec: %+v", n)
		s.formatSpec(n, false)
	default:
		log.Debugf(
			"Got a node type that can't be shortened: %+v",
			reflect.TypeOf(n),
		)
	}
}

// formatDecl formats an AST declaration node. These include function declarations,
// imports, and constants.
func (s *Shortener) formatDecl(decl dst.Decl) {
	switch d := decl.(type) {
	case *dst.FuncDecl:
		if !s.config.IgnoreBeforeIndentChange && HasAnnotationRecursive(decl) {
			if d.Type != nil && d.Type.Params != nil {
				s.formatFieldList(d.Type.Params)
			}
		}
		s.formatStmt(d.Body)
	case *dst.GenDecl:
		for _, spec := range d.Specs {
			s.formatSpec(spec, HasAnnotation(decl))
		}
	default:
		log.Debugf(
			"Got a declaration type that can't be shortened: %+v",
			reflect.TypeOf(d),
		)
	}
}

// formatFieldList formats a field list in a function declaration.
func (s *Shortener) formatFieldList(fieldList *dst.FieldList) {
	for f, field := range fieldList.List {
		if f == 0 {
			field.Decorations().Before = dst.NewLine
		} else {
			field.Decorations().Before = dst.None
		}

		field.Decorations().After = dst.NewLine
	}
}

// formatStmt formats an AST statement node. Among other examples, these include assignments,
// case clauses, for statements, if statements, and select statements.
func (s *Shortener) formatStmt(stmt dst.Stmt) {
	// Explicitly check for nil statements
	stmtType := reflect.TypeOf(stmt)
	if reflect.ValueOf(stmt) == reflect.Zero(stmtType) {
		return
	}

	shouldShorten := HasAnnotation(stmt)

	switch st := stmt.(type) {
	case *dst.AssignStmt:
		for _, expr := range st.Rhs {
			s.formatExpr(expr, shouldShorten, false)
		}
	case *dst.BlockStmt:
		for _, stmt := range st.List {
			s.formatStmt(stmt)
		}
	case *dst.CaseClause:
		if shouldShorten && !s.config.IgnoreBeforeIndentChange {
			for _, arg := range st.List {
				arg.Decorations().After = dst.NewLine
				s.formatExpr(arg, false, false)
			}
		}

		for _, stmt := range st.Body {
			s.formatStmt(stmt)
		}
	case *dst.CommClause:
		for _, stmt := range st.Body {
			s.formatStmt(stmt)
		}
	case *dst.DeclStmt:
		s.formatDecl(st.Decl)
	case *dst.DeferStmt:
		s.formatExpr(st.Call, shouldShorten, false)
	case *dst.ExprStmt:
		s.formatExpr(st.X, shouldShorten, false)
	case *dst.ForStmt:
		if !s.config.IgnoreBeforeIndentChange {
			s.formatStmt(st.Body)
		}
	case *dst.GoStmt:
		s.formatExpr(st.Call, shouldShorten, false)
	case *dst.IfStmt:
		if !s.config.IgnoreBeforeIndentChange {
			s.formatExpr(st.Cond, shouldShorten, false)
		}
		s.formatStmt(st.Body)
	case *dst.RangeStmt:
		s.formatStmt(st.Body)
	case *dst.ReturnStmt:
		for _, expr := range st.Results {
			s.formatExpr(expr, shouldShorten, false)
		}
	case *dst.SelectStmt:
		s.formatStmt(st.Body)
	case *dst.SwitchStmt:
		s.formatStmt(st.Body)
	default:
		if shouldShorten {
			log.Debugf(
				"Got a statement type that can't be shortened: %+v",
				reflect.TypeOf(st),
			)
		}
	}
}

// formatExpr formats an AST expression node. These include uniary and binary expressions, function
// literals, and key/value pair statements, among others.
func (s *Shortener) formatExpr(expr dst.Expr, force bool, isChain bool) {
	shouldShorten := force || HasAnnotation(expr)

	switch e := expr.(type) {
	case *dst.BinaryExpr:
		if (e.Op == token.LAND || e.Op == token.LOR) && shouldShorten {
			if e.Y.Decorations().Before == dst.NewLine {
				s.formatExpr(e.X, force, isChain)
			} else {
				e.Y.Decorations().Before = dst.NewLine
			}
		} else {
			s.formatExpr(e.X, shouldShorten, isChain)
			s.formatExpr(e.Y, shouldShorten, isChain)
		}
	case *dst.CallExpr:
		_, ok := e.Fun.(*dst.SelectorExpr)

		if ok &&
			s.config.ChainSplitDots &&
			(shouldShorten || HasAnnotationRecursive(e)) &&
			(isChain || s.chainLength(e) > 1) {
			e.Decorations().After = dst.NewLine

			for _, arg := range e.Args {
				s.formatExpr(arg, false, true)
			}

			s.formatExpr(e.Fun, shouldShorten, true)
		} else {
			shortenChildArgs := shouldShorten || HasAnnotationRecursive(e)

			for a, arg := range e.Args {
				if shortenChildArgs {
					if a == 0 {
						arg.Decorations().Before = dst.NewLine
					} else {
						arg.Decorations().After = dst.None
					}
					arg.Decorations().After = dst.NewLine
				}
				s.formatExpr(arg, false, isChain)
			}
			s.formatExpr(e.Fun, shouldShorten, isChain)
		}
	case *dst.CompositeLit:
		if shouldShorten {
			for i, element := range e.Elts {
				if i == 0 {
					element.Decorations().Before = dst.NewLine
				}
				element.Decorations().After = dst.NewLine
			}
		}

		for _, element := range e.Elts {
			s.formatExpr(element, false, isChain)
		}
	case *dst.FuncLit:
		s.formatStmt(e.Body)
	case *dst.FuncType:
		if shouldShorten {
			s.formatFieldList(e.Params)
		}
	case *dst.InterfaceType:
		if s.config.IgnoreBeforeIndentChange {
			return
		}
		for _, method := range e.Methods.List {
			if HasAnnotation(method) {
				s.formatExpr(method.Type, true, isChain)
			}
		}
	case *dst.KeyValueExpr:
		s.formatExpr(e.Value, shouldShorten, isChain)
	case *dst.SelectorExpr:
		s.formatExpr(e.X, shouldShorten, isChain)
	case *dst.StructType:
		if s.config.ReformatTags {
			FormatStructTags(e.Fields)
		}
	case *dst.UnaryExpr:
		s.formatExpr(e.X, shouldShorten, isChain)
	default:
		if shouldShorten {
			log.Debugf(
				"Got an expression type that can't be shortened: %+v",
				reflect.TypeOf(e),
			)
		}
	}
}

// formatSpec formats an AST spec node. These include type specifications, among other things.
func (s *Shortener) formatSpec(spec dst.Spec, force bool) {
	shouldShorten := HasAnnotation(spec) || force
	switch sp := spec.(type) {
	case *dst.ValueSpec:
		for _, expr := range sp.Values {
			s.formatExpr(expr, shouldShorten, false)
		}
	case *dst.TypeSpec:
		s.formatExpr(sp.Type, false, false)
	default:
		if shouldShorten {
			log.Debugf(
				"Got a spec type that can't be shortened: %+v",
				reflect.TypeOf(sp),
			)
		}
	}
}

// isGenerated checks whether the provided file bytes are from a generated file.
// This is done by looking for a set of typically-used strings in the first 5 lines.
func (s *Shortener) isGenerated(contents []byte) bool {
	scanner := bufio.NewScanner(bytes.NewBuffer(contents))

	for i := 0; scanner.Scan(); i++ {
		if i >= 5 {
			return false
		}

		for _, term := range generatedTerms {
			if strings.Contains(strings.ToLower(scanner.Text()), term) {
				return true
			}
		}
	}

	return false
}

// chainLength determines the length of the function call chain in an expression.
func (s *Shortener) chainLength(callExpr *dst.CallExpr) int {
	numCalls := 1
	currCall := callExpr

	for {
		selectorExpr, ok := currCall.Fun.(*dst.SelectorExpr)
		if !ok {
			break
		}
		currCall, ok = selectorExpr.X.(*dst.CallExpr)
		if !ok {
			break
		}
		numCalls++
	}

	return numCalls
}
