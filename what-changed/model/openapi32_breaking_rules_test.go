// Copyright 2022-2025 Princess Beef Heavy Industries, LLC / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	// additionalOperations follows the same pattern as other HTTP methods:
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
	// adding and removing are breaking, modifying is not
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

	// itemEncoding: all changes are breaking
	assert.True(t, *config.MediaType.ItemEncoding.Added)
	assert.True(t, *config.MediaType.ItemEncoding.Modified)
	assert.True(t, *config.MediaType.ItemEncoding.Removed)
}

func TestDefaultBreakingRules_SecurityScheme_OAuth2MetadataUrl(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.SecurityScheme)
	assert.NotNil(t, config.SecurityScheme.OAuth2MetadataUrl)

	// oauth2MetadataUrl changes are non-breaking
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Added)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Modified)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Removed)
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
	assert.False(t, *rule.Modified)
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

func TestGetRule_OpenAPI32_SecuritySchemeOAuth2MetadataUrl(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	rule := config.GetRule(CompSecurityScheme, PropOAuth2MetadataUrl)
	assert.NotNil(t, rule)
	assert.False(t, *rule.Added)
	assert.False(t, *rule.Modified)
	assert.False(t, *rule.Removed)
}

func TestIsBreaking_OpenAPI32_PathItemQuery(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	// adding query is not breaking
	assert.False(t, config.IsBreaking(CompPathItem, PropQuery, ChangeTypeAdded))
	// modifying query is not breaking
	assert.False(t, config.IsBreaking(CompPathItem, PropQuery, ChangeTypeModified))
	// removing query is breaking
	assert.True(t, config.IsBreaking(CompPathItem, PropQuery, ChangeTypeRemoved))
}

func TestIsBreaking_OpenAPI32_PathItemAdditionalOperations(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	// adding additionalOperations is not breaking
	assert.False(t, config.IsBreaking(CompPathItem, PropAdditionalOperations, ChangeTypeAdded))
	// modifying additionalOperations is not breaking
	assert.False(t, config.IsBreaking(CompPathItem, PropAdditionalOperations, ChangeTypeModified))
	// removing additionalOperations is breaking
	assert.True(t, config.IsBreaking(CompPathItem, PropAdditionalOperations, ChangeTypeRemoved))
}

func TestIsBreaking_OpenAPI32_MediaTypeItemSchema(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	// adding itemSchema is breaking
	assert.True(t, config.IsBreaking(CompMediaType, PropItemSchema, ChangeTypeAdded))
	// modifying itemSchema is not breaking
	assert.False(t, config.IsBreaking(CompMediaType, PropItemSchema, ChangeTypeModified))
	// removing itemSchema is breaking
	assert.True(t, config.IsBreaking(CompMediaType, PropItemSchema, ChangeTypeRemoved))
}

func TestIsBreaking_OpenAPI32_MediaTypeItemEncoding(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	// all changes to itemEncoding are breaking
	assert.True(t, config.IsBreaking(CompMediaType, PropItemEncoding, ChangeTypeAdded))
	assert.True(t, config.IsBreaking(CompMediaType, PropItemEncoding, ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompMediaType, PropItemEncoding, ChangeTypeRemoved))
}

func TestIsBreaking_OpenAPI32_SecuritySchemeOAuth2MetadataUrl(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	// all changes to oauth2MetadataUrl are non-breaking
	assert.False(t, config.IsBreaking(CompSecurityScheme, PropOAuth2MetadataUrl, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompSecurityScheme, PropOAuth2MetadataUrl, ChangeTypeModified))
	assert.False(t, config.IsBreaking(CompSecurityScheme, PropOAuth2MetadataUrl, ChangeTypeRemoved))
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
	// removal should remain unchanged
	assert.True(t, *config.PathItem.Query.Removed)
}

func TestMerge_OpenAPI32_AdditionalOperationsOverride(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	// override to make additionalOperations non-breaking for removal
	override := &BreakingRulesConfig{
		PathItem: &PathItemRules{
			AdditionalOperations: &BreakingChangeRule{
				Removed: boolPtr(false),
			},
		},
	}
	config.Merge(override)

	// removal is now non-breaking
	assert.False(t, *config.PathItem.AdditionalOperations.Removed)
	// added should remain unchanged
	assert.False(t, *config.PathItem.AdditionalOperations.Added)
}

