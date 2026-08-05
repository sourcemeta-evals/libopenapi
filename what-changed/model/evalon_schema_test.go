// Copyright 2022 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/datamodel/low/base"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/stretchr/testify/assert"
	yaml "go.yaml.in/yaml/v4"
)

// TestCompareSchemas_Comment_Added tests $comment addition detection
func TestEvalonGolden_CompareSchemas_Comment_Added(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `openapi: "3.1.0"
info:
  title: left
  version: "1.0"
components:
  schemas:
    Pet:
      type: object`

	right := `openapi: "3.1.0"
info:
  title: right
  version: "1.0"
components:
  schemas:
    Pet:
      $comment: "This is a comment"
      type: object`

	leftDoc, rightDoc := test_BuildDoc(left, right)

	lSchemaProxy := leftDoc.Components.Value.FindSchema("Pet").Value
	rSchemaProxy := rightDoc.Components.Value.FindSchema("Pet").Value

	changes := CompareSchemas(lSchemaProxy, rSchemaProxy)
	assert.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())

	found := false
	for _, change := range changes.Changes {
		if change.Property == PropComment {
			found = true
			assert.Equal(t, PropertyAdded, change.ChangeType)
			assert.Equal(t, "This is a comment", change.New)
			assert.False(t, change.Breaking)
			break
		}
	}
	assert.True(t, found, "Should find $comment property change")
}

// TestCompareSchemas_Comment_Removed tests $comment removal detection
func TestEvalonGolden_CompareSchemas_Comment_Removed(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `openapi: "3.1.0"
info:
  title: left
  version: "1.0"
components:
  schemas:
    Pet:
      $comment: "This is a comment"
      type: object`

	right := `openapi: "3.1.0"
info:
  title: right
  version: "1.0"
components:
  schemas:
    Pet:
      type: object`

	leftDoc, rightDoc := test_BuildDoc(left, right)

	lSchemaProxy := leftDoc.Components.Value.FindSchema("Pet").Value
	rSchemaProxy := rightDoc.Components.Value.FindSchema("Pet").Value

	changes := CompareSchemas(lSchemaProxy, rSchemaProxy)
	assert.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())

	found := false
	for _, change := range changes.Changes {
		if change.Property == PropComment {
			found = true
			assert.Equal(t, PropertyRemoved, change.ChangeType)
			assert.Equal(t, "This is a comment", change.Original)
			assert.False(t, change.Breaking)
			break
		}
	}
	assert.True(t, found, "Should find $comment property change")
}

// TestCompareSchemas_Comment_Modified tests $comment modification detection
func TestEvalonGolden_CompareSchemas_Comment_Modified(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `openapi: "3.1.0"
info:
  title: left
  version: "1.0"
components:
  schemas:
    Pet:
      $comment: "Original comment"
      type: object`

	right := `openapi: "3.1.0"
info:
  title: right
  version: "1.0"
components:
  schemas:
    Pet:
      $comment: "Modified comment"
      type: object`

	leftDoc, rightDoc := test_BuildDoc(left, right)

	lSchemaProxy := leftDoc.Components.Value.FindSchema("Pet").Value
	rSchemaProxy := rightDoc.Components.Value.FindSchema("Pet").Value

	changes := CompareSchemas(lSchemaProxy, rSchemaProxy)
	assert.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())

	found := false
	for _, change := range changes.Changes {
		if change.Property == PropComment {
			found = true
			assert.Equal(t, Modified, change.ChangeType)
			assert.Equal(t, "Original comment", change.Original)
			assert.Equal(t, "Modified comment", change.New)
			assert.False(t, change.Breaking)
			break
		}
	}
	assert.True(t, found, "Should find $comment property change")
}

// TestCompareSchemas_Comment_NoChange tests identical $comment produces no changes
func TestEvalonGolden_CompareSchemas_Comment_NoChange(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `openapi: "3.1.0"
info:
  title: left
  version: "1.0"
components:
  schemas:
    Pet:
      $comment: "Same comment"
      type: object`

	right := `openapi: "3.1.0"
info:
  title: right
  version: "1.0"
components:
  schemas:
    Pet:
      $comment: "Same comment"
      type: object`

	leftDoc, rightDoc := test_BuildDoc(left, right)

	lSchemaProxy := leftDoc.Components.Value.FindSchema("Pet").Value
	rSchemaProxy := rightDoc.Components.Value.FindSchema("Pet").Value

	changes := CompareSchemas(lSchemaProxy, rSchemaProxy)
	assert.Nil(t, changes)
}

