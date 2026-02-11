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

	assert.True(t, config.IsBreaking(CompPathItem, PropQuery, ChangeTypeRemoved))
	assert.False(t, config.IsBreaking(CompPathItem, PropQuery, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompPathItem, PropQuery, ChangeTypeModified))

	rule := config.GetRule(CompPathItem, PropQuery)
	assert.NotNil(t, rule)
	assert.False(t, *rule.Added)
	assert.False(t, *rule.Modified)
	assert.True(t, *rule.Removed)
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

	assert.True(t, config.IsBreaking(CompPathItem, PropAdditionalOperations, ChangeTypeRemoved))
	assert.False(t, config.IsBreaking(CompPathItem, PropAdditionalOperations, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompPathItem, PropAdditionalOperations, ChangeTypeModified))

	rule := config.GetRule(CompPathItem, PropAdditionalOperations)
	assert.NotNil(t, rule)
	assert.False(t, *rule.Added)
	assert.False(t, *rule.Modified)
	assert.True(t, *rule.Removed)
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

	assert.True(t, config.IsBreaking(CompMediaType, PropItemSchema, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompMediaType, PropItemSchema, ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompMediaType, PropItemSchema, ChangeTypeRemoved))

	rule := config.GetRule(CompMediaType, PropItemSchema)
	assert.NotNil(t, rule)
	assert.True(t, *rule.Added)
	assert.False(t, *rule.Modified)
	assert.True(t, *rule.Removed)
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

	assert.True(t, config.IsBreaking(CompMediaType, PropItemEncoding, ChangeTypeAdded))
	assert.True(t, config.IsBreaking(CompMediaType, PropItemEncoding, ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompMediaType, PropItemEncoding, ChangeTypeRemoved))

	rule := config.GetRule(CompMediaType, PropItemEncoding)
	assert.NotNil(t, rule)
	assert.True(t, *rule.Added)
	assert.True(t, *rule.Modified)
	assert.True(t, *rule.Removed)
}

func TestDefaultBreakingRules_SecurityScheme_OAuth2MetadataUrl(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.SecurityScheme)
	assert.NotNil(t, config.SecurityScheme.OAuth2MetadataUrl)

	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Added)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Modified)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Removed)

	assert.False(t, config.IsBreaking(CompSecurityScheme, PropOAuth2MetadataUrl, ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompSecurityScheme, PropOAuth2MetadataUrl, ChangeTypeModified))
	assert.False(t, config.IsBreaking(CompSecurityScheme, PropOAuth2MetadataUrl, ChangeTypeRemoved))

	rule := config.GetRule(CompSecurityScheme, PropOAuth2MetadataUrl)
	assert.NotNil(t, rule)
	assert.False(t, *rule.Added)
	assert.False(t, *rule.Modified)
	assert.False(t, *rule.Removed)
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
			OAuth2MetadataUrl: rule(true, true, true),
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

	assert.True(t, *config.SecurityScheme.OAuth2MetadataUrl.Added)
	assert.True(t, *config.SecurityScheme.OAuth2MetadataUrl.Modified)
	assert.True(t, *config.SecurityScheme.OAuth2MetadataUrl.Removed)

	assert.False(t, *config.PathItem.Get.Added)
	assert.True(t, *config.PathItem.Get.Removed)
	assert.True(t, *config.MediaType.Schema.Added)
	assert.True(t, *config.SecurityScheme.Type.Added)
}

func TestIsBreaking_OpenAPI32ChangeTypes(t *testing.T) {
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

func TestActiveConfig_OpenAPI32Properties(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	assert.True(t, BreakingRemoved(CompPathItem, PropQuery))
	assert.True(t, BreakingRemoved(CompPathItem, PropAdditionalOperations))
	assert.True(t, BreakingAdded(CompMediaType, PropItemSchema))
	assert.True(t, BreakingRemoved(CompMediaType, PropItemEncoding))
	assert.False(t, BreakingModified(CompSecurityScheme, PropOAuth2MetadataUrl))
}