func TestMerge_OpenAPI32_MediaTypeItemSchemaOverride(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	// override to make itemSchema modification breaking
	override := &BreakingRulesConfig{
		MediaType: &MediaTypeRules{
			ItemSchema: &BreakingChangeRule{
				Modified: boolPtr(true),
			},
		},
	}
	config.Merge(override)

	// modification is now breaking
	assert.True(t, *config.MediaType.ItemSchema.Modified)
	// other fields should remain unchanged
	assert.True(t, *config.MediaType.ItemSchema.Added)
	assert.True(t, *config.MediaType.ItemSchema.Removed)
}

func TestMerge_OpenAPI32_MediaTypeItemEncodingOverride(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	// override to make itemEncoding non-breaking
	override := &BreakingRulesConfig{
		MediaType: &MediaTypeRules{
			ItemEncoding: &BreakingChangeRule{
				Added:    boolPtr(false),
				Modified: boolPtr(false),
				Removed:  boolPtr(false),
			},
		},
	}
	config.Merge(override)

	assert.False(t, *config.MediaType.ItemEncoding.Added)
	assert.False(t, *config.MediaType.ItemEncoding.Modified)
	assert.False(t, *config.MediaType.ItemEncoding.Removed)
}

func TestMerge_OpenAPI32_SecuritySchemeOAuth2MetadataUrlOverride(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	// override to make oauth2MetadataUrl changes breaking
	override := &BreakingRulesConfig{
		SecurityScheme: &SecuritySchemeRules{
			OAuth2MetadataUrl: &BreakingChangeRule{
				Modified: boolPtr(true),
			},
		},
	}
	config.Merge(override)

	// modification is now breaking
	assert.True(t, *config.SecurityScheme.OAuth2MetadataUrl.Modified)
	// other fields should remain unchanged
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Added)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Removed)
}

func TestBreakingHelpers_OpenAPI32_Properties(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	// Test BreakingAdded/BreakingModified/BreakingRemoved helpers for new properties

	// PathItem query
	assert.False(t, BreakingAdded(CompPathItem, PropQuery))
	assert.False(t, BreakingModified(CompPathItem, PropQuery))
	assert.True(t, BreakingRemoved(CompPathItem, PropQuery))

	// PathItem additionalOperations
	assert.False(t, BreakingAdded(CompPathItem, PropAdditionalOperations))
	assert.False(t, BreakingModified(CompPathItem, PropAdditionalOperations))
	assert.True(t, BreakingRemoved(CompPathItem, PropAdditionalOperations))

	// MediaType itemSchema
	assert.True(t, BreakingAdded(CompMediaType, PropItemSchema))
	assert.False(t, BreakingModified(CompMediaType, PropItemSchema))
	assert.True(t, BreakingRemoved(CompMediaType, PropItemSchema))

	// MediaType itemEncoding
	assert.True(t, BreakingAdded(CompMediaType, PropItemEncoding))
	assert.True(t, BreakingModified(CompMediaType, PropItemEncoding))
	assert.True(t, BreakingRemoved(CompMediaType, PropItemEncoding))

	// SecurityScheme oauth2MetadataUrl
	assert.False(t, BreakingAdded(CompSecurityScheme, PropOAuth2MetadataUrl))
	assert.False(t, BreakingModified(CompSecurityScheme, PropOAuth2MetadataUrl))
	assert.False(t, BreakingRemoved(CompSecurityScheme, PropOAuth2MetadataUrl))
}

func TestPropertyConstants_OpenAPI32(t *testing.T) {
	// Verify the new property constants have the correct values
	assert.Equal(t, "query", PropQuery)
	assert.Equal(t, "additionalOperations", PropAdditionalOperations)
	assert.Equal(t, "itemSchema", PropItemSchema)
	assert.Equal(t, "itemEncoding", PropItemEncoding)
	assert.Equal(t, "oauth2MetadataUrl", PropOAuth2MetadataUrl)
}