// TestCompareSchemas_ContentSchema_Added tests contentSchema addition detection
func TestEvalonGolden_CompareSchemas_ContentSchema_Added(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `openapi: "3.1.0"
info:
  title: left
  version: "1.0"
components:
  schemas:
    Pet:
      type: string
      contentMediaType: application/json`

	right := `openapi: "3.1.0"
info:
  title: right
  version: "1.0"
components:
  schemas:
    Pet:
      type: string
      contentMediaType: application/json
      contentSchema:
        type: object`

	leftDoc, rightDoc := test_BuildDoc(left, right)

	lSchemaProxy := leftDoc.Components.Value.FindSchema("Pet").Value
	rSchemaProxy := rightDoc.Components.Value.FindSchema("Pet").Value

	changes := CompareSchemas(lSchemaProxy, rSchemaProxy)
	assert.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())

	// Search the aggregate visibility slice so the test does not require any
	// particular storage layout for the contentSchema addition record: it may
	// live in the top-level Changes slice or in a dedicated ContentSchemaChanges
	// container, as long as GetAllChanges surfaces it.
	found := false
	for _, change := range changes.GetAllChanges() {
		if change.Property == PropContentSchema {
			found = true
			assert.Contains(t, []int{PropertyAdded, ObjectAdded}, change.ChangeType,
				"contentSchema addition may be represented as either PropertyAdded or ObjectAdded")
			assert.True(t, change.Breaking)
			break
		}
	}
	assert.True(t, found, "Should find contentSchema property change via GetAllChanges()")
}

// TestCompareSchemas_ContentSchema_Removed tests contentSchema removal detection
func TestEvalonGolden_CompareSchemas_ContentSchema_Removed(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `openapi: "3.1.0"
info:
  title: left
  version: "1.0"
components:
  schemas:
    Pet:
      type: string
      contentMediaType: application/json
      contentSchema:
        type: object`

	right := `openapi: "3.1.0"
info:
  title: right
  version: "1.0"
components:
  schemas:
    Pet:
      type: string
      contentMediaType: application/json`

	leftDoc, rightDoc := test_BuildDoc(left, right)

	lSchemaProxy := leftDoc.Components.Value.FindSchema("Pet").Value
	rSchemaProxy := rightDoc.Components.Value.FindSchema("Pet").Value

	changes := CompareSchemas(lSchemaProxy, rSchemaProxy)
	assert.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())

	// Search the aggregate visibility slice so the test does not require any
	// particular storage layout for the contentSchema removal record: it may
	// live in the top-level Changes slice or in a dedicated ContentSchemaChanges
	// container, as long as GetAllChanges surfaces it.
	found := false
	for _, change := range changes.GetAllChanges() {
		if change.Property == PropContentSchema {
			found = true
			assert.Contains(t, []int{PropertyRemoved, ObjectRemoved}, change.ChangeType,
				"contentSchema removal may be represented as either PropertyRemoved or ObjectRemoved")
			assert.True(t, change.Breaking)
			break
		}
	}
	assert.True(t, found, "Should find contentSchema property change via GetAllChanges()")
}

// TestCompareSchemas_ContentSchema_Modified tests contentSchema modification detection
func TestEvalonGolden_CompareSchemas_ContentSchema_Modified(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `openapi: "3.1.0"
info:
  title: left
  version: "1.0"
components:
  schemas:
    Pet:
      type: string
      contentMediaType: application/json
      contentSchema:
        type: object`

	right := `openapi: "3.1.0"
info:
  title: right
  version: "1.0"
components:
  schemas:
    Pet:
      type: string
      contentMediaType: application/json
      contentSchema:
        type: array`

	leftDoc, rightDoc := test_BuildDoc(left, right)

	lSchemaProxy := leftDoc.Components.Value.FindSchema("Pet").Value
	rSchemaProxy := rightDoc.Components.Value.FindSchema("Pet").Value

	changes := CompareSchemas(lSchemaProxy, rSchemaProxy)
	assert.NotNil(t, changes)
	assert.NotNil(t, changes.ContentSchemaChanges)
	assert.Equal(t, 1, changes.ContentSchemaChanges.TotalChanges())
}

// TestCompareSchemas_Vocabulary_Added tests $vocabulary entry addition detection
func TestEvalonGolden_CompareSchemas_Vocabulary_Added(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `openapi: "3.1.0"
info:
  title: left
  version: "1.0"
components:
  schemas:
    Pet:
      $vocabulary:
        "https://json-schema.org/draft/2020-12/vocab/core": true
      type: object`

	right := `openapi: "3.1.0"
info:
  title: right
  version: "1.0"
components:
  schemas:
    Pet:
      $vocabulary:
        "https://json-schema.org/draft/2020-12/vocab/core": true
        "https://json-schema.org/draft/2020-12/vocab/validation": true
      type: object`

	leftDoc, rightDoc := test_BuildDoc(left, right)

	lSchemaProxy := leftDoc.Components.Value.FindSchema("Pet").Value
	rSchemaProxy := rightDoc.Components.Value.FindSchema("Pet").Value

	changes := CompareSchemas(lSchemaProxy, rSchemaProxy)
	assert.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Len(t, changes.VocabularyChanges, 1)
	assert.Equal(t, PropertyAdded, changes.VocabularyChanges[0].ChangeType)
	assert.True(t, changes.VocabularyChanges[0].Breaking)
}

