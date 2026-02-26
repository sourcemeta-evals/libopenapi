// Copyright 2022-2025 Princess Beef Heavy Industries, LLC / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultBreakingRules_OpenAPI32Properties(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config)
	assert.NotNil(t, config.PathItem)
	assert.NotNil(t, config.MediaType)
	assert.NotNil(t, config.SecurityScheme)

	assert.False(t, config.IsBreaking("pathItem", "query", ChangeTypeAdded))
	assert.False(t, config.IsBreaking("pathItem", "query", ChangeTypeModified))
	assert.True(t, config.IsBreaking("pathItem", "query", ChangeTypeRemoved))

	assert.False(t, config.IsBreaking("pathItem", "additionalOperations", ChangeTypeAdded))
	assert.False(t, config.IsBreaking("pathItem", "additionalOperations", ChangeTypeModified))
	assert.True(t, config.IsBreaking("pathItem", "additionalOperations", ChangeTypeRemoved))

	assert.True(t, config.IsBreaking("mediaType", "itemSchema", ChangeTypeAdded))
	assert.False(t, config.IsBreaking("mediaType", "itemSchema", ChangeTypeModified))
	assert.True(t, config.IsBreaking("mediaType", "itemSchema", ChangeTypeRemoved))

	assert.True(t, config.IsBreaking("mediaType", "itemEncoding", ChangeTypeAdded))
	assert.False(t, config.IsBreaking("mediaType", "itemEncoding", ChangeTypeModified))
	assert.True(t, config.IsBreaking("mediaType", "itemEncoding", ChangeTypeRemoved))

	assert.False(t, config.IsBreaking("securityScheme", "oauth2MetadataUrl", ChangeTypeAdded))
	assert.False(t, config.IsBreaking("securityScheme", "oauth2MetadataUrl", ChangeTypeModified))
	assert.False(t, config.IsBreaking("securityScheme", "oauth2MetadataUrl", ChangeTypeRemoved))
}

func TestGetRule_OpenAPI32Properties(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	tests := []struct {
		component string
		property  string
	}{
		{component: "pathItem", property: "query"},
		{component: "pathItem", property: "additionalOperations"},
		{component: "mediaType", property: "itemSchema"},
		{component: "mediaType", property: "itemEncoding"},
		{component: "securityScheme", property: "oauth2MetadataUrl"},
	}

	for _, tt := range tests {
		t.Run(tt.component+"/"+tt.property, func(t *testing.T) {
			rule := config.GetRule(tt.component, tt.property)
			assert.NotNil(t, rule, "expected rule for %s/%s", tt.component, tt.property)
		})
	}
}
