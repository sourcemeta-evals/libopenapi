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

func TestDefaultBreakingRules_SecurityScheme_OAuth2MetadataURL(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.SecurityScheme)
	assert.NotNil(t, config.SecurityScheme.OAuth2MetadataURL)

	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Added)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Modified)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Removed)
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

func TestGetRule_OpenAPI32_SecurityScheme_OAuth2MetadataURL(t *testing.T) {
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

func TestIsBreaking_OpenAPI32_SecurityScheme_OAuth2MetadataURL(t *testing.T) {
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

func TestMerge_OpenAPI32_SecurityScheme_OAuth2MetadataURL(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Added)

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

func TestGetRule_OpenAPI32_AllNewProperties(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	tests := []struct {
		component string
		property  string
	}{
		{"pathItem", "query"},
		{"pathItem", "additionalOperations"},
		{"mediaType", "itemSchema"},
		{"mediaType", "itemEncoding"},
		{"securityScheme", "oauth2MetadataUrl"},
	}

	for _, tt := range tests {
		t.Run(tt.component+"/"+tt.property, func(t *testing.T) {
			rule := config.GetRule(tt.component, tt.property)
			assert.NotNil(t, rule, "expected rule for %s/%s", tt.component, tt.property)
		})
	}
}

func TestIsBreaking_OpenAPI32_ChangeTypes(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	tests := []struct {
		component  string
		property   string
		changeType string
		expected   bool
	}{
		{"pathItem", "query", ChangeTypeAdded, false},
		{"pathItem", "query", ChangeTypeModified, false},
		{"pathItem", "query", ChangeTypeRemoved, true},

		{"pathItem", "additionalOperations", ChangeTypeAdded, false},
		{"pathItem", "additionalOperations", ChangeTypeModified, false},
		{"pathItem", "additionalOperations", ChangeTypeRemoved, true},

		{"mediaType", "itemSchema", ChangeTypeAdded, true},
		{"mediaType", "itemSchema", ChangeTypeModified, false},
		{"mediaType", "itemSchema", ChangeTypeRemoved, true},

		{"mediaType", "itemEncoding", ChangeTypeAdded, true},
		{"mediaType", "itemEncoding", ChangeTypeModified, false},
		{"mediaType", "itemEncoding", ChangeTypeRemoved, true},

		{"securityScheme", "oauth2MetadataUrl", ChangeTypeAdded, false},
		{"securityScheme", "oauth2MetadataUrl", ChangeTypeModified, false},
		{"securityScheme", "oauth2MetadataUrl", ChangeTypeRemoved, false},
	}

	for _, tt := range tests {
		name := tt.component + "/" + tt.property + "/" + tt.changeType
		t.Run(name, func(t *testing.T) {
			result := config.IsBreaking(tt.component, tt.property, tt.changeType)
			assert.Equal(t, tt.expected, result, "IsBreaking(%s, %s, %s)", tt.component, tt.property, tt.changeType)
		})
	}
}
