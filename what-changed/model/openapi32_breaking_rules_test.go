// Copyright 2022-2025 Princess Beef Heavy Industries, LLC / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for OpenAPI 3.2 properties in the breaking rules system.

func TestDefaultBreakingRules_PathItem_Query(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.PathItem)
	assert.NotNil(t, config.PathItem.Query)

	// query follows the same pattern as other HTTP methods:
	// adding is not breaking, removing is breaking
	assert.False(t, *config.PathItem.Query.Added)
	assert.False(t, *config.PathItem.Query.Modified)
	assert.True(t, *config.PathItem.Query.Removed)
}

func TestDefaultBreakingRules_PathItem_AdditionalOperations(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.PathItem)
	assert.NotNil(t, config.PathItem.AdditionalOperations)

	// additional operations follow the same pattern as other HTTP methods:
	// adding is not breaking, removing is breaking
	assert.False(t, *config.PathItem.AdditionalOperations.Added)
	assert.False(t, *config.PathItem.AdditionalOperations.Modified)
	assert.True(t, *config.PathItem.AdditionalOperations.Removed)
}

func TestDefaultBreakingRules_MediaType_ItemSchema(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.MediaType)
	assert.NotNil(t, config.MediaType.ItemSchema)

	// itemSchema follows the same pattern as schema:
	// adding is breaking, removing is breaking
	assert.True(t, *config.MediaType.ItemSchema.Added)
	assert.False(t, *config.MediaType.ItemSchema.Modified)
	assert.True(t, *config.MediaType.ItemSchema.Removed)
}

func TestDefaultBreakingRules_MediaType_ItemEncoding(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.MediaType)
	assert.NotNil(t, config.MediaType.ItemEncoding)

	// itemEncoding changes are all breaking
	assert.True(t, *config.MediaType.ItemEncoding.Added)
	assert.True(t, *config.MediaType.ItemEncoding.Modified)
	assert.True(t, *config.MediaType.ItemEncoding.Removed)
}

func TestDefaultBreakingRules_SecurityScheme_OAuth2MetadataURL(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.SecurityScheme)
	assert.NotNil(t, config.SecurityScheme.OAuth2MetadataURL)

	// oauth2MetadataUrl changes are non-breaking
	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Added)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Modified)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Removed)
}

func TestGetRule_OpenAPI32_PathItemQuery(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	rule := config.GetRule(CompPathItem, PropQuery)
	assert.NotNil(t, rule)
	assert.False(t, *rule.Added)
	assert.True(t, *rule.Removed)
}

func TestGetRule_OpenAPI32_PathItemAdditionalOperations(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	rule := config.GetRule(CompPathItem, PropAdditionalOperations)
	assert.NotNil(t, rule)
	assert.False(t, *rule.Added)
	assert.True(t, *rule.Removed)
}

func TestGetRule_OpenAPI32_MediaTypeItemSchema(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	rule := config.GetRule(CompMediaType, PropItemSchema)
	assert.NotNil(t, rule)
	assert.True(t, *rule.Added)
	assert.True(t, *rule.Removed)
}

func TestGetRule_OpenAPI32_MediaTypeItemEncoding(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	rule := config.GetRule(CompMediaType, PropItemEncoding)
	assert.NotNil(t, rule)
	assert.True(t, *rule.Added)
	assert.True(t, *rule.Modified)
	assert.True(t, *rule.Removed)
}

func TestGetRule_OpenAPI32_SecuritySchemeOAuth2MetadataURL(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	rule := config.GetRule(CompSecurityScheme, PropOAuth2MetadataURL)
	assert.NotNil(t, rule)
	assert.False(t, *rule.Added)
	assert.False(t, *rule.Modified)
	assert.False(t, *rule.Removed)
}

func TestIsBreaking_OpenAPI32_PathItemQuery(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	assert.False(t, config.IsBreaking(CompPathItem, PropQuery, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompPathItem, PropQuery, ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompPathItem, PropQuery, ChangeTypeRemoved))
}

func TestIsBreaking_OpenAPI32_PathItemAdditionalOperations(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	assert.False(t, config.IsBreaking(CompPathItem, PropAdditionalOperations, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompPathItem, PropAdditionalOperations, ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompPathItem, PropAdditionalOperations, ChangeTypeRemoved))
}

