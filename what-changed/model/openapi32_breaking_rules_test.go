// Copyright 2022-2025 Princess Beef Heavy Industries, LLC / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestOpenAPI32_PathItem_Query verifies that the query (OpenAPI 3.2+) property
// is present in PathItemRules with the expected default breaking change values.
func TestOpenAPI32_PathItem_Query(t *testing.T) {
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

	// GetRule lookup should work
	rule := config.GetRule(CompPathItem, PropQuery)
	assert.NotNil(t, rule)
	assert.False(t, *rule.Added)
	assert.False(t, *rule.Modified)
	assert.True(t, *rule.Removed)

	// IsBreaking checks
	assert.False(t, config.IsBreaking(CompPathItem, PropQuery, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompPathItem, PropQuery, ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompPathItem, PropQuery, ChangeTypeRemoved))
}

// TestOpenAPI32_PathItem_AdditionalOperations verifies that the additionalOperations
// (OpenAPI 3.2+) property is present in PathItemRules with expected defaults.
func TestOpenAPI32_PathItem_AdditionalOperations(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.PathItem)
	assert.NotNil(t, config.PathItem.AdditionalOperations)

	// additional operations follow the same pattern as other operations:
	// adding is not breaking, removing is breaking
	assert.False(t, *config.PathItem.AdditionalOperations.Added)
	assert.False(t, *config.PathItem.AdditionalOperations.Modified)
	assert.True(t, *config.PathItem.AdditionalOperations.Removed)

	// GetRule lookup should work
	rule := config.GetRule(CompPathItem, PropAdditionalOperations)
	assert.NotNil(t, rule)
	assert.False(t, *rule.Added)
	assert.False(t, *rule.Modified)
	assert.True(t, *rule.Removed)

	// IsBreaking checks
	assert.False(t, config.IsBreaking(CompPathItem, PropAdditionalOperations, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompPathItem, PropAdditionalOperations, ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompPathItem, PropAdditionalOperations, ChangeTypeRemoved))
}

// TestOpenAPI32_MediaType_ItemSchema verifies that the itemSchema (OpenAPI 3.2+)
// property is present in MediaTypeRules with expected defaults.
func TestOpenAPI32_MediaType_ItemSchema(t *testing.T) {
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

	// GetRule lookup should work
	rule := config.GetRule(CompMediaType, PropItemSchema)
	assert.NotNil(t, rule)
	assert.True(t, *rule.Added)
	assert.False(t, *rule.Modified)
	assert.True(t, *rule.Removed)

	// IsBreaking checks
	assert.True(t, config.IsBreaking(CompMediaType, PropItemSchema, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompMediaType, PropItemSchema, ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompMediaType, PropItemSchema, ChangeTypeRemoved))
}

// TestOpenAPI32_MediaType_ItemEncoding verifies that the itemEncoding (OpenAPI 3.2+)
// property is present in MediaTypeRules with expected defaults.
func TestOpenAPI32_MediaType_ItemEncoding(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.MediaType)
	assert.NotNil(t, config.MediaType.ItemEncoding)

	// itemEncoding changes are all breaking
	assert.True(t, *config.MediaType.ItemEncoding.Added)
	assert.True(t, *config.MediaType.ItemEncoding.Modified)
	assert.True(t, *config.MediaType.ItemEncoding.Removed)

	// GetRule lookup should work
	rule := config.GetRule(CompMediaType, PropItemEncoding)
	assert.NotNil(t, rule)
	assert.True(t, *rule.Added)
	assert.True(t, *rule.Modified)
	assert.True(t, *rule.Removed)

	// IsBreaking checks
	assert.True(t, config.IsBreaking(CompMediaType, PropItemEncoding, ChangeTypeAdded))
	assert.True(t, config.IsBreaking(CompMediaType, PropItemEncoding, ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompMediaType, PropItemEncoding, ChangeTypeRemoved))
}

// TestOpenAPI32_SecurityScheme_OAuth2MetadataUrl verifies that the oauth2MetadataUrl
// (OpenAPI 3.2+) property is present in SecuritySchemeRules with expected defaults.
func TestOpenAPI32_SecurityScheme_OAuth2MetadataUrl(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.SecurityScheme)
	assert.NotNil(t, config.SecurityScheme.OAuth2MetadataURL)

	// oauth2MetadataUrl changes are non-breaking (informational URL)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Added)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Modified)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Removed)

	// GetRule lookup should work
	rule := config.GetRule(CompSecurityScheme, PropOAuth2MetadataURL)
	assert.NotNil(t, rule)
	assert.False(t, *rule.Added)
	assert.False(t, *rule.Modified)
	assert.False(t, *rule.Removed)

	// IsBreaking checks
	assert.False(t, config.IsBreaking(CompSecurityScheme, PropOAuth2MetadataURL, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompSecurityScheme, PropOAuth2MetadataURL, ChangeTypeModified))
	assert.False(t, config.IsBreaking(CompSecurityScheme, PropOAuth2MetadataURL, ChangeTypeRemoved))
}

