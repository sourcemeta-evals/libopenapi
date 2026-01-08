// Copyright 2022-2025 Princess Beef Heavy Industries, LLC / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultBreakingRules_PathItem_Query(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.PathItem)
	assert.NotNil(t, config.PathItem.Query)

	assert.False(t, *config.PathItem.Query.Added)
	assert.False(t, *config.PathItem.Query.Modified)
	assert.True(t, *config.PathItem.Query.Removed)
}

func TestDefaultBreakingRules_PathItem_AdditionalOperations(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.PathItem)
	assert.NotNil(t, config.PathItem.AdditionalOperations)

	assert.False(t, *config.PathItem.AdditionalOperations.Added)
	assert.False(t, *config.PathItem.AdditionalOperations.Modified)
	assert.True(t, *config.PathItem.AdditionalOperations.Removed)
}

func TestDefaultBreakingRules_MediaType_ItemSchema(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.MediaType)
	assert.NotNil(t, config.MediaType.ItemSchema)

	assert.True(t, *config.MediaType.ItemSchema.Added)
	assert.False(t, *config.MediaType.ItemSchema.Modified)
	assert.True(t, *config.MediaType.ItemSchema.Removed)
}

func TestDefaultBreakingRules_MediaType_ItemEncoding(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.MediaType)
	assert.NotNil(t, config.MediaType.ItemEncoding)

	assert.True(t, *config.MediaType.ItemEncoding.Added)
	assert.False(t, *config.MediaType.ItemEncoding.Modified)
	assert.True(t, *config.MediaType.ItemEncoding.Removed)
}

func TestDefaultBreakingRules_SecurityScheme_OAuth2MetadataUrl(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.SecurityScheme)
	assert.NotNil(t, config.SecurityScheme.OAuth2MetadataUrl)

	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Added)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Modified)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Removed)
}

func TestGetRule_OpenAPI32_PathItem_Query(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	rule := config.GetRule("pathItem", "query")
	assert.NotNil(t, rule, "expected rule for pathItem/query")
	assert.False(t, *rule.Added)
	assert.False(t, *rule.Modified)
	assert.True(t, *rule.Removed)
}

func TestGetRule_OpenAPI32_PathItem_AdditionalOperations(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	rule := config.GetRule("pathItem", "additionalOperations")
	assert.NotNil(t, rule, "expected rule for pathItem/additionalOperations")
	assert.False(t, *rule.Added)
	assert.False(t, *rule.Modified)
	assert.True(t, *rule.Removed)
}

func TestGetRule_OpenAPI32_MediaType_ItemSchema(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	rule := config.GetRule("mediaType", "itemSchema")
	assert.NotNil(t, rule, "expected rule for mediaType/itemSchema")
	assert.True(t, *rule.Added)
	assert.False(t, *rule.Modified)
	assert.True(t, *rule.Removed)
}

func TestGetRule_OpenAPI32_MediaType_ItemEncoding(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	rule := config.GetRule("mediaType", "itemEncoding")
	assert.NotNil(t, rule, "expected rule for mediaType/itemEncoding")
	assert.True(t, *rule.Added)
	assert.False(t, *rule.Modified)
	assert.True(t, *rule.Removed)
}

func TestGetRule_OpenAPI32_SecurityScheme_OAuth2MetadataUrl(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	rule := config.GetRule("securityScheme", "oauth2MetadataUrl")
	assert.NotNil(t, rule, "expected rule for securityScheme/oauth2MetadataUrl")
	assert.False(t, *rule.Added)
	assert.False(t, *rule.Modified)
	assert.False(t, *rule.Removed)
}

func TestIsBreaking_OpenAPI32_PathItem_Query(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	assert.False(t, config.IsBreaking("pathItem", "query", ChangeTypeAdded))
	assert.False(t, config.IsBreaking("pathItem", "query", ChangeTypeModified))
	assert.True(t, config.IsBreaking("pathItem", "query", ChangeTypeRemoved))
}

func TestIsBreaking_OpenAPI32_PathItem_AdditionalOperations(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	assert.False(t, config.IsBreaking("pathItem", "additionalOperations", ChangeTypeAdded))
	assert.False(t, config.IsBreaking("pathItem", "additionalOperations", ChangeTypeModified))
	assert.True(t, config.IsBreaking("pathItem", "additionalOperations", ChangeTypeRemoved))
}

