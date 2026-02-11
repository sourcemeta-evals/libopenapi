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

	assert.False(t, *config.MediaType.ItemEncoding.Added)
	assert.False(t, *config.MediaType.ItemEncoding.Modified)
	assert.True(t, *config.MediaType.ItemEncoding.Removed)
}

func TestDefaultBreakingRules_SecurityScheme_OAuth2MetadataURL(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.SecurityScheme)
	assert.NotNil(t, config.SecurityScheme.OAuth2MetadataURL)

	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Added)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Modified)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Removed)
}

func TestGetRule_OpenAPI32_PathItemQuery(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	r := config.GetRule(CompPathItem, PropQuery)
	assert.NotNil(t, r)
	assert.False(t, *r.Added)
	assert.True(t, *r.Removed)
}

func TestGetRule_OpenAPI32_PathItemAdditionalOperations(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	r := config.GetRule(CompPathItem, PropAdditionalOperations)
	assert.NotNil(t, r)
	assert.False(t, *r.Added)
	assert.True(t, *r.Removed)
}

func TestGetRule_OpenAPI32_MediaTypeItemSchema(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	r := config.GetRule(CompMediaType, PropItemSchema)
	assert.NotNil(t, r)
	assert.True(t, *r.Added)
	assert.True(t, *r.Removed)
}

func TestGetRule_OpenAPI32_MediaTypeItemEncoding(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	r := config.GetRule(CompMediaType, PropItemEncoding)
	assert.NotNil(t, r)
	assert.False(t, *r.Added)
	assert.True(t, *r.Removed)
}

func TestGetRule_OpenAPI32_SecuritySchemeOAuth2MetadataURL(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	r := config.GetRule(CompSecurityScheme, PropOAuth2MetadataURL)
	assert.NotNil(t, r)
	assert.False(t, *r.Added)
	assert.False(t, *r.Modified)
	assert.False(t, *r.Removed)
}

func TestIsBreaking_OpenAPI32_PathItemQuery(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	assert.False(t, config.IsBreaking(CompPathItem, PropQuery, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompPathItem, PropQuery, ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompPathItem, PropQuery, ChangeTypeRemoved))
}

func TestIsBreaking_OpenAPI32_PathItemAdditionalOperations(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	assert.False(t, config.IsBreaking(CompPathItem, PropAdditionalOperations, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompPathItem, PropAdditionalOperations, ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompPathItem, PropAdditionalOperations, ChangeTypeRemoved))
}

func TestIsBreaking_OpenAPI32_MediaTypeItemSchema(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	assert.True(t, config.IsBreaking(CompMediaType, PropItemSchema, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompMediaType, PropItemSchema, ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompMediaType, PropItemSchema, ChangeTypeRemoved))
}

func TestIsBreaking_OpenAPI32_MediaTypeItemEncoding(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	assert.False(t, config.IsBreaking(CompMediaType, PropItemEncoding, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompMediaType, PropItemEncoding, ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompMediaType, PropItemEncoding, ChangeTypeRemoved))
}

func TestIsBreaking_OpenAPI32_SecuritySchemeOAuth2MetadataURL(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	assert.False(t, config.IsBreaking(CompSecurityScheme, PropOAuth2MetadataURL, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompSecurityScheme, PropOAuth2MetadataURL, ChangeTypeModified))
	assert.False(t, config.IsBreaking(CompSecurityScheme, PropOAuth2MetadataURL, ChangeTypeRemoved))
}

func TestMerge_OpenAPI32_PathItemQueryOverride(t *testing.T) {
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
	assert.True(t, *config.PathItem.Query.Removed)
}

func TestMerge_OpenAPI32_AdditionalOperationsOverride(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	override := &BreakingRulesConfig{
		PathItem: &PathItemRules{
			AdditionalOperations: &BreakingChangeRule{
				Added:   boolPtr(true),
				Removed: boolPtr(false),
			},
		},
	}
	config.Merge(override)

	assert.True(t, *config.PathItem.AdditionalOperations.Added)
	assert.False(t, *config.PathItem.AdditionalOperations.Removed)
}

func TestMerge_OpenAPI32_MediaTypeItemSchemaOverride(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	override := &BreakingRulesConfig{
		MediaType: &MediaTypeRules{
			ItemSchema: &BreakingChangeRule{
				Added: boolPtr(false),
			},
		},
	}
	config.Merge(override)

	assert.False(t, *config.MediaType.ItemSchema.Added)
	assert.True(t, *config.MediaType.ItemSchema.Removed)
}

func TestMerge_OpenAPI32_MediaTypeItemEncodingOverride(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	override := &BreakingRulesConfig{
		MediaType: &MediaTypeRules{
			ItemEncoding: &BreakingChangeRule{
				Modified: boolPtr(true),
			},
		},
	}
	config.Merge(override)

	assert.True(t, *config.MediaType.ItemEncoding.Modified)
	assert.True(t, *config.MediaType.ItemEncoding.Removed)
}

func TestMerge_OpenAPI32_SecuritySchemeOAuth2MetadataURLOverride(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	override := &BreakingRulesConfig{
		SecurityScheme: &SecuritySchemeRules{
			OAuth2MetadataURL: &BreakingChangeRule{
				Modified: boolPtr(true),
			},
		},
	}
	config.Merge(override)

	assert.True(t, *config.SecurityScheme.OAuth2MetadataURL.Modified)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Added)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Removed)
}