// TestCompareSchemas_Vocabulary_Removed tests $vocabulary entry removal detection
func TestEvalonGolden_CompareSchemas_Vocabulary_Removed(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `openapi: "3.1.0"
info:
  title: left
  version: "1.0"
components:
  schemas:
    Pet:
      $vocabulary:
        "https://json-schema.org/draft/2020-12/vocab/core": true
        "https://json-schema.org/draft/2020-12/vocab/validation": true
      type: object`

	right := `openapi: "3.1.0"
info:
  title: right
  version: "1.0"
components:
  schemas:
    Pet:
      $vocabulary:
        "https://json-schema.org/draft/2020-12/vocab/core": true
      type: object`

	leftDoc, rightDoc := test_BuildDoc(left, right)

	lSchemaProxy := leftDoc.Components.Value.FindSchema("Pet").Value
	rSchemaProxy := rightDoc.Components.Value.FindSchema("Pet").Value

	changes := CompareSchemas(lSchemaProxy, rSchemaProxy)
	assert.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Len(t, changes.VocabularyChanges, 1)
	assert.Equal(t, PropertyRemoved, changes.VocabularyChanges[0].ChangeType)
	assert.True(t, changes.VocabularyChanges[0].Breaking)
}

// TestCompareSchemas_Vocabulary_Modified tests $vocabulary value modification detection
func TestEvalonGolden_CompareSchemas_Vocabulary_Modified(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `openapi: "3.1.0"
info:
  title: left
  version: "1.0"
components:
  schemas:
    Pet:
      $vocabulary:
        "https://json-schema.org/draft/2020-12/vocab/core": true
      type: object`

	right := `openapi: "3.1.0"
info:
  title: right
  version: "1.0"
components:
  schemas:
    Pet:
      $vocabulary:
        "https://json-schema.org/draft/2020-12/vocab/core": false
      type: object`

	leftDoc, rightDoc := test_BuildDoc(left, right)

	lSchemaProxy := leftDoc.Components.Value.FindSchema("Pet").Value
	rSchemaProxy := rightDoc.Components.Value.FindSchema("Pet").Value

	changes := CompareSchemas(lSchemaProxy, rSchemaProxy)
	assert.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Len(t, changes.VocabularyChanges, 1)
	assert.Equal(t, Modified, changes.VocabularyChanges[0].ChangeType)
	assert.True(t, changes.VocabularyChanges[0].Breaking)
}

// TestCompareSchemas_Vocabulary_NoChange tests identical $vocabulary produces no changes
func TestEvalonGolden_CompareSchemas_Vocabulary_NoChange(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `openapi: "3.1.0"
info:
  title: left
  version: "1.0"
components:
  schemas:
    Pet:
      $vocabulary:
        "https://json-schema.org/draft/2020-12/vocab/core": true
      type: object`

	right := `openapi: "3.1.0"
info:
  title: right
  version: "1.0"
components:
  schemas:
    Pet:
      $vocabulary:
        "https://json-schema.org/draft/2020-12/vocab/core": true
      type: object`

	leftDoc, rightDoc := test_BuildDoc(left, right)

	lSchemaProxy := leftDoc.Components.Value.FindSchema("Pet").Value
	rSchemaProxy := rightDoc.Components.Value.FindSchema("Pet").Value

	changes := CompareSchemas(lSchemaProxy, rSchemaProxy)
	assert.Nil(t, changes)
}

// TestCompareSchemas_Vocabulary_AddedFromNil tests $vocabulary added where none existed
func TestEvalonGolden_CompareSchemas_Vocabulary_AddedFromNil(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `openapi: "3.1.0"
info:
  title: left
  version: "1.0"
components:
  schemas:
    Pet:
      type: object`

	right := `openapi: "3.1.0"
info:
  title: right
  version: "1.0"
components:
  schemas:
    Pet:
      $vocabulary:
        "https://json-schema.org/draft/2020-12/vocab/core": true
      type: object`

	leftDoc, rightDoc := test_BuildDoc(left, right)

	lSchemaProxy := leftDoc.Components.Value.FindSchema("Pet").Value
	rSchemaProxy := rightDoc.Components.Value.FindSchema("Pet").Value

	changes := CompareSchemas(lSchemaProxy, rSchemaProxy)
	assert.NotNil(t, changes)
	assert.Len(t, changes.VocabularyChanges, 1)
	assert.Equal(t, PropertyAdded, changes.VocabularyChanges[0].ChangeType)
}

