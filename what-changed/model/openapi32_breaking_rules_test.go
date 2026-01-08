// Copyright 2024-2025 Princess Beef Heavy Industries, LLC / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestQueryMethodConfigurableBreakingRules tests that the Query method's breaking change
// behavior can be configured via the breaking rules system.
func TestQueryMethodConfigurableBreakingRules(t *testing.T) {
	ResetDefaultBreakingRules()
	config := GenerateDefaultBreakingRules()

	// Verify Query rule exists in PathItem
	assert.NotNil(t, config.PathItem)
	assert.NotNil(t, config.PathItem.Query)
}

// TestQueryMethodAddedConfigurable tests that adding a Query operation can be configured
// to be breaking or non-breaking.
func TestQueryMethodAddedConfigurable(t *testing.T) {
	ResetDefaultBreakingRules()
	config := GenerateDefaultBreakingRules()

	// Verify Query rule has Added field
	assert.NotNil(t, config.PathItem.Query.Added)
}

// TestAdditionalOperationsConfigurableBreakingRules tests that AdditionalOperations
// breaking change behavior can be configured.
func TestAdditionalOperationsConfigurableBreakingRules(t *testing.T) {
	ResetDefaultBreakingRules()
	config := GenerateDefaultBreakingRules()

	// Verify AdditionalOperations rule exists
	assert.NotNil(t, config.PathItem.AdditionalOperations)
}

// TestItemSchemaConfigurableBreakingRules tests that the ItemSchema field in MediaType
// has configurable breaking change behavior.
func TestItemSchemaConfigurableBreakingRules(t *testing.T) {
	ResetDefaultBreakingRules()
	config := GenerateDefaultBreakingRules()

	// Verify ItemSchema rule exists in MediaType
	assert.NotNil(t, config.MediaType)
	assert.NotNil(t, config.MediaType.ItemSchema)
}

// TestItemSchemaAddedConfigurable tests that adding itemSchema can be configured.
func TestItemSchemaAddedConfigurable(t *testing.T) {
	ResetDefaultBreakingRules()
	config := GenerateDefaultBreakingRules()

	// Verify ItemSchema rule has Added field
	assert.NotNil(t, config.MediaType.ItemSchema.Added)
}

// TestOAuth2MetadataUrlConfigurableBreakingRules tests that the oauth2MetadataUrl field
// in SecurityScheme has configurable breaking change behavior.
func TestOAuth2MetadataUrlConfigurableBreakingRules(t *testing.T) {
	ResetDefaultBreakingRules()
	config := GenerateDefaultBreakingRules()

	// Verify OAuth2MetadataUrl rule exists in SecurityScheme
	assert.NotNil(t, config.SecurityScheme)
	assert.NotNil(t, config.SecurityScheme.OAuth2MetadataUrl)
}

// TestDefaultBreakingRules_Query verifies the default breaking rules for Query method
func TestDefaultBreakingRules_Query(t *testing.T) {
	ResetDefaultBreakingRules()
	config := GenerateDefaultBreakingRules()

	// Query should have same defaults as other HTTP methods
	assert.NotNil(t, config.PathItem.Query)
	assert.False(t, *config.PathItem.Query.Added, "Adding query should not be breaking by default")
	assert.False(t, *config.PathItem.Query.Modified, "Modifying query should not be breaking by default")
	assert.True(t, *config.PathItem.Query.Removed, "Removing query should be breaking by default")
}

// TestDefaultBreakingRules_AdditionalOperations verifies the default breaking rules for AdditionalOperations
func TestDefaultBreakingRules_AdditionalOperations(t *testing.T) {
	ResetDefaultBreakingRules()
	config := GenerateDefaultBreakingRules()

	assert.NotNil(t, config.PathItem.AdditionalOperations)
	assert.False(t, *config.PathItem.AdditionalOperations.Added, "Adding additionalOperations should not be breaking by default")
	assert.False(t, *config.PathItem.AdditionalOperations.Modified, "Modifying additionalOperations should not be breaking by default")
	assert.True(t, *config.PathItem.AdditionalOperations.Removed, "Removing additionalOperations should be breaking by default")
}

// TestDefaultBreakingRules_ItemSchema verifies the default breaking rules for ItemSchema
func TestDefaultBreakingRules_ItemSchema(t *testing.T) {
	ResetDefaultBreakingRules()
	config := GenerateDefaultBreakingRules()

	assert.NotNil(t, config.MediaType.ItemSchema)
	assert.True(t, *config.MediaType.ItemSchema.Added, "Adding itemSchema should be breaking by default")
	assert.False(t, *config.MediaType.ItemSchema.Modified, "Modifying itemSchema should not be breaking by default")
	assert.True(t, *config.MediaType.ItemSchema.Removed, "Removing itemSchema should be breaking by default")
}

// TestDefaultBreakingRules_ItemEncoding verifies the default breaking rules for ItemEncoding
func TestDefaultBreakingRules_ItemEncoding(t *testing.T) {
	ResetDefaultBreakingRules()
	config := GenerateDefaultBreakingRules()

	assert.NotNil(t, config.MediaType.ItemEncoding)
	assert.False(t, *config.MediaType.ItemEncoding.Added, "Adding itemEncoding should not be breaking by default")
	assert.False(t, *config.MediaType.ItemEncoding.Modified, "Modifying itemEncoding should not be breaking by default")
	assert.True(t, *config.MediaType.ItemEncoding.Removed, "Removing itemEncoding should be breaking by default")
}

// TestDefaultBreakingRules_OAuth2MetadataUrl verifies the default breaking rules for OAuth2MetadataUrl
func TestDefaultBreakingRules_OAuth2MetadataUrl(t *testing.T) {
	ResetDefaultBreakingRules()
	config := GenerateDefaultBreakingRules()

	assert.NotNil(t, config.SecurityScheme.OAuth2MetadataUrl)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Added, "Adding oauth2MetadataUrl should not be breaking by default")
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Modified, "Modifying oauth2MetadataUrl should not be breaking by default")
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Removed, "Removing oauth2MetadataUrl should not be breaking by default")
}

// TestBreakingHelperFunctions_OpenAPI32 tests the BreakingAdded/Modified/Removed helper functions
// work correctly with the new OpenAPI 3.2 properties.
func TestBreakingHelperFunctions_OpenAPI32(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	// Verify helper functions exist and can be called
	_ = BreakingAdded(CompPathItem, PropQuery)
	_ = BreakingModified(CompPathItem, PropQuery)
	_ = BreakingRemoved(CompPathItem, PropQuery)

	_ = BreakingAdded(CompPathItem, PropAdditionalOperations)
	_ = BreakingModified(CompPathItem, PropAdditionalOperations)
	_ = BreakingRemoved(CompPathItem, PropAdditionalOperations)

	_ = BreakingAdded(CompMediaType, PropItemSchema)
	_ = BreakingModified(CompMediaType, PropItemSchema)
	_ = BreakingRemoved(CompMediaType, PropItemSchema)

	_ = BreakingAdded(CompMediaType, PropItemEncoding)
	_ = BreakingModified(CompMediaType, PropItemEncoding)
	_ = BreakingRemoved(CompMediaType, PropItemEncoding)

	_ = BreakingAdded(CompSecurityScheme, PropOAuth2MetadataUrl)
	_ = BreakingModified(CompSecurityScheme, PropOAuth2MetadataUrl)
	_ = BreakingRemoved(CompSecurityScheme, PropOAuth2MetadataUrl)
}
