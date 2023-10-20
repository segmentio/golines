package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrettyDiff(t *testing.T) {
	// For now, this just tests that the script runs without error
	err := PrettyDiff(
		"test_path.txt",
		[]byte("line 1\nline 2"),
		[]byte("line 1\nline 2 modified"),
	)
	assert.Nil(t, err)
}