// TestOpenAPI32_MergeOverride_NewProperties verifies that the new OpenAPI 3.2+
// properties can be overridden via the Merge function.
func TestOpenAPI32_MergeOverride_NewProperties(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	override := &BreakingRulesConfig{
		PathItem: &PathItemRules{
			Query:                rule(true, true, true),
			AdditionalOperations: rule(true, true, true),
		},
		MediaType: &MediaTypeRules{
			ItemSchema:   rule(false, false, false),
			ItemEncoding: rule(false, false, false),
		},
		SecurityScheme: &SecuritySchemeRules{
			OAuth2MetadataURL: rule(true, true, true),
		},
	}
	config.Merge(override)

	// PathItem query should now be all breaking
	assert.True(t, *config.PathItem.Query.Added)
	assert.True(t, *config.PathItem.Query.Modified)
	assert.True(t, *config.PathItem.Query.Removed)

	// PathItem additionalOperations should now be all breaking
	assert.True(t, *config.PathItem.AdditionalOperations.Added)
	assert.True(t, *config.PathItem.AdditionalOperations.Modified)
	assert.True(t, *config.PathItem.AdditionalOperations.Removed)

	// MediaType itemSchema should now be all non-breaking
	assert.False(t, *config.MediaType.ItemSchema.Added)
	assert.False(t, *config.MediaType.ItemSchema.Modified)
	assert.False(t, *config.MediaType.ItemSchema.Removed)

	// MediaType itemEncoding should now be all non-breaking
	assert.False(t, *config.MediaType.ItemEncoding.Added)
	assert.False(t, *config.MediaType.ItemEncoding.Modified)
	assert.False(t, *config.MediaType.ItemEncoding.Removed)

	// SecurityScheme oauth2MetadataUrl should now be all breaking
	assert.True(t, *config.SecurityScheme.OAuth2MetadataURL.Added)
	assert.True(t, *config.SecurityScheme.OAuth2MetadataURL.Modified)
	assert.True(t, *config.SecurityScheme.OAuth2MetadataURL.Removed)
}

// TestOpenAPI32_IsBreaking_NewProperties tests IsBreaking with change type combinations
// for all new OpenAPI 3.2+ properties.
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
		// query
		{CompPathItem, PropQuery, ChangeTypeAdded, false},
		{CompPathItem, PropQuery, ChangeTypeModified, false},
		{CompPathItem, PropQuery, ChangeTypeRemoved, true},

		// additionalOperations
		{CompPathItem, PropAdditionalOperations, ChangeTypeAdded, false},
		{CompPathItem, PropAdditionalOperations, ChangeTypeModified, false},
		{CompPathItem, PropAdditionalOperations, ChangeTypeRemoved, true},

		// itemSchema
		{CompMediaType, PropItemSchema, ChangeTypeAdded, true},
		{CompMediaType, PropItemSchema, ChangeTypeModified, false},
		{CompMediaType, PropItemSchema, ChangeTypeRemoved, true},

		// itemEncoding
		{CompMediaType, PropItemEncoding, ChangeTypeAdded, true},
		{CompMediaType, PropItemEncoding, ChangeTypeModified, true},
		{CompMediaType, PropItemEncoding, ChangeTypeRemoved, true},

		// oauth2MetadataUrl
		{CompSecurityScheme, PropOAuth2MetadataURL, ChangeTypeAdded, false},
		{CompSecurityScheme, PropOAuth2MetadataURL, ChangeTypeModified, false},
		{CompSecurityScheme, PropOAuth2MetadataURL, ChangeTypeRemoved, false},
	}

	for _, tt := range tests {
		name := tt.component + "/" + tt.property + "/" + tt.changeType
		t.Run(name, func(t *testing.T) {
			result := config.IsBreaking(tt.component, tt.property, tt.changeType)
			assert.Equal(t, tt.expected, result, "IsBreaking(%s, %s, %s)", tt.component, tt.property, tt.changeType)
		})
	}
}

// TestOpenAPI32_GetRule_NewProperties verifies that GetRule returns non-nil
// rules for all new OpenAPI 3.2+ properties.
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
		{CompSecurityScheme, PropOAuth2MetadataURL},
	}

	for _, tt := range tests {
		t.Run(tt.component+"/"+tt.property, func(t *testing.T) {
			rule := config.GetRule(tt.component, tt.property)
			assert.NotNil(t, rule, "expected rule for %s/%s", tt.component, tt.property)
			// every default rule should have all three fields set
			assert.NotNil(t, rule.Added, "expected Added for %s/%s", tt.component, tt.property)
			assert.NotNil(t, rule.Modified, "expected Modified for %s/%s", tt.component, tt.property)
			assert.NotNil(t, rule.Removed, "expected Removed for %s/%s", tt.component, tt.property)
		})
	}
}
