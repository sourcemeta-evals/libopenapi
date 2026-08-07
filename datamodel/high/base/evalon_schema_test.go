// Copyright 2022 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package base

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "go.yaml.in/yaml/v4"
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

// TestEvalonGolden_NewSchema_JSONYAMLKeywordSerialization pins the exact
// JSON and YAML top-level keys emitted for the five new/adjacent high-level
// Schema fields ($comment, contentSchema, $vocabulary, contentEncoding,
// contentMediaType) so that a struct-tag mistake affecting only the
// serialized key names (not projected Go field values) is caught even when
// projection-only tests would otherwise pass. The rendered output is
// unmarshaled into a root map so that presence assertions verify TOP-LEVEL
// placement, not merely substring occurrence anywhere in nested content.
func TestEvalonGolden_NewSchema_JSONYAMLKeywordSerialization(t *testing.T) {
	yml := `type: object
$comment: an example comment
contentEncoding: base64
contentMediaType: application/json
contentSchema:
  type: array
$vocabulary:
  https://example.com/vocab/core: true`

	highSchema := getHighSchema(t, yml)

	yamlBytes, err := highSchema.Render()
	require.NoError(t, err, "Render must succeed for a populated high-level Schema")

	var yamlRoot map[string]any
	require.NoError(t, yaml.Unmarshal(yamlBytes, &yamlRoot),
		"rendered YAML must be parseable back into a root map")

	// YAML output must carry all five case-exact keys at the TOP level.
	requiredKeys := []string{"$comment", "contentEncoding", "contentMediaType", "contentSchema", "$vocabulary"}
	for _, key := range requiredKeys {
		_, present := yamlRoot[key]
		assert.True(t, present,
			"YAML root map must contain the top-level key %q (rendered output: %s)", key, string(yamlBytes))
	}

	// YAML output must NOT carry any of the common misnamed alias forms at the root.
	forbiddenKeys := []string{"comment", "vocabulary", "ContentSchema", "ContentEncoding", "ContentMediaType"}
	for _, key := range forbiddenKeys {
		_, present := yamlRoot[key]
		assert.False(t, present,
			"YAML root map must not contain the alias key %q", key)
	}

	jsonBytes, err := highSchema.MarshalJSON()
	require.NoError(t, err, "MarshalJSON must succeed for a populated high-level Schema")

	var jsonRoot map[string]any
	require.NoError(t, json.Unmarshal(jsonBytes, &jsonRoot),
		"rendered JSON must be parseable back into a root map")

	// JSON output must carry all five case-exact keys at the TOP level.
	for _, key := range requiredKeys {
		_, present := jsonRoot[key]
		assert.True(t, present,
			"JSON root map must contain the top-level key %q (rendered output: %s)", key, string(jsonBytes))
	}

	// JSON output must not carry the common misnamed alias forms at the root.
	for _, key := range forbiddenKeys {
		_, present := jsonRoot[key]
		assert.False(t, present,
			"JSON root map must not contain the alias key %q", key)
	}
}
