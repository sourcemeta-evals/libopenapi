// Copyright 2022-2025 Princess Beef Heavy Industries, LLC / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenAPI32BreakingRules_PathItemProperties(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.PathItem)

	assert.NotNil(t, config.PathItem.Query)
	assert.False(t, *config.PathItem.Query.Added)
	assert.False(t, *config.PathItem.Query.Modified)
	assert.True(t, *config.PathItem.Query.Removed)

	assert.NotNil(t, config.PathItem.AdditionalOperations)
	assert.False(t, *config.PathItem.AdditionalOperations.Added)
	assert.False(t, *config.PathItem.AdditionalOperations.Modified)
	assert.True(t, *config.PathItem.AdditionalOperations.Removed)

	assert.NotNil(t, config.GetRule(CompPathItem, PropQuery))
	assert.NotNil(t, config.GetRule(CompPathItem, PropAdditionalOperations))

	assert.False(t, config.IsBreaking(CompPathItem, PropQuery, ChangeTypeAdded))
	assert.True(t, config.IsBreaking(CompPathItem, PropQuery, ChangeTypeRemoved))
	assert.False(t, config.IsBreaking(CompPathItem, PropAdditionalOperations, ChangeTypeAdded))
	assert.True(t, config.IsBreaking(CompPathItem, PropAdditionalOperations, ChangeTypeRemoved))
}

func TestOpenAPI32BreakingRules_MediaTypeProperties(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.MediaType)

	assert.NotNil(t, config.MediaType.ItemSchema)
	assert.True(t, *config.MediaType.ItemSchema.Added)
	assert.False(t, *config.MediaType.ItemSchema.Modified)
	assert.True(t, *config.MediaType.ItemSchema.Removed)

	assert.NotNil(t, config.MediaType.ItemEncoding)
	assert.False(t, *config.MediaType.ItemEncoding.Added)
	assert.False(t, *config.MediaType.ItemEncoding.Modified)
	assert.True(t, *config.MediaType.ItemEncoding.Removed)

	assert.NotNil(t, config.GetRule(CompMediaType, PropItemSchema))
	assert.NotNil(t, config.GetRule(CompMediaType, PropItemEncoding))

	assert.True(t, config.IsBreaking(CompMediaType, PropItemSchema, ChangeTypeAdded))
	assert.True(t, config.IsBreaking(CompMediaType, PropItemSchema, ChangeTypeRemoved))
	assert.False(t, config.IsBreaking(CompMediaType, PropItemEncoding, ChangeTypeAdded))
	assert.True(t, config.IsBreaking(CompMediaType, PropItemEncoding, ChangeTypeRemoved))
}

func TestOpenAPI32BreakingRules_SecuritySchemeProperties(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.SecurityScheme)

	assert.NotNil(t, config.SecurityScheme.OAuth2MetadataURL)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Added)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Modified)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Removed)

	rule := config.GetRule(CompSecurityScheme, PropOAuth2MetadataURL)
	assert.NotNil(t, rule)
	assert.Equal(t, config.SecurityScheme.OAuth2MetadataURL, rule)

	assert.False(t, config.IsBreaking(CompSecurityScheme, PropOAuth2MetadataURL, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompSecurityScheme, PropOAuth2MetadataURL, ChangeTypeModified))
	assert.False(t, config.IsBreaking(CompSecurityScheme, PropOAuth2MetadataURL, ChangeTypeRemoved))
}
