// Copyright 2022-2025 Princess Beef Heavy Industries, LLC / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDefaultBreakingRules_PathItem_Query tests the Query property in PathItemRules (OpenAPI 3.2+)
func TestDefaultBreakingRules_PathItem_Query(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.PathItem)
	assert.NotNil(t, config.PathItem.Query)

	// Query operation: adding is fine, removing is breaking (same as other HTTP methods)
	assert.False(t, *config.PathItem.Query.Added)
	assert.False(t, *config.PathItem.Query.Modified)
	assert.True(t, *config.PathItem.Query.Removed)
}

// TestDefaultBreakingRules_PathItem_AdditionalOperations tests the AdditionalOperations property in PathItemRules (OpenAPI 3.2+)
func TestDefaultBreakingRules_PathItem_AdditionalOperations(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.PathItem)
	assert.NotNil(t, config.PathItem.AdditionalOperations)

	// Additional operations: adding is fine, removing is breaking
	assert.False(t, *config.PathItem.AdditionalOperations.Added)
	assert.False(t, *config.PathItem.AdditionalOperations.Modified)
	assert.True(t, *config.PathItem.AdditionalOperations.Removed)
}

// TestDefaultBreakingRules_MediaType_ItemSchema tests the ItemSchema property in MediaTypeRules (OpenAPI 3.2+)
func TestDefaultBreakingRules_MediaType_ItemSchema(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.MediaType)
	assert.NotNil(t, config.MediaType.ItemSchema)

	// ItemSchema: adding is breaking, removing is breaking (same as Schema)
	assert.True(t, *config.MediaType.ItemSchema.Added)
	assert.False(t, *config.MediaType.ItemSchema.Modified)
	assert.True(t, *config.MediaType.ItemSchema.Removed)
}

// TestDefaultBreakingRules_MediaType_ItemEncoding tests the ItemEncoding property in MediaTypeRules (OpenAPI 3.2+)
func TestDefaultBreakingRules_MediaType_ItemEncoding(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.MediaType)
	assert.NotNil(t, config.MediaType.ItemEncoding)

	// ItemEncoding: adding is fine, removing is breaking
	assert.False(t, *config.MediaType.ItemEncoding.Added)
	assert.False(t, *config.MediaType.ItemEncoding.Modified)
	assert.True(t, *config.MediaType.ItemEncoding.Removed)
}

// TestDefaultBreakingRules_SecurityScheme_OAuth2MetadataUrl tests the OAuth2MetadataUrl property in SecuritySchemeRules (OpenAPI 3.2+)
func TestDefaultBreakingRules_SecurityScheme_OAuth2MetadataUrl(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.SecurityScheme)
	assert.NotNil(t, config.SecurityScheme.OAuth2MetadataUrl)

	// OAuth2MetadataUrl: non-breaking for all change types (informational URL)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Added)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Modified)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Removed)
}

// TestIsBreaking_OpenAPI32_PathItem tests IsBreaking for OpenAPI 3.2 PathItem properties
func TestIsBreaking_OpenAPI32_PathItem(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	// Query operation
	assert.False(t, config.IsBreaking(CompPathItem, PropQuery, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompPathItem, PropQuery, ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompPathItem, PropQuery, ChangeTypeRemoved))

	// Additional operations
	assert.False(t, config.IsBreaking(CompPathItem, PropAdditionalOperations, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompPathItem, PropAdditionalOperations, ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompPathItem, PropAdditionalOperations, ChangeTypeRemoved))
}

// TestIsBreaking_OpenAPI32_MediaType tests IsBreaking for OpenAPI 3.2 MediaType properties
func TestIsBreaking_OpenAPI32_MediaType(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	// ItemSchema
	assert.True(t, config.IsBreaking(CompMediaType, PropItemSchema, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompMediaType, PropItemSchema, ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompMediaType, PropItemSchema, ChangeTypeRemoved))

	// ItemEncoding
	assert.False(t, config.IsBreaking(CompMediaType, PropItemEncoding, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompMediaType, PropItemEncoding, ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompMediaType, PropItemEncoding, ChangeTypeRemoved))
}

