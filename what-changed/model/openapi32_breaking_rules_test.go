// Copyright 2022-2025 Princess Beef Heavy Industries, LLC / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultBreakingRules_PathItem_Query(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.PathItem)
	assert.NotNil(t, config.PathItem.Query)

	assert.False(t, *config.PathItem.Query.Added)
	assert.False(t, *config.PathItem.Query.Modified)
	assert.True(t, *config.PathItem.Query.Removed)
}

func TestDefaultBreakingRules_PathItem_AdditionalOperations(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.PathItem)
	assert.NotNil(t, config.PathItem.AdditionalOperations)

	assert.False(t, *config.PathItem.AdditionalOperations.Added)
	assert.False(t, *config.PathItem.AdditionalOperations.Modified)
	assert.True(t, *config.PathItem.AdditionalOperations.Removed)
}

func TestDefaultBreakingRules_MediaType_ItemSchema(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.MediaType)
	assert.NotNil(t, config.MediaType.ItemSchema)

	assert.True(t, *config.MediaType.ItemSchema.Added)
	assert.False(t, *config.MediaType.ItemSchema.Modified)
	assert.True(t, *config.MediaType.ItemSchema.Removed)
}

func TestDefaultBreakingRules_MediaType_ItemEncoding(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.MediaType)
	assert.NotNil(t, config.MediaType.ItemEncoding)

	assert.True(t, *config.MediaType.ItemEncoding.Added)
	assert.True(t, *config.MediaType.ItemEncoding.Modified)
	assert.True(t, *config.MediaType.ItemEncoding.Removed)
}

func TestDefaultBreakingRules_SecurityScheme_OAuth2MetadataURL(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.SecurityScheme)
	assert.NotNil(t, config.SecurityScheme.OAuth2MetadataURL)

	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Added)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Modified)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Removed)
}

func TestGetRule_OpenAPI32Properties(t *testing.T) {
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
		})
	}
}

func TestIsBreaking_OpenAPI32Properties(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	tests := []struct {
		component  string
		property   string
		changeType string
		expected   bool
	}{
		{CompPathItem, PropQuery, ChangeTypeAdded, false},
		{CompPathItem, PropQuery, ChangeTypeModified, false},
		{CompPathItem, PropQuery, ChangeTypeRemoved, true},

		{CompPathItem, PropAdditionalOperations, ChangeTypeAdded, false},
		{CompPathItem, PropAdditionalOperations, ChangeTypeModified, false},
		{CompPathItem, PropAdditionalOperations, ChangeTypeRemoved, true},

		{CompMediaType, PropItemSchema, ChangeTypeAdded, true},
		{CompMediaType, PropItemSchema, ChangeTypeModified, false},
		{CompMediaType, PropItemSchema, ChangeTypeRemoved, true},

		{CompMediaType, PropItemEncoding, ChangeTypeAdded, true},
		{CompMediaType, PropItemEncoding, ChangeTypeModified, true},
		{CompMediaType, PropItemEncoding, ChangeTypeRemoved, true},

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

func TestMerge_OpenAPI32Properties(t *testing.T) {
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

	assert.True(t, *config.PathItem.Query.Added)
	assert.True(t, *config.PathItem.Query.Modified)
	assert.True(t, *config.PathItem.Query.Removed)

	assert.True(t, *config.PathItem.AdditionalOperations.Added)
	assert.True(t, *config.PathItem.AdditionalOperations.Modified)
	assert.True(t, *config.PathItem.AdditionalOperations.Removed)

	assert.False(t, *config.MediaType.ItemSchema.Added)
	assert.False(t, *config.MediaType.ItemSchema.Modified)
	assert.False(t, *config.MediaType.ItemSchema.Removed)

	assert.False(t, *config.MediaType.ItemEncoding.Added)
	assert.False(t, *config.MediaType.ItemEncoding.Modified)
	assert.False(t, *config.MediaType.ItemEncoding.Removed)

	assert.True(t, *config.SecurityScheme.OAuth2MetadataURL.Added)
	assert.True(t, *config.SecurityScheme.OAuth2MetadataURL.Modified)
	assert.True(t, *config.SecurityScheme.OAuth2MetadataURL.Removed)

	assert.False(t, *config.PathItem.Get.Added)
	assert.True(t, *config.PathItem.Get.Removed)
	assert.True(t, *config.MediaType.Schema.Added)
	assert.True(t, *config.SecurityScheme.Type.Added)
}
