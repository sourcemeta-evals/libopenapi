// Copyright 2022 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package base

import (
	"context"
	"testing"

	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/stretchr/testify/assert"
	"go.yaml.in/yaml/v4"
)

func TestEvalonGolden_Schema_Comment(t *testing.T) {
	yml := `type: object
$comment: This is a comment that explains the schema purpose
description: A schema with $comment`

	var idxNode yaml.Node
	_ = yaml.Unmarshal([]byte(yml), &idxNode)

	var sch Schema
	err := low.BuildModel(idxNode.Content[0], &sch)
	assert.NoError(t, err)

	err = sch.Build(context.Background(), idxNode.Content[0], nil)
	assert.NoError(t, err)

	assert.Equal(t, "This is a comment that explains the schema purpose", sch.Comment.Value)
}

func TestEvalonGolden_Schema_Comment_Empty(t *testing.T) {
	yml := `type: object
description: A schema without $comment`

	var idxNode yaml.Node
	_ = yaml.Unmarshal([]byte(yml), &idxNode)

	var sch Schema
	err := low.BuildModel(idxNode.Content[0], &sch)
	assert.NoError(t, err)

	err = sch.Build(context.Background(), idxNode.Content[0], nil)
	assert.NoError(t, err)

	assert.True(t, sch.Comment.IsEmpty())
}

func TestEvalonGolden_Schema_ContentSchema(t *testing.T) {
	yml := `type: string
contentMediaType: application/jwt
contentSchema:
  type: object
  properties:
    iss:
      type: string
    exp:
      type: integer`

	var idxNode yaml.Node
	_ = yaml.Unmarshal([]byte(yml), &idxNode)

	var sch Schema
	err := low.BuildModel(idxNode.Content[0], &sch)
	assert.NoError(t, err)

	err = sch.Build(context.Background(), idxNode.Content[0], nil)
	assert.NoError(t, err)

	assert.False(t, sch.ContentSchema.IsEmpty())
	assert.NotNil(t, sch.ContentSchema.Value)

	contentSch := sch.ContentSchema.Value.Schema()
	assert.NotNil(t, contentSch)
	assert.Equal(t, "object", contentSch.Type.Value.A)
}

func TestEvalonGolden_Schema_ContentSchema_Empty(t *testing.T) {
	yml := `type: string
contentMediaType: text/plain`

	var idxNode yaml.Node
	_ = yaml.Unmarshal([]byte(yml), &idxNode)

	var sch Schema
	err := low.BuildModel(idxNode.Content[0], &sch)
	assert.NoError(t, err)

	err = sch.Build(context.Background(), idxNode.Content[0], nil)
	assert.NoError(t, err)

	assert.True(t, sch.ContentSchema.IsEmpty())
}

func TestEvalonGolden_Schema_Vocabulary(t *testing.T) {
	yml := `$vocabulary:
  https://json-schema.org/draft/2020-12/vocab/core: true
  https://json-schema.org/draft/2020-12/vocab/applicator: true
  https://json-schema.org/draft/2020-12/vocab/validation: false
type: object`

	var idxNode yaml.Node
	_ = yaml.Unmarshal([]byte(yml), &idxNode)

	var sch Schema
	err := low.BuildModel(idxNode.Content[0], &sch)
	assert.NoError(t, err)

	err = sch.Build(context.Background(), idxNode.Content[0], nil)
	assert.NoError(t, err)

	assert.NotNil(t, sch.Vocabulary.Value)
	assert.Equal(t, 3, sch.Vocabulary.Value.Len())

	for k, v := range sch.Vocabulary.Value.FromOldest() {
		switch k.Value {
		case "https://json-schema.org/draft/2020-12/vocab/core":
			assert.True(t, v.Value)
		case "https://json-schema.org/draft/2020-12/vocab/applicator":
			assert.True(t, v.Value)
		case "https://json-schema.org/draft/2020-12/vocab/validation":
			assert.False(t, v.Value)
		}
	}
}

func TestEvalonGolden_Schema_Vocabulary_Empty(t *testing.T) {
	yml := `type: object
description: A regular schema without $vocabulary`

	var idxNode yaml.Node
	_ = yaml.Unmarshal([]byte(yml), &idxNode)

	var sch Schema
	err := low.BuildModel(idxNode.Content[0], &sch)
	assert.NoError(t, err)

	err = sch.Build(context.Background(), idxNode.Content[0], nil)
	assert.NoError(t, err)

	assert.Nil(t, sch.Vocabulary.Value)
}