// TestIsBreaking_OpenAPI32_SecurityScheme tests IsBreaking for OpenAPI 3.2 SecurityScheme properties
func TestIsBreaking_OpenAPI32_SecurityScheme(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	// OAuth2MetadataUrl - all non-breaking
	assert.False(t, config.IsBreaking(CompSecurityScheme, PropOAuth2MetadataUrl, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompSecurityScheme, PropOAuth2MetadataUrl, ChangeTypeModified))
	assert.False(t, config.IsBreaking(CompSecurityScheme, PropOAuth2MetadataUrl, ChangeTypeRemoved))
}

// TestGetRule_OpenAPI32_Properties tests GetRule for all OpenAPI 3.2 properties
func TestGetRule_OpenAPI32_Properties(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	tests := []struct {
		component string
		property  string
	}{
		{CompPathItem, PropQuery},
		{CompPathItem, PropAdditionalOperations},
		{CompMediaType, PropItemSchema},
		{CompMediaType, PropItemEncoding},
		{CompSecurityScheme, PropOAuth2MetadataUrl},
	}

	for _, tt := range tests {
		t.Run(tt.component+"/"+tt.property, func(t *testing.T) {
			rule := config.GetRule(tt.component, tt.property)
			assert.NotNil(t, rule, "expected rule for %s/%s", tt.component, tt.property)
		})
	}
}

// TestMerge_OpenAPI32_PathItemOverride tests merging overrides for OpenAPI 3.2 PathItem properties
func TestMerge_OpenAPI32_PathItemOverride(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	// Verify defaults
	assert.False(t, *config.PathItem.Query.Added)
	assert.True(t, *config.PathItem.Query.Removed)
	assert.False(t, *config.PathItem.AdditionalOperations.Added)
	assert.True(t, *config.PathItem.AdditionalOperations.Removed)

	// Override to make query additions breaking
	override := &BreakingRulesConfig{
		PathItem: &PathItemRules{
			Query: &BreakingChangeRule{
				Added: boolPtr(true),
			},
			AdditionalOperations: &BreakingChangeRule{
				Removed: boolPtr(false),
			},
		},
	}
	config.Merge(override)

	// Query added is now breaking
	assert.True(t, *config.PathItem.Query.Added)
	// Query removed unchanged
	assert.True(t, *config.PathItem.Query.Removed)
	// AdditionalOperations removed is now non-breaking
	assert.False(t, *config.PathItem.AdditionalOperations.Removed)
}

// TestMerge_OpenAPI32_MediaTypeOverride tests merging overrides for OpenAPI 3.2 MediaType properties
func TestMerge_OpenAPI32_MediaTypeOverride(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	// Verify defaults
	assert.True(t, *config.MediaType.ItemSchema.Added)
	assert.True(t, *config.MediaType.ItemSchema.Removed)
	assert.False(t, *config.MediaType.ItemEncoding.Added)
	assert.True(t, *config.MediaType.ItemEncoding.Removed)

	// Override to make itemSchema additions non-breaking
	override := &BreakingRulesConfig{
		MediaType: &MediaTypeRules{
			ItemSchema: &BreakingChangeRule{
				Added: boolPtr(false),
			},
			ItemEncoding: &BreakingChangeRule{
				Modified: boolPtr(true),
			},
		},
	}
	config.Merge(override)

	// ItemSchema added is now non-breaking
	assert.False(t, *config.MediaType.ItemSchema.Added)
	// ItemSchema removed unchanged
	assert.True(t, *config.MediaType.ItemSchema.Removed)
	// ItemEncoding modified is now breaking
	assert.True(t, *config.MediaType.ItemEncoding.Modified)
}

// TestMerge_OpenAPI32_SecuritySchemeOverride tests merging overrides for OpenAPI 3.2 SecurityScheme properties
func TestMerge_OpenAPI32_SecuritySchemeOverride(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	// Verify defaults - all non-breaking
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Added)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Modified)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Removed)

	// Override to make oauth2MetadataUrl changes breaking
	override := &BreakingRulesConfig{
		SecurityScheme: &SecuritySchemeRules{
			OAuth2MetadataUrl: &BreakingChangeRule{
				Modified: boolPtr(true),
				Removed:  boolPtr(true),
			},
		},
	}
	config.Merge(override)

	// OAuth2MetadataUrl modified is now breaking
	assert.True(t, *config.SecurityScheme.OAuth2MetadataUrl.Modified)
	// OAuth2MetadataUrl removed is now breaking
	assert.True(t, *config.SecurityScheme.OAuth2MetadataUrl.Removed)
	// OAuth2MetadataUrl added unchanged
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Added)
}