// TestCompareSchemas_Vocabulary_RemovedToNil tests $vocabulary removed to nil
func TestEvalonGolden_CompareSchemas_Vocabulary_RemovedToNil(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `openapi: "3.1.0"
info:
  title: left
  version: "1.0"
components:
  schemas:
    Pet:
      $vocabulary:
        "https://json-schema.org/draft/2020-12/vocab/core": true
      type: object`

	right := `openapi: "3.1.0"
info:
  title: right
  version: "1.0"
components:
  schemas:
    Pet:
      type: object`

	leftDoc, rightDoc := test_BuildDoc(left, right)

	lSchemaProxy := leftDoc.Components.Value.FindSchema("Pet").Value
	rSchemaProxy := rightDoc.Components.Value.FindSchema("Pet").Value

	changes := CompareSchemas(lSchemaProxy, rSchemaProxy)
	assert.NotNil(t, changes)
	assert.Len(t, changes.VocabularyChanges, 1)
	assert.Equal(t, PropertyRemoved, changes.VocabularyChanges[0].ChangeType)
}

// TestCheckVocabularyChanges_BothNil tests the checkVocabularyChanges helper with both nil
func TestEvalonGolden_CheckVocabularyChanges_BothNil(t *testing.T) {
	changes := checkVocabularyChanges(nil, nil)
	assert.Nil(t, changes)
}

// TestCompareSchemas_Vocabulary_MultipleChanges tests multiple vocabulary changes at once
func TestEvalonGolden_CompareSchemas_Vocabulary_MultipleChanges(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `openapi: "3.1.0"
info:
  title: left
  version: "1.0"
components:
  schemas:
    Pet:
      $vocabulary:
        "https://json-schema.org/draft/2020-12/vocab/core": true
        "https://json-schema.org/draft/2020-12/vocab/validation": true
      type: object`

	right := `openapi: "3.1.0"
info:
  title: right
  version: "1.0"
components:
  schemas:
    Pet:
      $vocabulary:
        "https://json-schema.org/draft/2020-12/vocab/core": false
        "https://json-schema.org/draft/2020-12/vocab/applicator": true
      type: object`

	leftDoc, rightDoc := test_BuildDoc(left, right)

	lSchemaProxy := leftDoc.Components.Value.FindSchema("Pet").Value
	rSchemaProxy := rightDoc.Components.Value.FindSchema("Pet").Value

	changes := CompareSchemas(lSchemaProxy, rSchemaProxy)
	assert.NotNil(t, changes)
	// Should have 3 changes: core modified, validation removed, applicator added
	assert.Equal(t, 3, changes.TotalChanges())
	assert.Len(t, changes.VocabularyChanges, 3)
}

// TestGetAllChanges_IncludesContentSchemaAndVocabulary verifies that
// SchemaChanges.GetAllChanges() surfaces both recursive contentSchema
// changes and every $vocabulary entry change alongside the other change
// containers, per the task contract's aggregate-accessor requirement.
func TestEvalonGolden_GetAllChanges_IncludesContentSchemaAndVocabulary(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `openapi: "3.1.0"
info:
  title: left
  version: "1.0"
components:
  schemas:
    Pet:
      type: string
      contentMediaType: application/json
      contentSchema:
        type: object
      $vocabulary:
        https://example.com/vocab/keep: true
        https://example.com/vocab/toggle: true
        https://example.com/vocab/removed: false`

	right := `openapi: "3.1.0"
info:
  title: right
  version: "1.0"
components:
  schemas:
    Pet:
      type: string
      contentMediaType: application/json
      contentSchema:
        type: array
      $vocabulary:
        https://example.com/vocab/keep: true
        https://example.com/vocab/toggle: false
        https://example.com/vocab/added: true`

	leftDoc, rightDoc := test_BuildDoc(left, right)
	lSchemaProxy := leftDoc.Components.Value.FindSchema("Pet").Value
	rSchemaProxy := rightDoc.Components.Value.FindSchema("Pet").Value

	changes := CompareSchemas(lSchemaProxy, rSchemaProxy)
	assert.NotNil(t, changes)

	// Container-level sanity: the new containers picked up their changes.
	assert.NotNil(t, changes.ContentSchemaChanges,
		"contentSchema change must be captured in ContentSchemaChanges")
	assert.GreaterOrEqual(t, changes.ContentSchemaChanges.TotalChanges(), 1,
		"nested contentSchema type change must contribute to ContentSchemaChanges")
	assert.Len(t, changes.VocabularyChanges, 3,
		"three $vocabulary entry changes expected: one added, one removed, one modified")

	// Aggregate accessor must surface every change from both new containers.
	all := changes.GetAllChanges()
	for _, c := range changes.ContentSchemaChanges.GetAllChanges() {
		assert.Contains(t, all, c,
			"GetAllChanges must include recursive contentSchema changes")
	}
	for _, c := range changes.VocabularyChanges {
		assert.Contains(t, all, c,
			"GetAllChanges must include every $vocabulary entry change")
	}
}