func TestIsBreaking_OpenAPI32_MediaType_ItemSchema(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	assert.True(t, config.IsBreaking("mediaType", "itemSchema", ChangeTypeAdded))
	assert.False(t, config.IsBreaking("mediaType", "itemSchema", ChangeTypeModified))
	assert.True(t, config.IsBreaking("mediaType", "itemSchema", ChangeTypeRemoved))
}

func TestIsBreaking_OpenAPI32_MediaType_ItemEncoding(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	assert.True(t, config.IsBreaking("mediaType", "itemEncoding", ChangeTypeAdded))
	assert.False(t, config.IsBreaking("mediaType", "itemEncoding", ChangeTypeModified))
	assert.True(t, config.IsBreaking("mediaType", "itemEncoding", ChangeTypeRemoved))
}

func TestIsBreaking_OpenAPI32_SecurityScheme_OAuth2MetadataUrl(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	assert.False(t, config.IsBreaking("securityScheme", "oauth2MetadataUrl", ChangeTypeAdded))
	assert.False(t, config.IsBreaking("securityScheme", "oauth2MetadataUrl", ChangeTypeModified))
	assert.False(t, config.IsBreaking("securityScheme", "oauth2MetadataUrl", ChangeTypeRemoved))
}

func TestMerge_OpenAPI32_PathItem_Query(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	assert.False(t, *config.PathItem.Query.Added)

	override := &BreakingRulesConfig{
		PathItem: &PathItemRules{
			Query: &BreakingChangeRule{
				Added: boolPtr(true),
			},
		},
	}
	config.Merge(override)

	assert.True(t, *config.PathItem.Query.Added)
	assert.False(t, *config.PathItem.Query.Modified)
	assert.True(t, *config.PathItem.Query.Removed)
}

func TestMerge_OpenAPI32_PathItem_AdditionalOperations(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	assert.False(t, *config.PathItem.AdditionalOperations.Added)

	override := &BreakingRulesConfig{
		PathItem: &PathItemRules{
			AdditionalOperations: &BreakingChangeRule{
				Added: boolPtr(true),
			},
		},
	}
	config.Merge(override)

	assert.True(t, *config.PathItem.AdditionalOperations.Added)
	assert.False(t, *config.PathItem.AdditionalOperations.Modified)
	assert.True(t, *config.PathItem.AdditionalOperations.Removed)
}

func TestMerge_OpenAPI32_MediaType_ItemSchema(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	assert.True(t, *config.MediaType.ItemSchema.Added)

	override := &BreakingRulesConfig{
		MediaType: &MediaTypeRules{
			ItemSchema: &BreakingChangeRule{
				Added: boolPtr(false),
			},
		},
	}
	config.Merge(override)

	assert.False(t, *config.MediaType.ItemSchema.Added)
	assert.False(t, *config.MediaType.ItemSchema.Modified)
	assert.True(t, *config.MediaType.ItemSchema.Removed)
}

func TestMerge_OpenAPI32_MediaType_ItemEncoding(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	assert.True(t, *config.MediaType.ItemEncoding.Added)

	override := &BreakingRulesConfig{
		MediaType: &MediaTypeRules{
			ItemEncoding: &BreakingChangeRule{
				Added: boolPtr(false),
			},
		},
	}
	config.Merge(override)

	assert.False(t, *config.MediaType.ItemEncoding.Added)
	assert.False(t, *config.MediaType.ItemEncoding.Modified)
	assert.True(t, *config.MediaType.ItemEncoding.Removed)
}

func TestMerge_OpenAPI32_SecurityScheme_OAuth2MetadataUrl(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Added)

	override := &BreakingRulesConfig{
		SecurityScheme: &SecuritySchemeRules{
			OAuth2MetadataUrl: &BreakingChangeRule{
				Added:    boolPtr(true),
				Modified: boolPtr(true),
				Removed:  boolPtr(true),
			},
		},
	}
	config.Merge(override)

	assert.True(t, *config.SecurityScheme.OAuth2MetadataUrl.Added)
	assert.True(t, *config.SecurityScheme.OAuth2MetadataUrl.Modified)
	assert.True(t, *config.SecurityScheme.OAuth2MetadataUrl.Removed)
}