func TestIsBreaking_OpenAPI32_MediaTypeItemSchema(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	assert.True(t, config.IsBreaking(CompMediaType, PropItemSchema, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompMediaType, PropItemSchema, ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompMediaType, PropItemSchema, ChangeTypeRemoved))
}

func TestIsBreaking_OpenAPI32_MediaTypeItemEncoding(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	assert.True(t, config.IsBreaking(CompMediaType, PropItemEncoding, ChangeTypeAdded))
	assert.True(t, config.IsBreaking(CompMediaType, PropItemEncoding, ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompMediaType, PropItemEncoding, ChangeTypeRemoved))
}

func TestIsBreaking_OpenAPI32_SecuritySchemeOAuth2MetadataURL(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	assert.False(t, config.IsBreaking(CompSecurityScheme, PropOAuth2MetadataURL, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompSecurityScheme, PropOAuth2MetadataURL, ChangeTypeModified))
	assert.False(t, config.IsBreaking(CompSecurityScheme, PropOAuth2MetadataURL, ChangeTypeRemoved))
}

func TestMerge_OpenAPI32_PathItemQueryOverride(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	// verify defaults
	assert.False(t, *config.PathItem.Query.Added)
	assert.True(t, *config.PathItem.Query.Removed)

	// override to make query addition breaking
	override := &BreakingRulesConfig{
		PathItem: &PathItemRules{
			Query: &BreakingChangeRule{
				Added: boolPtr(true),
			},
		},
	}
	config.Merge(override)

	// query addition is now breaking
	assert.True(t, *config.PathItem.Query.Added)
	// removed should remain unchanged
	assert.True(t, *config.PathItem.Query.Removed)
}

func TestMerge_OpenAPI32_AdditionalOperationsOverride(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	// override to make additional operations non-breaking on removal
	override := &BreakingRulesConfig{
		PathItem: &PathItemRules{
			AdditionalOperations: &BreakingChangeRule{
				Removed: boolPtr(false),
			},
		},
	}
	config.Merge(override)

	assert.False(t, *config.PathItem.AdditionalOperations.Removed)
	// other fields should remain unchanged
	assert.False(t, *config.PathItem.AdditionalOperations.Added)
}

func TestMerge_OpenAPI32_MediaTypeItemSchemaOverride(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	// override to make itemSchema non-breaking
	override := &BreakingRulesConfig{
		MediaType: &MediaTypeRules{
			ItemSchema: &BreakingChangeRule{
				Added:   boolPtr(false),
				Removed: boolPtr(false),
			},
		},
	}
	config.Merge(override)

	assert.False(t, *config.MediaType.ItemSchema.Added)
	assert.False(t, *config.MediaType.ItemSchema.Removed)
}

func TestMerge_OpenAPI32_MediaTypeItemEncodingOverride(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	// override to make itemEncoding modification non-breaking
	override := &BreakingRulesConfig{
		MediaType: &MediaTypeRules{
			ItemEncoding: &BreakingChangeRule{
				Modified: boolPtr(false),
			},
		},
	}
	config.Merge(override)

	assert.True(t, *config.MediaType.ItemEncoding.Added)     // unchanged
	assert.False(t, *config.MediaType.ItemEncoding.Modified) // overridden
	assert.True(t, *config.MediaType.ItemEncoding.Removed)   // unchanged
}

func TestMerge_OpenAPI32_SecuritySchemeOAuth2MetadataURLOverride(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	// override to make oauth2MetadataUrl breaking
	override := &BreakingRulesConfig{
		SecurityScheme: &SecuritySchemeRules{
			OAuth2MetadataURL: &BreakingChangeRule{
				Added:    boolPtr(true),
				Modified: boolPtr(true),
				Removed:  boolPtr(true),
			},
		},
	}
	config.Merge(override)

	assert.True(t, *config.SecurityScheme.OAuth2MetadataURL.Added)
	assert.True(t, *config.SecurityScheme.OAuth2MetadataURL.Modified)
	assert.True(t, *config.SecurityScheme.OAuth2MetadataURL.Removed)
}
