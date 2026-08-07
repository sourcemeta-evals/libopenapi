// Copyright 2023 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "go.yaml.in/yaml/v4"
)

func TestSchema2020Keywords_Projection(t *testing.T) {
	highSchema := getHighSchema(t, `type: object
$comment: projected comment
contentSchema:
  type: array
$vocabulary:
  https://example.com/vocab/core: true
  https://example.com/vocab/extra: false`)

	assert.Equal(t, "projected comment", highSchema.Comment)

	require.NotNil(t, highSchema.ContentSchema)
	nested := highSchema.ContentSchema.Schema()
	require.NotNil(t, nested)
	assert.Equal(t, []string{"array"}, nested.Type)

	require.NotNil(t, highSchema.Vocabulary)
	values := map[string]bool{}
	for uri, enabled := range highSchema.Vocabulary.FromOldest() {
		values[uri] = enabled
	}
	assert.Equal(t, map[string]bool{
		"https://example.com/vocab/core":  true,
		"https://example.com/vocab/extra": false,
	}, values)
}

func TestSchema2020Keywords_RenderedTopLevelKeys(t *testing.T) {
	highSchema := getHighSchema(t, `type: object
$comment: rendered comment
contentSchema:
  type: array
$vocabulary:
  https://example.com/vocab/core: true`)

	rendered, err := highSchema.Render()
	require.NoError(t, err)

	var root map[string]any
	require.NoError(t, yaml.Unmarshal(rendered, &root))

	for _, key := range []string{"$comment", "contentSchema", "$vocabulary"} {
		_, present := root[key]
		assert.True(t, present, "rendered output must contain the top-level key %q", key)
	}
}
