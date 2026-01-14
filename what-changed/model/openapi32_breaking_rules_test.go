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

	assert.True(t, *config.MediaType.ItemEncoding.Added)
	assert.False(t, *config.MediaType.ItemEncoding.Modified)
	assert.True(t, *config.MediaType.ItemEncoding.Removed)
}

func TestDefaultBreakingRules_SecurityScheme_OAuth2MetadataUrl(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.SecurityScheme)
	assert.NotNil(t, config.SecurityScheme.OAuth2MetadataUrl)

	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Added)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Modified)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Removed)
}

func TestGetRule_OpenAPI32_PathItemProperties(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	queryRule := config.GetRule("pathItem", "query")
	assert.NotNil(t, queryRule, "expected rule for pathItem/query")
	assert.False(t, *queryRule.Added)
	assert.True(t, *queryRule.Removed)

	additionalOpsRule := config.GetRule("pathItem", "additionalOperations")
	assert.NotNil(t, additionalOpsRule, "expected rule for pathItem/additionalOperations")
	assert.False(t, *additionalOpsRule.Added)
	assert.True(t, *additionalOpsRule.Removed)
}

func TestGetRule_OpenAPI32_MediaTypeProperties(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	itemSchemaRule := config.GetRule("mediaType", "itemSchema")
	assert.NotNil(t, itemSchemaRule, "expected rule for mediaType/itemSchema")
	assert.True(t, *itemSchemaRule.Added)
	assert.True(t, *itemSchemaRule.Removed)

	itemEncodingRule := config.GetRule("mediaType", "itemEncoding")
	assert.NotNil(t, itemEncodingRule, "expected rule for mediaType/itemEncoding")
	assert.True(t, *itemEncodingRule.Added)
	assert.True(t, *itemEncodingRule.Removed)
}

func TestGetRule_OpenAPI32_SecuritySchemeProperties(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	oauth2MetadataUrlRule := config.GetRule("securityScheme", "oauth2MetadataUrl")
	assert.NotNil(t, oauth2MetadataUrlRule, "expected rule for securityScheme/oauth2MetadataUrl")
	assert.False(t, *oauth2MetadataUrlRule.Added)
	assert.False(t, *oauth2MetadataUrlRule.Modified)
	assert.False(t, *oauth2MetadataUrlRule.Removed)
}

func TestIsBreaking_OpenAPI32_PathItem(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	assert.False(t, config.IsBreaking("pathItem", "query", ChangeTypeAdded))
	assert.False(t, config.IsBreaking("pathItem", "query", ChangeTypeModified))
	assert.True(t, config.IsBreaking("pathItem", "query", ChangeTypeRemoved))

	assert.False(t, config.IsBreaking("pathItem", "additionalOperations", ChangeTypeAdded))
	assert.False(t, config.IsBreaking("pathItem", "additionalOperations", ChangeTypeModified))
	assert.True(t, config.IsBreaking("pathItem", "additionalOperations", ChangeTypeRemoved))
}

func TestIsBreaking_OpenAPI32_MediaType(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	assert.True(t, config.IsBreaking("mediaType", "itemSchema", ChangeTypeAdded))
	assert.False(t, config.IsBreaking("mediaType", "itemSchema", ChangeTypeModified))
	assert.True(t, config.IsBreaking("mediaType", "itemSchema", ChangeTypeRemoved))

	assert.True(t, config.IsBreaking("mediaType", "itemEncoding", ChangeTypeAdded))
	assert.False(t, config.IsBreaking("mediaType", "itemEncoding", ChangeTypeModified))
	assert.True(t, config.IsBreaking("mediaType", "itemEncoding", ChangeTypeRemoved))
}

func TestIsBreaking_OpenAPI32_SecurityScheme(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	assert.False(t, config.IsBreaking("securityScheme", "oauth2MetadataUrl", ChangeTypeAdded))
	assert.False(t, config.IsBreaking("securityScheme", "oauth2MetadataUrl", ChangeTypeModified))
	assert.False(t, config.IsBreaking("securityScheme", "oauth2MetadataUrl", ChangeTypeRemoved))
}

func TestMerge_OpenAPI32_PathItemOverride(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	assert.False(t, *config.PathItem.Query.Added)
	assert.True(t, *config.PathItem.Query.Removed)

	override := &BreakingRulesConfig{
		PathItem: &PathItemRules{
			Query: &BreakingChangeRule{
				Added:   boolPtr(true),
				Removed: boolPtr(false),
			},
		},
	}
	config.Merge(override)

	assert.True(t, *config.PathItem.Query.Added)
	assert.False(t, *config.PathItem.Query.Removed)
}

func TestMerge_OpenAPI32_MediaTypeOverride(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	assert.True(t, *config.MediaType.ItemSchema.Added)
	assert.True(t, *config.MediaType.ItemSchema.Removed)

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

func TestMerge_OpenAPI32_SecuritySchemeOverride(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Modified)

	override := &BreakingRulesConfig{
		SecurityScheme: &SecuritySchemeRules{
			OAuth2MetadataUrl: &BreakingChangeRule{
				Modified: boolPtr(true),
			},
		},
	}
	config.Merge(override)

	assert.True(t, *config.SecurityScheme.OAuth2MetadataUrl.Modified)
}

func TestBreakingHelpers_OpenAPI32(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	defer ResetDefaultBreakingRules()
	defer ResetActiveBreakingRulesConfig()

	assert.False(t, BreakingAdded(CompPathItem, "query"))
	assert.True(t, BreakingRemoved(CompPathItem, "query"))

	assert.False(t, BreakingAdded(CompPathItem, "additionalOperations"))
	assert.True(t, BreakingRemoved(CompPathItem, "additionalOperations"))

	assert.True(t, BreakingAdded(CompMediaType, "itemSchema"))
	assert.True(t, BreakingRemoved(CompMediaType, "itemSchema"))

	assert.True(t, BreakingAdded(CompMediaType, "itemEncoding"))
	assert.True(t, BreakingRemoved(CompMediaType, "itemEncoding"))

	assert.False(t, BreakingAdded(CompSecurityScheme, "oauth2MetadataUrl"))
	assert.False(t, BreakingModified(CompSecurityScheme, "oauth2MetadataUrl"))
	assert.False(t, BreakingRemoved(CompSecurityScheme, "oauth2MetadataUrl"))
}