// TestEvalonGolden_CheckVocabularyChanges_DeterministicOrdering verifies that
// the emitted $vocabulary change records are independent of the source
// insertion order of the input maps. The test constructs two semantically
// equivalent left maps with different insertion orders (and likewise two
// equivalent right maps), invokes the comparator once per pair, and asserts
// that both pairs produce the same sequence of records in canonical URI order.
// This catches implementations that iterate a native Go map or preserve
// insertion order instead of applying a stable sort to the URI keys.
func TestEvalonGolden_CheckVocabularyChanges_DeterministicOrdering(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	buildVocab := func(entries [][2]string) *orderedmap.Map[low.KeyReference[string], low.ValueReference[bool]] {
		m := orderedmap.New[low.KeyReference[string], low.ValueReference[bool]]()
		for _, entry := range entries {
			m.Set(
				low.KeyReference[string]{Value: entry[0]},
				low.ValueReference[bool]{Value: entry[1] == "true"},
			)
		}
		return m
	}

	// leftEntriesA and leftEntriesB carry the SAME URI/value set but in
	// DIFFERENT insertion orders. Same for rightEntriesA and rightEntriesB.
	leftEntriesA := [][2]string{
		{"https://example.com/vocab/c", "true"},
		{"https://example.com/vocab/shared", "true"},
		{"https://example.com/vocab/a", "false"},
	}
	leftEntriesB := [][2]string{
		{"https://example.com/vocab/a", "false"},
		{"https://example.com/vocab/c", "true"},
		{"https://example.com/vocab/shared", "true"},
	}
	rightEntriesA := [][2]string{
		{"https://example.com/vocab/shared", "false"},
		{"https://example.com/vocab/b", "true"},
		{"https://example.com/vocab/d", "false"},
	}
	rightEntriesB := [][2]string{
		{"https://example.com/vocab/d", "false"},
		{"https://example.com/vocab/shared", "false"},
		{"https://example.com/vocab/b", "true"},
	}

	pairA := checkVocabularyChanges(buildVocab(leftEntriesA), buildVocab(rightEntriesA))
	pairB := checkVocabularyChanges(buildVocab(leftEntriesB), buildVocab(rightEntriesB))

	assert.NotEmpty(t, pairA, "expected at least one vocabulary change record")
	assert.Equal(t, len(pairA), len(pairB),
		"record count must be independent of source insertion order")

	for i := range pairA {
		assert.Equal(t, pairA[i].Property, pairB[i].Property,
			"record[%d] Property must be independent of source insertion order", i)
		assert.Equal(t, pairA[i].ChangeType, pairB[i].ChangeType,
			"record[%d] ChangeType must be independent of source insertion order", i)
		assert.Equal(t, pairA[i].Original, pairB[i].Original,
			"record[%d] Original payload must be independent of source insertion order", i)
		assert.Equal(t, pairA[i].New, pairB[i].New,
			"record[%d] New payload must be independent of source insertion order", i)
	}
}

// TestSchemaChanges_MarshalJSONYAML_NewContainers verifies that the two new
// change containers (ContentSchemaChanges and VocabularyChanges) serialize
// under their contract keyword keys in both JSON and YAML output. This
// protects the public serialized API against tag drift.
func TestEvalonGolden_SchemaChanges_MarshalJSONYAML_NewContainers(t *testing.T) {
	sc := &SchemaChanges{
		VocabularyChanges: []*Change{
			{
				Property:   base.VocabularyLabel,
				ChangeType: PropertyAdded,
				New:        "https://example.com/vocab=true",
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
					},
				},
			},
		},
	}

	jsonBytes, err := json.Marshal(sc)
	assert.NoError(t, err, "SchemaChanges must marshal to JSON without error")
	jsonOutput := string(jsonBytes)
	assert.Contains(t, jsonOutput, `"$vocabulary"`,
		"JSON output must use $vocabulary as the container key")
	assert.Contains(t, jsonOutput, `"contentSchema"`,
		"JSON output must use contentSchema as the container key")
	assert.NotContains(t, strings.ReplaceAll(jsonOutput, "$vocabulary", ""), `"vocabulary"`,
		"JSON output must NOT use plain vocabulary (must retain the leading $)")

	yamlBytes, err := yaml.Marshal(sc)
	assert.NoError(t, err, "SchemaChanges must marshal to YAML without error")
	yamlOutput := string(yamlBytes)
	assert.Contains(t, yamlOutput, "$vocabulary:",
		"YAML output must use $vocabulary as the container key")
	assert.Contains(t, yamlOutput, "contentSchema:",
		"YAML output must use contentSchema as the container key")
}

