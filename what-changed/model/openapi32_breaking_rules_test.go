// Copyright 2022-2025 Princess Beef Heavy Industries, LLC / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultBreakingRules_OpenAPI32Properties(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	assert.NotNil(t, config.PathItem.Query)
	assert.False(t, *config.PathItem.Query.Added)
	assert.False(t, *config.PathItem.Query.Modified)
	assert.True(t, *config.PathItem.Query.Removed)

	assert.NotNil(t, config.PathItem.AdditionalOperations)
	assert.False(t, *config.PathItem.AdditionalOperations.Added)
	assert.False(t, *config.PathItem.AdditionalOperations.Modified)
	assert.True(t, *config.PathItem.AdditionalOperations.Removed)

	assert.NotNil(t, config.MediaType.ItemSchema)
	assert.True(t, *config.MediaType.ItemSchema.Added)
	assert.False(t, *config.MediaType.ItemSchema.Modified)
	assert.True(t, *config.MediaType.ItemSchema.Removed)

	assert.NotNil(t, config.MediaType.ItemEncoding)
	assert.False(t, *config.MediaType.ItemEncoding.Added)
	assert.False(t, *config.MediaType.ItemEncoding.Modified)
	assert.True(t, *config.MediaType.ItemEncoding.Removed)

	assert.NotNil(t, config.SecurityScheme.OAuth2MetadataURL)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Added)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Modified)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Removed)
}

func TestGetRule_OpenAPI32Properties(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	tests := []struct {
		component string
		property  string
		added     bool
		modified  bool
		removed   bool
	}{
		{CompPathItem, PropQuery, false, false, true},
		{CompPathItem, PropAdditionalOperations, false, false, true},
		{CompMediaType, PropItemSchema, true, false, true},
		{CompMediaType, PropItemEncoding, false, false, true},
		{CompSecurityScheme, PropOAuth2MetadataURL, false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.component+"/"+tt.property, func(t *testing.T) {
			rule := config.GetRule(tt.component, tt.property)
			assert.NotNil(t, rule)
			assert.Equal(t, tt.added, *rule.Added)
			assert.Equal(t, tt.modified, *rule.Modified)
			assert.Equal(t, tt.removed, *rule.Removed)

			assert.Equal(t, tt.added, config.IsBreaking(tt.component, tt.property, ChangeTypeAdded))
			assert.Equal(t, tt.modified, config.IsBreaking(tt.component, tt.property, ChangeTypeModified))
			assert.Equal(t, tt.removed, config.IsBreaking(tt.component, tt.property, ChangeTypeRemoved))
		})
	}
}
