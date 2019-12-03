package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/dave/dst"
)

var tagKeyRegexp = regexp.MustCompile("([a-zA-Z0-9_-]+):")

// FormatStructTags formats struct tags so that the keys within each block of fields are aligned.
// It's not technically a shortening (and it usually makes these tags longer), so it's being
// kept separate from the core shortening logic for now.
//
// See the struct_tags fixture for examples.
func FormatStructTags(fieldList *dst.FieldList) {
	if fieldList == nil || len(fieldList.List) == 0 {
		return
	}

	blockFields := []*dst.Field{}

	// Divide fields into "blocks" so that we don't do alignments across blank lines
	for f, field := range fieldList.List {
		if f == 0 || field.Decorations().Before == dst.EmptyLine {
			alignTags(blockFields)
			blockFields = blockFields[:0]
		}

		blockFields = append(blockFields, field)
	}

	alignTags(blockFields)
}

// alignTags formats the struct tags within a single field block.
func alignTags(fields []*dst.Field) {
	if len(fields) == 0 {
		return
	}

	maxTagWidths := map[string]int{}
	tagKeys := []string{}
	tagKVs := make([]map[string]string, len(fields))

	// First, scan over all field tags so that we can understand their values and widths
	for f, field := range fields {
		if field.Tag == nil {
			continue
		}

		tagValue := field.Tag.Value

		// The dst library doesn't strip off the backticks, so we need to do this manually
		if tagValue[0] != '`' || tagValue[len(tagValue)-1] != '`' {
			continue
		}

		tagValue = tagValue[1 : len(tagValue)-1]
		structTag := reflect.StructTag(tagValue)

		keyMatches := tagKeyRegexp.FindAllStringSubmatch(tagValue, -1)

		for _, keyMatch := range keyMatches {
			key := keyMatch[1]

			value := structTag.Get(key)

			// Tag is key, value, and some extra chars (two quotes + one colon)
			width := len(key) + len(value) + 3

			if _, ok := maxTagWidths[key]; !ok {
				maxTagWidths[key] = width
				tagKeys = append(tagKeys, key)
			} else if width > maxTagWidths[key] {
				maxTagWidths[key] = width
			}

			if tagKVs[f] == nil {
				tagKVs[f] = map[string]string{}
			}

			tagKVs[f][key] = value
		}
	}

	// Go over all the fields again, replacing each tag with a reformatted one
	for f, field := range fields {
		if tagKVs[f] == nil {
			continue
		}

		tagComponents := []string{}

		for _, key := range tagKeys {
			value, ok := tagKVs[f][key]
			lenUsed := 0

			if ok {
				tagComponents = append(tagComponents, fmt.Sprintf("%s:\"%s\"", key, value))
				lenUsed += len(key) + len(value) + 3
			} else {
				tagComponents = append(tagComponents, "")
			}

			lenRemaining := maxTagWidths[key] - lenUsed

			for i := 0; i < lenRemaining; i++ {
				tagComponents[len(tagComponents)-1] += " "
			}
		}

		updatedTagValue := strings.TrimRight(strings.Join(tagComponents, " "), " ")
		field.Tag.Value = fmt.Sprintf("`%s`", updatedTagValue)
	}
}