// TestCheckVocabularyChanges_TransitionMatrix exhaustively exercises the
// $vocabulary comparator across the full transition matrix: both-nil,
// left-nil-with-right-entries, right-nil-with-left-entries, per-entry
// added / removed / modified inside otherwise non-empty maps, identical
// content, and a mixed simultaneous scenario. For every case the test
// asserts (a) no panic on nil-receiver iterators, (b) the correct number
// of emitted records, (c) every record's Property is exactly base.VocabularyLabel,
// (d) each ChangeType matches its transition, and (e) the applicable Original
// and New payloads preserve both the URI and the boolean value.
func TestEvalonGolden_CheckVocabularyChanges_TransitionMatrix(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	buildVocab := func(entries [][2]any) *orderedmap.Map[low.KeyReference[string], low.ValueReference[bool]] {
		m := orderedmap.New[low.KeyReference[string], low.ValueReference[bool]]()
		for _, entry := range entries {
			m.Set(
				low.KeyReference[string]{Value: entry[0].(string)},
				low.ValueReference[bool]{Value: entry[1].(bool)},
			)
		}
		return m
	}

	type expected struct {
		uri        string
		changeType int
		original   string
		new        string
	}

	cases := []struct {
		name          string
		left          *orderedmap.Map[low.KeyReference[string], low.ValueReference[bool]]
		right         *orderedmap.Map[low.KeyReference[string], low.ValueReference[bool]]
		expectRecords []expected
	}{
		{
			name:          "both nil emits no records and does not panic",
			left:          nil,
			right:         nil,
			expectRecords: nil,
		},
		{
			name: "left nil with right entries emits per-entry additions",
			left: nil,
			right: buildVocab([][2]any{
				{"https://example.com/vocab/a", true},
				{"https://example.com/vocab/b", false},
			}),
			expectRecords: []expected{
				{uri: "https://example.com/vocab/a", changeType: PropertyAdded, new: "https://example.com/vocab/a=true"},
				{uri: "https://example.com/vocab/b", changeType: PropertyAdded, new: "https://example.com/vocab/b=false"},
			},
		},
		{
			name: "right nil with left entries emits per-entry removals",
			left: buildVocab([][2]any{
				{"https://example.com/vocab/a", true},
				{"https://example.com/vocab/b", false},
			}),
			right: nil,
			expectRecords: []expected{
				{uri: "https://example.com/vocab/a", changeType: PropertyRemoved, original: "https://example.com/vocab/a=true"},
				{uri: "https://example.com/vocab/b", changeType: PropertyRemoved, original: "https://example.com/vocab/b=false"},
			},
		},
		{
			name: "per entry added inside otherwise non-empty map",
			left: buildVocab([][2]any{
				{"https://example.com/vocab/shared", true},
			}),
			right: buildVocab([][2]any{
				{"https://example.com/vocab/shared", true},
				{"https://example.com/vocab/added", false},
			}),
			expectRecords: []expected{
				{uri: "https://example.com/vocab/added", changeType: PropertyAdded, new: "https://example.com/vocab/added=false"},
			},
		},
		{
			name: "per entry removed inside otherwise non-empty map",
			left: buildVocab([][2]any{
				{"https://example.com/vocab/shared", true},
				{"https://example.com/vocab/removed", false},
			}),
			right: buildVocab([][2]any{
				{"https://example.com/vocab/shared", true},
			}),
			expectRecords: []expected{
				{uri: "https://example.com/vocab/removed", changeType: PropertyRemoved, original: "https://example.com/vocab/removed=false"},
			},
		},
		{
			name: "per entry modified when boolean value differs",
			left: buildVocab([][2]any{
				{"https://example.com/vocab/toggle", true},
			}),
			right: buildVocab([][2]any{
				{"https://example.com/vocab/toggle", false},
			}),
			expectRecords: []expected{
				{
					uri:        "https://example.com/vocab/toggle",
					changeType: Modified,
					original:   "https://example.com/vocab/toggle=true",
					new:        "https://example.com/vocab/toggle=false",
				},
			},
		},
		{
			name: "identical maps emit no records",
			left: buildVocab([][2]any{
				{"https://example.com/vocab/a", true},
				{"https://example.com/vocab/b", false},
			}),
			right: buildVocab([][2]any{
				{"https://example.com/vocab/a", true},
				{"https://example.com/vocab/b", false},
			}),
			expectRecords: nil,
		},
		{
			name: "mixed simultaneous add remove and modify alongside unchanged entry",
			left: buildVocab([][2]any{
				{"https://example.com/vocab/unchanged", true},
				{"https://example.com/vocab/removed", true},
				{"https://example.com/vocab/toggle", true},
			}),
			right: buildVocab([][2]any{
				{"https://example.com/vocab/unchanged", true},
				{"https://example.com/vocab/added", false},
				{"https://example.com/vocab/toggle", false},
			}),
			expectRecords: []expected{
				{uri: "https://example.com/vocab/added", changeType: PropertyAdded, new: "https://example.com/vocab/added=false"},
				{uri: "https://example.com/vocab/removed", changeType: PropertyRemoved, original: "https://example.com/vocab/removed=true"},
				{
					uri:        "https://example.com/vocab/toggle",
					changeType: Modified,
					original:   "https://example.com/vocab/toggle=true",
					new:        "https://example.com/vocab/toggle=false",
				},
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			var records []*Change
			assert.NotPanics(t, func() {
				records = checkVocabularyChanges(testCase.left, testCase.right)
			}, "checkVocabularyChanges must not panic on nil receivers")

			assert.Len(t, records, len(testCase.expectRecords),
				"unexpected number of vocabulary change records")

			byURI := make(map[string]*Change, len(records))
			for _, record := range records {
				assert.Equal(t, base.VocabularyLabel, record.Property,
					"every vocabulary record must carry the $vocabulary label")
				var uri string
				switch record.ChangeType {
				case PropertyAdded:
					uri = strings.SplitN(record.New, "=", 2)[0]
				case PropertyRemoved:
					uri = strings.SplitN(record.Original, "=", 2)[0]
				case Modified:
					uri = strings.SplitN(record.Original, "=", 2)[0]
				}
				byURI[uri] = record
			}

			for _, want := range testCase.expectRecords {
				got, ok := byURI[want.uri]
				assert.True(t, ok, "missing record for URI %q", want.uri)
				if !ok {
					continue
				}
				assert.Equal(t, want.changeType, got.ChangeType,
					"unexpected change type for URI %q", want.uri)
				if want.original != "" {
					assert.Equal(t, want.original, got.Original,
						"Original payload for URI %q must preserve both URI and boolean value", want.uri)
				}
				if want.new != "" {
					assert.Equal(t, want.new, got.New,
						"New payload for URI %q must preserve both URI and boolean value", want.uri)
				}
			}
		})
	}
}