// TestEvalonGolden_Schema_Vocabulary_NonBooleanScalarsDefaultToFalse tests
// that $vocabulary map values which are not actual YAML boolean scalars
// (numbers, unquoted alternative truthy tokens, etc.) silently default to
// false per the task contract, rather than being permissively coerced via
// strconv.ParseBool.
func TestEvalonGolden_Schema_Vocabulary_NonBooleanScalarsDefaultToFalse(t *testing.T) {
	yml := `type: object
$vocabulary:
  "https://example.com/vocab/one": 1
  "https://example.com/vocab/zero": 0
  "https://example.com/vocab/t": t
  "https://example.com/vocab/f": f
  "https://example.com/vocab/TRUE": TRUE
  "https://example.com/vocab/FALSE": FALSE`

	var idxNode yaml.Node
	_ = yaml.Unmarshal([]byte(yml), &idxNode)

	var sch Schema
	err := low.BuildModel(idxNode.Content[0], &sch)
	assert.NoError(t, err)

	err = sch.Build(context.Background(), idxNode.Content[0], nil)
	assert.NoError(t, err)

	assert.NotNil(t, sch.Vocabulary.Value)
	assert.Equal(t, 6, sch.Vocabulary.Value.Len())

	for k, v := range sch.Vocabulary.Value.FromOldest() {
		assert.False(t, v.Value, "%s should default to false because it is not a YAML boolean scalar", k.Value)
	}
}

// TestEvalonGolden_Schema_Vocabulary_InvalidBooleanDefaultsToFalse tests that
// invalid boolean values default to false when parsed with strconv.ParseBool.
func TestEvalonGolden_Schema_Vocabulary_InvalidBooleanDefaultsToFalse(t *testing.T) {
	yml := `type: object
$vocabulary:
  "https://example.com/vocab/invalid": notaboolean
  "https://example.com/vocab/valid": true`

	var idxNode yaml.Node
	_ = yaml.Unmarshal([]byte(yml), &idxNode)

	var sch Schema
	err := low.BuildModel(idxNode.Content[0], &sch)
	assert.NoError(t, err)

	err = sch.Build(context.Background(), idxNode.Content[0], nil)
	assert.NoError(t, err)

	assert.NotNil(t, sch.Vocabulary.Value)
	assert.Equal(t, 2, sch.Vocabulary.Value.Len())

	for k, v := range sch.Vocabulary.Value.FromOldest() {
		switch k.Value {
		case "https://example.com/vocab/invalid":
			assert.False(t, v.Value, "Invalid boolean should default to false")
		case "https://example.com/vocab/valid":
			assert.True(t, v.Value, "true should parse as true")
		}
	}
}

func TestEvalonGolden_Schema_Hash_IncludesNewFields(t *testing.T) {
	yml1 := `type: object
$comment: Comment 1`

	yml2 := `type: object
$comment: Comment 2`

	var node1, node2 yaml.Node
	_ = yaml.Unmarshal([]byte(yml1), &node1)
	_ = yaml.Unmarshal([]byte(yml2), &node2)

	var sch1, sch2 Schema
	_ = low.BuildModel(node1.Content[0], &sch1)
	_ = sch1.Build(context.Background(), node1.Content[0], nil)

	_ = low.BuildModel(node2.Content[0], &sch2)
	_ = sch2.Build(context.Background(), node2.Content[0], nil)

	hash1 := sch1.Hash()
	hash2 := sch2.Hash()

	assert.NotEqual(t, hash1, hash2)
}

// TestEvalonGolden_Schema_Hash_VocabularyDifferent tests that different
// vocabulary values produce different hashes.
func TestEvalonGolden_Schema_Hash_VocabularyDifferent(t *testing.T) {
	yml1 := `type: object
$vocabulary:
  "https://example.com/vocab/core": true`

	yml2 := `type: object
$vocabulary:
  "https://example.com/vocab/core": false`

	var node1, node2 yaml.Node
	_ = yaml.Unmarshal([]byte(yml1), &node1)
	_ = yaml.Unmarshal([]byte(yml2), &node2)

	var sch1, sch2 Schema
	_ = low.BuildModel(node1.Content[0], &sch1)
	_ = sch1.Build(context.Background(), node1.Content[0], nil)

	_ = low.BuildModel(node2.Content[0], &sch2)
	_ = sch2.Build(context.Background(), node2.Content[0], nil)

	hash1 := sch1.Hash()
	hash2 := sch2.Hash()

	assert.NotEqual(t, hash1, hash2)
}

