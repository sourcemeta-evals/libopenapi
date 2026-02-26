// Copyright 2022-2025 Princess Beef Heavy Industries, LLC / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestOpenAPI32_PathItem_Query tests that the query property (OpenAPI 3.2+) is properly
// represented in the breaking rules for Path Item objects.
func TestOpenAPI32_PathItem_Query(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.PathItem)
	assert.NotNil(t, config.PathItem.Query)

	// query follows the same pattern as other HTTP methods: adding is not breaking, removing is breaking
	assert.False(t, *config.PathItem.Query.Added)
	assert.False(t, *config.PathItem.Query.Modified)
	assert.True(t, *config.PathItem.Query.Removed)

	// verify via GetRule lookup
	rule := config.GetRule(CompPathItem, PropQuery)
	assert.NotNil(t, rule)
	assert.False(t, *rule.Added)
	assert.False(t, *rule.Modified)
	assert.True(t, *rule.Removed)

	// verify via IsBreaking
	assert.False(t, config.IsBreaking(CompPathItem, PropQuery, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompPathItem, PropQuery, ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompPathItem, PropQuery, ChangeTypeRemoved))
}

// TestOpenAPI32_PathItem_AdditionalOperations tests that the additionalOperations property
// (OpenAPI 3.2+) is properly represented in the breaking rules for Path Item objects.
func TestOpenAPI32_PathItem_AdditionalOperations(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.PathItem)
	assert.NotNil(t, config.PathItem.AdditionalOperations)

	// additionalOperations follows the same pattern: adding is not breaking, removing is breaking
	assert.False(t, *config.PathItem.AdditionalOperations.Added)
	assert.False(t, *config.PathItem.AdditionalOperations.Modified)
	assert.True(t, *config.PathItem.AdditionalOperations.Removed)

	// verify via GetRule lookup
	rule := config.GetRule(CompPathItem, PropAdditionalOperations)
	assert.NotNil(t, rule)
	assert.False(t, *rule.Added)
	assert.False(t, *rule.Modified)
	assert.True(t, *rule.Removed)

	// verify via IsBreaking
	assert.False(t, config.IsBreaking(CompPathItem, PropAdditionalOperations, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompPathItem, PropAdditionalOperations, ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompPathItem, PropAdditionalOperations, ChangeTypeRemoved))
}

// TestOpenAPI32_MediaType_ItemSchema tests that the itemSchema property (OpenAPI 3.2+)
// is properly represented in the breaking rules for Media Type objects.
func TestOpenAPI32_MediaType_ItemSchema(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.MediaType)
	assert.NotNil(t, config.MediaType.ItemSchema)

	// itemSchema follows the same pattern as schema: adding is breaking, removing is breaking
	assert.True(t, *config.MediaType.ItemSchema.Added)
	assert.False(t, *config.MediaType.ItemSchema.Modified)
	assert.True(t, *config.MediaType.ItemSchema.Removed)

	// verify via GetRule lookup
	rule := config.GetRule(CompMediaType, PropItemSchema)
	assert.NotNil(t, rule)
	assert.True(t, *rule.Added)
	assert.False(t, *rule.Modified)
	assert.True(t, *rule.Removed)

	// verify via IsBreaking
	assert.True(t, config.IsBreaking(CompMediaType, PropItemSchema, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompMediaType, PropItemSchema, ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompMediaType, PropItemSchema, ChangeTypeRemoved))
}

// TestOpenAPI32_MediaType_ItemEncoding tests that the itemEncoding property (OpenAPI 3.2+)
// is properly represented in the breaking rules for Media Type objects.
func TestOpenAPI32_MediaType_ItemEncoding(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.MediaType)
	assert.NotNil(t, config.MediaType.ItemEncoding)

	// itemEncoding: removing is breaking
	assert.False(t, *config.MediaType.ItemEncoding.Added)
	assert.False(t, *config.MediaType.ItemEncoding.Modified)
	assert.True(t, *config.MediaType.ItemEncoding.Removed)

	// verify via GetRule lookup
	rule := config.GetRule(CompMediaType, PropItemEncoding)
	assert.NotNil(t, rule)
	assert.False(t, *rule.Added)
	assert.False(t, *rule.Modified)
	assert.True(t, *rule.Removed)

	// verify via IsBreaking
	assert.False(t, config.IsBreaking(CompMediaType, PropItemEncoding, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompMediaType, PropItemEncoding, ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompMediaType, PropItemEncoding, ChangeTypeRemoved))
}

// TestOpenAPI32_SecurityScheme_OAuth2MetadataUrl tests that the oauth2MetadataUrl property
// (OpenAPI 3.2+) is properly represented in the breaking rules for Security Scheme objects.
func TestOpenAPI32_SecurityScheme_OAuth2MetadataUrl(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.SecurityScheme)
	assert.NotNil(t, config.SecurityScheme.OAuth2MetadataUrl)

	// oauth2MetadataUrl is informational and not breaking
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Added)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Modified)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Removed)

	// verify via GetRule lookup
	rule := config.GetRule(CompSecurityScheme, PropOAuth2MetadataUrl)
	assert.NotNil(t, rule)
	assert.False(t, *rule.Added)
	assert.False(t, *rule.Modified)
	assert.False(t, *rule.Removed)

	// verify via IsBreaking
	assert.False(t, config.IsBreaking(CompSecurityScheme, PropOAuth2MetadataUrl, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompSecurityScheme, PropOAuth2MetadataUrl, ChangeTypeModified))
	assert.False(t, config.IsBreaking(CompSecurityScheme, PropOAuth2MetadataUrl, ChangeTypeRemoved))
}

