// Copyright 2022-2025 Princess Beef Heavy Industries, LLC / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultBreakingRules_OpenAPI32_PathItemQuery(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config)
	assert.NotNil(t, config.PathItem)

	assert.NotNil(t, config.PathItem.Query)
	assert.False(t, *config.PathItem.Query.Added)
	assert.False(t, *config.PathItem.Query.Modified)
	assert.True(t, *config.PathItem.Query.Removed)

	rule := config.GetRule(CompPathItem, PropQuery)
	assert.Same(t, config.PathItem.Query, rule)
}

func TestDefaultBreakingRules_OpenAPI32_PathItemAdditionalOperations(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config)
	assert.NotNil(t, config.PathItem)

	assert.NotNil(t, config.PathItem.AdditionalOperations)
	assert.False(t, *config.PathItem.AdditionalOperations.Added)
	assert.False(t, *config.PathItem.AdditionalOperations.Modified)
	assert.True(t, *config.PathItem.AdditionalOperations.Removed)

	rule := config.GetRule(CompPathItem, PropAdditionalOperations)
	assert.Same(t, config.PathItem.AdditionalOperations, rule)
}

func TestDefaultBreakingRules_OpenAPI32_MediaTypeItemSchema(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config)
	assert.NotNil(t, config.MediaType)

	assert.NotNil(t, config.MediaType.ItemSchema)
	assert.True(t, *config.MediaType.ItemSchema.Added)
	assert.False(t, *config.MediaType.ItemSchema.Modified)
	assert.True(t, *config.MediaType.ItemSchema.Removed)

	rule := config.GetRule(CompMediaType, PropItemSchema)
	assert.Same(t, config.MediaType.ItemSchema, rule)
}

func TestDefaultBreakingRules_OpenAPI32_MediaTypeItemEncoding(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config)
	assert.NotNil(t, config.MediaType)

	assert.NotNil(t, config.MediaType.ItemEncoding)
	assert.False(t, *config.MediaType.ItemEncoding.Added)
	assert.False(t, *config.MediaType.ItemEncoding.Modified)
	assert.True(t, *config.MediaType.ItemEncoding.Removed)

	rule := config.GetRule(CompMediaType, PropItemEncoding)
	assert.Same(t, config.MediaType.ItemEncoding, rule)
}

func TestDefaultBreakingRules_OpenAPI32_SecuritySchemeOAuth2MetadataURL(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config)
	assert.NotNil(t, config.SecurityScheme)

	assert.NotNil(t, config.SecurityScheme.OAuth2MetadataURL)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Added)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Modified)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Removed)

	rule := config.GetRule(CompSecurityScheme, PropOAuth2MetadataURL)
	assert.Same(t, config.SecurityScheme.OAuth2MetadataURL, rule)
}
