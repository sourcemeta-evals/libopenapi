// Copyright 2022 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package base

import (
	"context"
	"testing"

	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/stretchr/testify/assert"
	yaml "go.yaml.in/yaml/v4"
)

func build2020Schema(t *testing.T, yml string) *Schema {
	var node yaml.Node
	assert.NoError(t, yaml.Unmarshal([]byte(yml), &node))
	var sch Schema
	assert.NoError(t, low.BuildModel(node.Content[0], &sch))
	assert.NoError(t, sch.Build(context.Background(), node.Content[0], nil))
	return &sch
}

func TestSchema2020Keywords_CommentAndContentSchemaPresent(t *testing.T) {
	sch := build2020Schema(t, `type: string
$comment: internal usage note
contentMediaType: application/json
contentSchema:
  type: object`)

	assert.Equal(t, "internal usage note", sch.Comment.Value)
	assert.False(t, sch.ContentSchema.IsEmpty())
	assert.NotNil(t, sch.ContentSchema.Value)
	nested := sch.ContentSchema.Value.Schema()
	assert.NotNil(t, nested)
	assert.Equal(t, "object", nested.Type.Value.A)
}

func TestSchema2020Keywords_CommentAbsent(t *testing.T) {
	sch := build2020Schema(t, `type: object
description: no 2020-12 keywords here`)
	assert.True(t, sch.Comment.IsEmpty())
}

func TestSchema2020Keywords_ContentSchemaAbsent(t *testing.T) {
	sch := build2020Schema(t, `type: string
description: no 2020-12 keywords here`)
	assert.True(t, sch.ContentSchema.IsEmpty())
}

func TestSchema2020Keywords_VocabularyPresent(t *testing.T) {
	sch := build2020Schema(t, `type: object
$vocabulary:
  https://example.com/vocab/a: true
  https://example.com/vocab/b: false`)

	assert.NotNil(t, sch.Vocabulary.Value)
	values := map[string]bool{}
	for k, v := range sch.Vocabulary.Value.FromOldest() {
		values[k.Value] = v.Value
	}
	assert.Equal(t, map[string]bool{
		"https://example.com/vocab/a": true,
		"https://example.com/vocab/b": false,
	}, values)
}

func TestSchema2020Keywords_VocabularyAbsent(t *testing.T) {
	sch := build2020Schema(t, `type: object`)
	assert.Nil(t, sch.Vocabulary.Value)
}

func TestSchema2020Keywords_VocabularyScalarIgnored(t *testing.T) {
	sch := build2020Schema(t, `type: object
$vocabulary: true`)
	assert.Nil(t, sch.Vocabulary.Value)
}

func TestSchema2020Keywords_VocabularySequenceIgnored(t *testing.T) {
	sch := build2020Schema(t, `type: object
$vocabulary:
  - https://example.com/vocab/a`)
	assert.Nil(t, sch.Vocabulary.Value)
}

func TestSchema2020Keywords_VocabularyNonBooleanScalarBecomesFalse(t *testing.T) {
	sch := build2020Schema(t, `type: object
$vocabulary:
  https://example.com/vocab/a: enabled`)

	assert.NotNil(t, sch.Vocabulary.Value)
	for k, v := range sch.Vocabulary.Value.FromOldest() {
		assert.Equal(t, "https://example.com/vocab/a", k.Value)
		assert.False(t, v.Value)
	}
}

func TestSchema2020Keywords_HashComment(t *testing.T) {
	first := build2020Schema(t, `type: object
$comment: one`)
	second := build2020Schema(t, `type: object
$comment: two`)
	same := build2020Schema(t, `type: object
$comment: one`)

	assert.NotEqual(t, first.Hash(), second.Hash())
	assert.Equal(t, first.Hash(), same.Hash())
}

func TestSchema2020Keywords_HashContentSchema(t *testing.T) {
	first := build2020Schema(t, `type: string
contentSchema:
  type: object`)
	second := build2020Schema(t, `type: string
contentSchema:
  type: array`)

	assert.NotEqual(t, first.Hash(), second.Hash())
}

func TestSchema2020Keywords_HashVocabulary(t *testing.T) {
	first := build2020Schema(t, `type: object
$vocabulary:
  https://example.com/vocab/a: true`)
	second := build2020Schema(t, `type: object
$vocabulary:
  https://example.com/vocab/a: false`)

	assert.NotEqual(t, first.Hash(), second.Hash())
}
