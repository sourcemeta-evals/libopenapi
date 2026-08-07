// Copyright 2023 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"strings"
	"testing"

	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/datamodel/low/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func compare2020Schemas(t *testing.T, left, right string) *SchemaChanges {
	low.ClearHashCache()
	lDoc, rDoc := test_BuildDoc(wrap2020Schema(left), wrap2020Schema(right))
	lProxy := lDoc.Components.Value.FindSchema("Target").Value
	rProxy := rDoc.Components.Value.FindSchema("Target").Value
	return CompareSchemas(lProxy, rProxy)
}

func wrap2020Schema(schema string) string {
	var indented []string
	for _, line := range strings.Split(schema, "\n") {
		indented = append(indented, "      "+line)
	}
	return `openapi: "3.1.0"
info:
  title: schema comparison fixture
  version: "1.0"
components:
  schemas:
    Target:
` + strings.Join(indented, "\n")
}

func changesForProperty(changes *SchemaChanges, property string) []*Change {
	var matched []*Change
	if changes == nil {
		return matched
	}
	for _, change := range changes.GetAllChanges() {
		if change.Property == property {
			matched = append(matched, change)
		}
	}
	return matched
}

func TestSchema2020Keywords_CommentAdded(t *testing.T) {
	changes := compare2020Schemas(t,
		"type: object",
		"type: object\n$comment: fresh note")

	require.NotNil(t, changes)
	records := changesForProperty(changes, base.CommentLabel)
	require.Len(t, records, 1)
	assert.Equal(t, PropertyAdded, records[0].ChangeType)
	assert.Equal(t, "fresh note", records[0].New)
}

func TestSchema2020Keywords_CommentRemoved(t *testing.T) {
	changes := compare2020Schemas(t,
		"type: object\n$comment: fresh note",
		"type: object")

	require.NotNil(t, changes)
	records := changesForProperty(changes, base.CommentLabel)
	require.Len(t, records, 1)
	assert.Equal(t, PropertyRemoved, records[0].ChangeType)
	assert.Equal(t, "fresh note", records[0].Original)
}

func TestSchema2020Keywords_CommentModified(t *testing.T) {
	changes := compare2020Schemas(t,
		"type: object\n$comment: before",
		"type: object\n$comment: after")

	require.NotNil(t, changes)
	records := changesForProperty(changes, base.CommentLabel)
	require.Len(t, records, 1)
	assert.Equal(t, Modified, records[0].ChangeType)
	assert.Equal(t, "before", records[0].Original)
	assert.Equal(t, "after", records[0].New)
}

func TestSchema2020Keywords_CommentEqualEmitsNothing(t *testing.T) {
	changes := compare2020Schemas(t,
		"type: object\n$comment: identical note",
		"type: object\n$comment: identical note")

	if changes != nil {
		assert.Empty(t, changesForProperty(changes, base.CommentLabel))
		assert.Equal(t, 0, changes.TotalChanges())
	}
}

func TestSchema2020Keywords_VocabularyEntryTransitions(t *testing.T) {
	changes := compare2020Schemas(t,
		`type: object
$vocabulary:
  https://example.com/vocab/kept: true
  https://example.com/vocab/removed: true
  https://example.com/vocab/flipped: true`,
		`type: object
$vocabulary:
  https://example.com/vocab/kept: true
  https://example.com/vocab/added: false
  https://example.com/vocab/flipped: false`)

	require.NotNil(t, changes)
	records := changesForProperty(changes, base.VocabularyLabel)
	require.Len(t, records, 3)

	observed := map[string]int{}
	for _, record := range records {
		switch record.ChangeType {
		case PropertyAdded:
			observed[record.New] = PropertyAdded
		case PropertyRemoved:
			observed[record.Original] = PropertyRemoved
		case Modified:
			assert.Equal(t, "https://example.com/vocab/flipped=true", record.Original)
			assert.Equal(t, "https://example.com/vocab/flipped=false", record.New)
			observed[record.New] = Modified
		}
	}
	assert.Equal(t, map[string]int{
		"https://example.com/vocab/added=false":   PropertyAdded,
		"https://example.com/vocab/removed=true":  PropertyRemoved,
		"https://example.com/vocab/flipped=false": Modified,
	}, observed)

	assert.Len(t, changes.GetAllChanges(), 3)
	assert.Equal(t, 3, changes.TotalChanges())
}