// TestOpenAPI32_MergeOverride_NewProperties tests that the new OpenAPI 3.2 properties
// can be overridden via the Merge mechanism.
func TestOpenAPI32_MergeOverride_NewProperties(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	// override all new properties
	override := &BreakingRulesConfig{
		PathItem: &PathItemRules{
			Query:                rule(true, true, false),
			AdditionalOperations: rule(true, true, false),
		},
		MediaType: &MediaTypeRules{
			ItemSchema:   rule(false, false, false),
			ItemEncoding: rule(true, true, true),
		},
		SecurityScheme: &SecuritySchemeRules{
			OAuth2MetadataUrl: rule(true, true, true),
		},
	}
	config.Merge(override)

	// verify PathItem overrides applied
	assert.True(t, *config.PathItem.Query.Added)
	assert.True(t, *config.PathItem.Query.Modified)
	assert.False(t, *config.PathItem.Query.Removed)

	assert.True(t, *config.PathItem.AdditionalOperations.Added)
	assert.True(t, *config.PathItem.AdditionalOperations.Modified)
	assert.False(t, *config.PathItem.AdditionalOperations.Removed)

	// verify MediaType overrides applied
	assert.False(t, *config.MediaType.ItemSchema.Added)
	assert.False(t, *config.MediaType.ItemSchema.Modified)
	assert.False(t, *config.MediaType.ItemSchema.Removed)

	assert.True(t, *config.MediaType.ItemEncoding.Added)
	assert.True(t, *config.MediaType.ItemEncoding.Modified)
	assert.True(t, *config.MediaType.ItemEncoding.Removed)

	// verify SecurityScheme override applied
	assert.True(t, *config.SecurityScheme.OAuth2MetadataUrl.Added)
	assert.True(t, *config.SecurityScheme.OAuth2MetadataUrl.Modified)
	assert.True(t, *config.SecurityScheme.OAuth2MetadataUrl.Removed)

	// verify existing properties remain unchanged
	assert.False(t, *config.PathItem.Get.Added)
	assert.True(t, *config.PathItem.Get.Removed)
	assert.True(t, *config.MediaType.Schema.Added)
	assert.True(t, *config.SecurityScheme.Type.Added)
}

// TestOpenAPI32_GetRule_NewProperties tests that GetRule returns the correct rules
// for all new OpenAPI 3.2 properties.
func TestOpenAPI32_GetRule_NewProperties(t *testing.T) {
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
			assert.NotNil(t, rule.Added, "expected Added for %s/%s", tt.component, tt.property)
			assert.NotNil(t, rule.Modified, "expected Modified for %s/%s", tt.component, tt.property)
			assert.NotNil(t, rule.Removed, "expected Removed for %s/%s", tt.component, tt.property)
		})
	}
}

// TestOpenAPI32_IsBreaking_NewProperties tests IsBreaking for the new OpenAPI 3.2
// properties across all change types.
func TestOpenAPI32_IsBreaking_NewProperties(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	tests := []struct {
		component  string
		property   string
		changeType string
		expected   bool
	}{
		// query - like other HTTP methods
		{CompPathItem, PropQuery, ChangeTypeAdded, false},
		{CompPathItem, PropQuery, ChangeTypeModified, false},
		{CompPathItem, PropQuery, ChangeTypeRemoved, true},

		// additionalOperations - like other HTTP methods
		{CompPathItem, PropAdditionalOperations, ChangeTypeAdded, false},
		{CompPathItem, PropAdditionalOperations, ChangeTypeModified, false},
		{CompPathItem, PropAdditionalOperations, ChangeTypeRemoved, true},

		// itemSchema - like schema
		{CompMediaType, PropItemSchema, ChangeTypeAdded, true},
		{CompMediaType, PropItemSchema, ChangeTypeModified, false},
		{CompMediaType, PropItemSchema, ChangeTypeRemoved, true},

		// itemEncoding - removing is breaking
		{CompMediaType, PropItemEncoding, ChangeTypeAdded, false},
		{CompMediaType, PropItemEncoding, ChangeTypeModified, false},
		{CompMediaType, PropItemEncoding, ChangeTypeRemoved, true},

		// oauth2MetadataUrl - not breaking
		{CompSecurityScheme, PropOAuth2MetadataUrl, ChangeTypeAdded, false},
		{CompSecurityScheme, PropOAuth2MetadataUrl, ChangeTypeModified, false},
		{CompSecurityScheme, PropOAuth2MetadataUrl, ChangeTypeRemoved, false},
	}

	for _, tt := range tests {
		name := tt.component + "/" + tt.property + "/" + tt.changeType
		t.Run(name, func(t *testing.T) {
			result := config.IsBreaking(tt.component, tt.property, tt.changeType)
			assert.Equal(t, tt.expected, result,
				"IsBreaking(%s, %s, %s) = %v, want %v",
				tt.component, tt.property, tt.changeType, result, tt.expected)
		})
	}
}