// TestEvalonGolden_CompareSchemas_ContentSchema_NoChange asserts that
// identical contentSchema values on both sides emit no change record for
// contentSchema. This complements the Added/Removed/Modified cases and
// closes the four-branch transition matrix for contentSchema.
func TestEvalonGolden_CompareSchemas_ContentSchema_NoChange(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `openapi: "3.1.0"
info:
  title: left
  version: "1.0"
components:
  schemas:
    Pet:
      type: string
      contentMediaType: application/json
      contentSchema:
        type: object`

	right := `openapi: "3.1.0"
info:
  title: right
  version: "1.0"
components:
  schemas:
    Pet:
      type: string
      contentMediaType: application/json
      contentSchema:
        type: object`

	leftDoc, rightDoc := test_BuildDoc(left, right)

	lSchemaProxy := leftDoc.Components.Value.FindSchema("Pet").Value
	rSchemaProxy := rightDoc.Components.Value.FindSchema("Pet").Value

	changes := CompareSchemas(lSchemaProxy, rSchemaProxy)
	if changes != nil {
		for _, change := range changes.GetAllChanges() {
			assert.NotEqual(t, PropContentSchema, change.Property,
				"identical contentSchema values must not emit a change record")
		}
	}
}

// TestEvalonGolden_SchemaChanges_AggregateCounts covers the aggregate-count
// contract for the two new containers: TotalChanges must include contributions
// from ContentSchemaChanges (recursively) and every entry in VocabularyChanges,
// while TotalBreakingChanges must include only the breaking-flagged entries.
// The test uses a mix of populated and nil containers with mixed breaking
// flags to prove no omission, no double-count, and no panic on nil.
func TestEvalonGolden_SchemaChanges_AggregateCounts(t *testing.T) {
	nonBreakingContentChange := &Change{
		Property:   PropContentSchema,
		ChangeType: Modified,
		Breaking:   false,
	}
	breakingContentChange := &Change{
		Property:   PropContentSchema,
		ChangeType: Modified,
		Breaking:   true,
	}
	breakingVocabAdd := &Change{
		Property:   base.VocabularyLabel,
		ChangeType: PropertyAdded,
		Breaking:   true,
		New:        "https://example.com/vocab/a=true",
	}
	nonBreakingVocabAdd := &Change{
		Property:   base.VocabularyLabel,
		ChangeType: PropertyAdded,
		Breaking:   false,
		New:        "https://example.com/vocab/b=false",
	}
	breakingVocabModify := &Change{
		Property:   base.VocabularyLabel,
		ChangeType: Modified,
		Breaking:   true,
		Original:   "https://example.com/vocab/c=true",
		New:        "https://example.com/vocab/c=false",
	}

	type totals struct {
		total    int
		breaking int
	}

	// newSchemaChanges mirrors the initialization CompareSchemas performs on
	// the outer container so tests exercise the two new containers (not the
	// unrelated nil-safety of the embedded PropertyChanges struct).
	newSchemaChanges := func() *SchemaChanges {
		return &SchemaChanges{PropertyChanges: &PropertyChanges{}}
	}

	cases := []struct {
		name     string
		build    func() *SchemaChanges
		expected totals
	}{
		{
			name: "both new containers absent produces zero totals and does not panic",
			build: func() *SchemaChanges {
				return newSchemaChanges()
			},
			expected: totals{total: 0, breaking: 0},
		},
		{
			name: "contentSchema-only recursive contribution counts once",
			build: func() *SchemaChanges {
				sc := newSchemaChanges()
				sc.ContentSchemaChanges = &SchemaChanges{
					PropertyChanges: &PropertyChanges{
						Changes: []*Change{breakingContentChange, nonBreakingContentChange},
					},
				}
				return sc
			},
			expected: totals{total: 2, breaking: 1},
		},
		{
			name: "vocabulary-only contribution counts every entry with only breaking flagged",
			build: func() *SchemaChanges {
				sc := newSchemaChanges()
				sc.VocabularyChanges = []*Change{
					breakingVocabAdd,
					nonBreakingVocabAdd,
					breakingVocabModify,
				}
				return sc
			},
			expected: totals{total: 3, breaking: 2},
		},
		{
			name: "both containers together sum without double-counting",
			build: func() *SchemaChanges {
				sc := newSchemaChanges()
				sc.ContentSchemaChanges = &SchemaChanges{
					PropertyChanges: &PropertyChanges{
						Changes: []*Change{breakingContentChange, nonBreakingContentChange},
					},
				}
				sc.VocabularyChanges = []*Change{
					breakingVocabAdd,
					nonBreakingVocabAdd,
					breakingVocabModify,
				}
				return sc
			},
			expected: totals{total: 5, breaking: 3},
		},
		{
			name: "empty vocabulary slice contributes nothing",
			build: func() *SchemaChanges {
				sc := newSchemaChanges()
				sc.VocabularyChanges = []*Change{}
				return sc
			},
			expected: totals{total: 0, breaking: 0},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			sc := testCase.build()
			var total, breaking int
			assert.NotPanics(t, func() {
				total = sc.TotalChanges()
				breaking = sc.TotalBreakingChanges()
			}, "aggregate accessors must be nil-safe")

			assert.Equal(t, testCase.expected.total, total,
				"TotalChanges must include contentSchema (recursive) plus every vocabulary entry")
			assert.Equal(t, testCase.expected.breaking, breaking,
				"TotalBreakingChanges must include recursive contentSchema breaking plus only breaking vocabulary entries")
		})
	}
}