func TestSchema2020Keywords_VocabularyWholeMapAdded(t *testing.T) {
	changes := compare2020Schemas(t,
		"type: object",
		`type: object
$vocabulary:
  https://example.com/vocab/a: true
  https://example.com/vocab/b: false`)

	require.NotNil(t, changes)
	records := changesForProperty(changes, base.VocabularyLabel)
	require.Len(t, records, 2)
	for _, record := range records {
		assert.Equal(t, PropertyAdded, record.ChangeType)
	}
}

func TestSchema2020Keywords_VocabularyWholeMapRemoved(t *testing.T) {
	changes := compare2020Schemas(t,
		`type: object
$vocabulary:
  https://example.com/vocab/a: true
  https://example.com/vocab/b: false`,
		"type: object")

	require.NotNil(t, changes)
	records := changesForProperty(changes, base.VocabularyLabel)
	require.Len(t, records, 2)
	for _, record := range records {
		assert.Equal(t, PropertyRemoved, record.ChangeType)
	}
}

func TestSchema2020Keywords_ContentSchemaAdded(t *testing.T) {
	changes := compare2020Schemas(t,
		"type: string",
		`type: string
contentSchema:
  type: object`)

	require.NotNil(t, changes)
	records := changesForProperty(changes, PropContentSchema)
	require.Len(t, records, 1)
	assert.Contains(t, []int{PropertyAdded, ObjectAdded}, records[0].ChangeType)
}

func TestSchema2020Keywords_ContentSchemaRemoved(t *testing.T) {
	changes := compare2020Schemas(t,
		`type: string
contentSchema:
  type: object`,
		"type: string")

	require.NotNil(t, changes)
	records := changesForProperty(changes, PropContentSchema)
	require.Len(t, records, 1)
	assert.Contains(t, []int{PropertyRemoved, ObjectRemoved}, records[0].ChangeType)
}

func TestSchema2020Keywords_ContentSchemaModifiedAggregates(t *testing.T) {
	changes := compare2020Schemas(t,
		`type: string
contentSchema:
  type: object`,
		`type: string
contentSchema:
  type: array`)

	require.NotNil(t, changes)
	require.NotNil(t, changes.ContentSchemaChanges)
	assert.Equal(t, 1, changes.TotalChanges())

	all := changes.GetAllChanges()
	require.Len(t, all, 1)
	assert.Equal(t, Modified, all[0].ChangeType)
	assert.Equal(t, "object", all[0].Original)
	assert.Equal(t, "array", all[0].New)
}

func TestSchema2020Keywords_ContentSchemaEqualEmitsNothing(t *testing.T) {
	changes := compare2020Schemas(t,
		`type: string
contentSchema:
  type: object`,
		`type: string
contentSchema:
  type: object`)

	if changes != nil {
		assert.Nil(t, changes.ContentSchemaChanges)
		assert.Equal(t, 0, changes.TotalChanges())
	}
}

func TestSchema2020Keywords_MixedBreakingTotals(t *testing.T) {
	changes := &SchemaChanges{
		PropertyChanges: &PropertyChanges{},
		VocabularyChanges: []*Change{
			{
				Property:   base.VocabularyLabel,
				ChangeType: PropertyRemoved,
				Original:   "https://example.com/vocab/a=true",
				Breaking:   true,
			},
			{
				Property:   base.VocabularyLabel,
				ChangeType: PropertyAdded,
				New:        "https://example.com/vocab/b=false",
				Breaking:   false,
			},
		},
		ContentSchemaChanges: &SchemaChanges{
			PropertyChanges: &PropertyChanges{
				Changes: []*Change{
					{
						Property:   "type",
						ChangeType: Modified,
						Original:   "object",
						New:        "array",
						Breaking:   true,
					},
					{
						Property:   "description",
						ChangeType: Modified,
						Original:   "before",
						New:        "after",
						Breaking:   false,
					},
				},
			},
		},
	}

	assert.Equal(t, 4, changes.TotalChanges())
	assert.Equal(t, 2, changes.TotalBreakingChanges())
}
