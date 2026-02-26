// Copyright 2022-2025 Princess Beef Heavy Industries, LLC / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultBreakingRules_OpenAPI32Properties(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	tests := []struct {
		name             string
		component        string
		property         string
		expectedAdded    bool
		expectedModified bool
		expectedRemoved  bool
	}{
		{
			name:             "pathItem/query",
			component:        CompPathItem,
			property:         PropQuery,
			expectedAdded:    false,
			expectedModified: false,
			expectedRemoved:  true,
		},
		{
			name:             "pathItem/additionalOperations",
			component:        CompPathItem,
			property:         PropAdditionalOperations,
			expectedAdded:    false,
			expectedModified: false,
			expectedRemoved:  true,
		},
		{
			name:             "mediaType/itemSchema",
			component:        CompMediaType,
			property:         PropItemSchema,
			expectedAdded:    true,
			expectedModified: false,
			expectedRemoved:  true,
		},
		{
			name:             "mediaType/itemEncoding",
			component:        CompMediaType,
			property:         PropItemEncoding,
			expectedAdded:    false,
			expectedModified: false,
			expectedRemoved:  true,
		},
		{
			name:             "securityScheme/oauth2MetadataUrl",
			component:        CompSecurityScheme,
			property:         PropOAuth2MetadataURL,
			expectedAdded:    false,
			expectedModified: false,
			expectedRemoved:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := config.GetRule(tt.component, tt.property)
			if assert.NotNil(t, rule) {
				assert.Equal(t, tt.expectedAdded, *rule.Added)
				assert.Equal(t, tt.expectedModified, *rule.Modified)
				assert.Equal(t, tt.expectedRemoved, *rule.Removed)
			}

			assert.Equal(t, tt.expectedAdded, config.IsBreaking(tt.component, tt.property, ChangeTypeAdded))
			assert.Equal(t, tt.expectedModified, config.IsBreaking(tt.component, tt.property, ChangeTypeModified))
			assert.Equal(t, tt.expectedRemoved, config.IsBreaking(tt.component, tt.property, ChangeTypeRemoved))
		})
	}
}

func TestDefaultBreakingRules_OpenAPI32PropertyFieldsPresent(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	assert.NotNil(t, config.PathItem.Query)
	assert.NotNil(t, config.PathItem.AdditionalOperations)
	assert.NotNil(t, config.MediaType.ItemSchema)
	assert.NotNil(t, config.MediaType.ItemEncoding)
	assert.NotNil(t, config.SecurityScheme.OAuth2MetadataURL)
}