// TestEvalonGolden_DefaultBreakingRules_NewKeywordPolarities asserts the exact
// polarity matrix for the three new Schema keywords: $comment is
// documentation-only so every transition is non-breaking, while contentSchema
// and $vocabulary affect validation so every transition is breaking. This
// prevents accidental factory drift without inspecting the construction path.
func TestEvalonGolden_DefaultBreakingRules_NewKeywordPolarities(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config, "default rules factory must return a non-nil config")
	assert.NotNil(t, config.Schema, "default config must include Schema rules")

	polarity := func(rule *BreakingChangeRule, field string) bool {
		var ptr *bool
		switch field {
		case "Added":
			ptr = rule.Added
		case "Modified":
			ptr = rule.Modified
		case "Removed":
			ptr = rule.Removed
		}
		if ptr == nil {
			t.Fatalf("%s field must be set explicitly, not nil", field)
		}
		return *ptr
	}

	assert.NotNil(t, config.Schema.Comment, "Schema.Comment rule must be present")
	assert.False(t, polarity(config.Schema.Comment, "Added"),
		"$comment addition is documentation-only and must be non-breaking")
	assert.False(t, polarity(config.Schema.Comment, "Modified"),
		"$comment modification is documentation-only and must be non-breaking")
	assert.False(t, polarity(config.Schema.Comment, "Removed"),
		"$comment removal is documentation-only and must be non-breaking")

	assert.NotNil(t, config.Schema.ContentSchema, "Schema.ContentSchema rule must be present")
	assert.True(t, polarity(config.Schema.ContentSchema, "Added"),
		"contentSchema addition affects validation and must be breaking")
	assert.True(t, polarity(config.Schema.ContentSchema, "Modified"),
		"contentSchema modification affects validation and must be breaking")
	assert.True(t, polarity(config.Schema.ContentSchema, "Removed"),
		"contentSchema removal affects validation and must be breaking")

	assert.NotNil(t, config.Schema.Vocabulary, "Schema.Vocabulary rule must be present")
	assert.True(t, polarity(config.Schema.Vocabulary, "Added"),
		"$vocabulary addition affects validation and must be breaking")
	assert.True(t, polarity(config.Schema.Vocabulary, "Modified"),
		"$vocabulary modification affects validation and must be breaking")
	assert.True(t, polarity(config.Schema.Vocabulary, "Removed"),
		"$vocabulary removal affects validation and must be breaking")
}