// TestEvalonGolden_Schema_Hash_ContentSchemaDifferent tests that different
// contentSchema produces different hashes.
func TestEvalonGolden_Schema_Hash_ContentSchemaDifferent(t *testing.T) {
	yml1 := `type: string
contentMediaType: application/json
contentSchema:
  type: object`

	yml2 := `type: string
contentMediaType: application/json
contentSchema:
  type: array`

	var node1, node2 yaml.Node
	_ = yaml.Unmarshal([]byte(yml1), &node1)
	_ = yaml.Unmarshal([]byte(yml2), &node2)

	var sch1, sch2 Schema
	_ = low.BuildModel(node1.Content[0], &sch1)
	_ = sch1.Build(context.Background(), node1.Content[0], nil)

	_ = low.BuildModel(node2.Content[0], &sch2)
	_ = sch2.Build(context.Background(), node2.Content[0], nil)

	hash1 := sch1.Hash()
	hash2 := sch2.Hash()

	assert.NotEqual(t, hash1, hash2)
}

// TestEvalonGolden_Schema_Hash_VocabularyOrderIndependent asserts that two
// schemas with the SAME $vocabulary entries declared in DIFFERENT source order
// produce equal hashes. This catches implementations that iterate a native Go
// map or preserve YAML insertion order when hashing the vocabulary, allowing
// otherwise-equivalent schemas to compare unequal.
func TestEvalonGolden_Schema_Hash_VocabularyOrderIndependent(t *testing.T) {
	yml1 := `type: object
$vocabulary:
  "https://example.com/vocab/alpha": true
  "https://example.com/vocab/beta": false
  "https://example.com/vocab/gamma": true`

	yml2 := `type: object
$vocabulary:
  "https://example.com/vocab/gamma": true
  "https://example.com/vocab/alpha": true
  "https://example.com/vocab/beta": false`

	var node1, node2 yaml.Node
	_ = yaml.Unmarshal([]byte(yml1), &node1)
	_ = yaml.Unmarshal([]byte(yml2), &node2)

	var sch1, sch2 Schema
	_ = low.BuildModel(node1.Content[0], &sch1)
	_ = sch1.Build(context.Background(), node1.Content[0], nil)

	_ = low.BuildModel(node2.Content[0], &sch2)
	_ = sch2.Build(context.Background(), node2.Content[0], nil)

	hash1 := sch1.Hash()
	hash2 := sch2.Hash()

	assert.Equal(t, hash1, hash2,
		"equivalent $vocabulary entries in different source order must hash equally")
}

// TestEvalonGolden_Schema_Vocabulary_StrTaggedTrueDefaultsToFalse asserts that
// an explicit YAML !!str-tagged "true" scalar in a $vocabulary map silently
// parses as `false` per the strict-YAML boolean contract, while a plain
// canonical unquoted `true` in the same map still parses as `true`.
// Implementations that inspect scalar kind, style, or text without checking
// the YAML tag (e.g. via utils.IsNodeBoolValue which requires the !!bool tag)
// incorrectly treat !!str true as boolean true and fail this test.
func TestEvalonGolden_Schema_Vocabulary_StrTaggedTrueDefaultsToFalse(t *testing.T) {
	yml := `type: object
$vocabulary:
  "https://example.com/vocab/canonical": true
  "https://example.com/vocab/str-tagged": !!str true`

	var idxNode yaml.Node
	_ = yaml.Unmarshal([]byte(yml), &idxNode)

	var sch Schema
	err := low.BuildModel(idxNode.Content[0], &sch)
	assert.NoError(t, err, "BuildModel must succeed on the $vocabulary fixture")

	err = sch.Build(context.Background(), idxNode.Content[0], nil)
	assert.NoError(t, err, "Schema.Build must succeed on the $vocabulary fixture and not error on a non-boolean scalar")

	assert.NotNil(t, sch.Vocabulary.Value, "Vocabulary must be populated")
	assert.Equal(t, 2, sch.Vocabulary.Value.Len())

	for k, v := range sch.Vocabulary.Value.FromOldest() {
		switch k.Value {
		case "https://example.com/vocab/canonical":
			assert.True(t, v.Value,
				"a plain canonical unquoted true entry must remain true")
		case "https://example.com/vocab/str-tagged":
			assert.False(t, v.Value,
				"an explicit !!str-tagged 'true' scalar must silently default to false; implementations that miss the tag check and accept the string as boolean fail here")
		}
	}
}
