package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasMultiKeyTags(t *testing.T) {
	assert.False(t, HasMultiKeyTags([]string{"xxxxx"}))
	assert.False(t, HasMultiKeyTags([]string{"key   `xxxxx yyyy zzzz key:`"}))
	assert.False(t, HasMultiKeyTags([]string{"key   `tagKey:\"tag value\"`"}))
	assert.False(t, HasMultiKeyTags([]string{"key   `  tagKey:\"tag value\"  `"}))
	assert.True(
		t,
		HasMultiKeyTags(
			[]string{
				"xxxx",
				"key   `tagKey1:\"tag value1\"  tagKey2:\"tag value2\" `",
			},
		),
	)
	assert.True(
		t,
		HasMultiKeyTags(
			[]string{
				"key   `  tagKey1:\"tag value1\" tagKey2:\"tag value2\"   tagKey3:\"tag value3\" `",
				"zzzz",
			},
		),
	)
}