// TestOpenAPI32_AllNewPropertiesHaveRules ensures all OpenAPI 3.2 properties have complete rules
func TestOpenAPI32_AllNewPropertiesHaveRules(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	// PathItem.Query
	assert.NotNil(t, config.PathItem.Query.Added)
	assert.NotNil(t, config.PathItem.Query.Modified)
	assert.NotNil(t, config.PathItem.Query.Removed)

	// PathItem.AdditionalOperations
	assert.NotNil(t, config.PathItem.AdditionalOperations.Added)
	assert.NotNil(t, config.PathItem.AdditionalOperations.Modified)
	assert.NotNil(t, config.PathItem.AdditionalOperations.Removed)

	// MediaType.ItemSchema
	assert.NotNil(t, config.MediaType.ItemSchema.Added)
	assert.NotNil(t, config.MediaType.ItemSchema.Modified)
	assert.NotNil(t, config.MediaType.ItemSchema.Removed)

	// MediaType.ItemEncoding
	assert.NotNil(t, config.MediaType.ItemEncoding.Added)
	assert.NotNil(t, config.MediaType.ItemEncoding.Modified)
	assert.NotNil(t, config.MediaType.ItemEncoding.Removed)

	// SecurityScheme.OAuth2MetadataUrl
	assert.NotNil(t, config.SecurityScheme.OAuth2MetadataUrl.Added)
	assert.NotNil(t, config.SecurityScheme.OAuth2MetadataUrl.Modified)
	assert.NotNil(t, config.SecurityScheme.OAuth2MetadataUrl.Removed)
}

// TestIsBreaking_OpenAPI32_ChangeTypes tests all change types for OpenAPI 3.2 properties
func TestIsBreaking_OpenAPI32_ChangeTypes(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	tests := []struct {
		component  string
		property   string
		changeType string
		expected   bool
	}{
		// PathItem.Query
		{CompPathItem, PropQuery, ChangeTypeAdded, false},
		{CompPathItem, PropQuery, ChangeTypeModified, false},
		{CompPathItem, PropQuery, ChangeTypeRemoved, true},

		// PathItem.AdditionalOperations
		{CompPathItem, PropAdditionalOperations, ChangeTypeAdded, false},
		{CompPathItem, PropAdditionalOperations, ChangeTypeModified, false},
		{CompPathItem, PropAdditionalOperations, ChangeTypeRemoved, true},

		// MediaType.ItemSchema
		{CompMediaType, PropItemSchema, ChangeTypeAdded, true},
		{CompMediaType, PropItemSchema, ChangeTypeModified, false},
		{CompMediaType, PropItemSchema, ChangeTypeRemoved, true},

		// MediaType.ItemEncoding
		{CompMediaType, PropItemEncoding, ChangeTypeAdded, false},
		{CompMediaType, PropItemEncoding, ChangeTypeModified, false},
		{CompMediaType, PropItemEncoding, ChangeTypeRemoved, true},

		// SecurityScheme.OAuth2MetadataUrl
		{CompSecurityScheme, PropOAuth2MetadataUrl, ChangeTypeAdded, false},
		{CompSecurityScheme, PropOAuth2MetadataUrl, ChangeTypeModified, false},
		{CompSecurityScheme, PropOAuth2MetadataUrl, ChangeTypeRemoved, false},
	}

	for _, tt := range tests {
		name := tt.component + "/" + tt.property + "/" + tt.changeType
		t.Run(name, func(t *testing.T) {
			result := config.IsBreaking(tt.component, tt.property, tt.changeType)
			assert.Equal(t, tt.expected, result, "IsBreaking(%s, %s, %s)", tt.component, tt.property, tt.changeType)
		})
	}
}
