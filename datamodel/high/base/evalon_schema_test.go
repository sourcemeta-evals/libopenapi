// Copyright 2022 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEvalonGolden_NewSchema_JSONSchema202012Keywords verifies that NewSchema
// exposes the three JSON Schema 2020-12 keywords added by this task ($comment,
// contentSchema, $vocabulary) on the high-level Schema when they are present
// on the source YAML. It also pins the accompanying content-family keywords
// (contentEncoding, contentMediaType) so the sibling assignments in NewSchema
// are covered by an executed assertion.
func TestEvalonGolden_NewSchema_JSONSchema202012Keywords(t *testing.T) {
	yml := `type: object
$comment: an example comment
contentEncoding: base64
contentMediaType: application/json
contentSchema:
  type: array
$vocabulary:
  https://example.com/vocab/core: true
  https://example.com/vocab/optional: false`

	highSchema := getHighSchema(t, yml)

	assert.Equal(t, "an example comment", highSchema.Comment)

	assert.Equal(t, "base64", highSchema.ContentEncoding,
		"ContentEncoding must be projected from the low-level Schema onto the high-level view")
	assert.Equal(t, "application/json", highSchema.ContentMediaType,
		"ContentMediaType must be projected from the low-level Schema onto the high-level view")

	assert.NotNil(t, highSchema.ContentSchema, "ContentSchema must be non-nil when present")
	nested := highSchema.ContentSchema.Schema()
	assert.NotNil(t, nested, "ContentSchema.Schema() must resolve the nested schema")
	assert.Contains(t, nested.Type, "array", "nested ContentSchema should carry its type")

	assert.NotNil(t, highSchema.Vocabulary)
	assert.Equal(t, 2, highSchema.Vocabulary.Len())
	core, coreOk := highSchema.Vocabulary.Get("https://example.com/vocab/core")
	assert.True(t, coreOk, "vocab/core entry expected in the high-level Vocabulary map")
	assert.True(t, core, "vocab/core value must be true")
	opt, optOk := highSchema.Vocabulary.Get("https://example.com/vocab/optional")
	assert.True(t, optOk, "vocab/optional entry expected in the high-level Vocabulary map")
	assert.False(t, opt, "vocab/optional value must be false")
}

// TestEvalonGolden_NewSchema_JSONSchema202012Keywords_Empty verifies that the
// three new high-level fields plus the two sibling content-family fields are
// left at their zero values when the source schema omits every content and
// vocabulary keyword. In particular ContentSchema must not be a spurious proxy
// over an absent low-level reference.
func TestEvalonGolden_NewSchema_JSONSchema202012Keywords_Empty(t *testing.T) {
	yml := `type: object`
	highSchema := getHighSchema(t, yml)

	assert.Equal(t, "", highSchema.Comment, "Comment must be the zero string when $comment is absent")
	assert.Equal(t, "", highSchema.ContentEncoding, "ContentEncoding must be the zero string when contentEncoding is absent")
	assert.Equal(t, "", highSchema.ContentMediaType, "ContentMediaType must be the zero string when contentMediaType is absent")
	assert.Nil(t, highSchema.ContentSchema, "ContentSchema must be nil when absent, not a spurious proxy")
	assert.Nil(t, highSchema.Vocabulary, "Vocabulary must be nil when absent")
}
