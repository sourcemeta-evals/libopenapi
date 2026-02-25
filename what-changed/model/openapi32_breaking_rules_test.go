// Copyright 2022-2025 Princess Beef Heavy Industries, LLC / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRule_OpenAPI32Properties(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	tests := []struct {
		component string
		property  string
	}{
		{CompPathItem, "query"},
		{CompPathItem, "additionalOperations"},
		{CompMediaType, "itemSchema"},
		{CompMediaType, "itemEncoding"},
		{CompSecurityScheme, "oauth2MetadataUrl"},
	}

	for _, tt := range tests {
		t.Run(tt.component+"/"+tt.property, func(t *testing.T) {
			assert.NotNil(t, config.GetRule(tt.component, tt.property))
		})
	}
}

func TestDefaultBreakingRules_OpenAPI32PropertyDefaults(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	assert.False(t, config.IsBreaking(CompPathItem, "query", ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompPathItem, "query", ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompPathItem, "query", ChangeTypeRemoved))

	assert.False(t, config.IsBreaking(CompPathItem, "additionalOperations", ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompPathItem, "additionalOperations", ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompPathItem, "additionalOperations", ChangeTypeRemoved))

	assert.True(t, config.IsBreaking(CompMediaType, "itemSchema", ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompMediaType, "itemSchema", ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompMediaType, "itemSchema", ChangeTypeRemoved))

	assert.False(t, config.IsBreaking(CompMediaType, "itemEncoding", ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompMediaType, "itemEncoding", ChangeTypeModified))
	assert.True(t, config.IsBreaking(CompMediaType, "itemEncoding", ChangeTypeRemoved))

	assert.False(t, config.IsBreaking(CompSecurityScheme, "oauth2MetadataUrl", ChangeTypeAdded))
	assert.False(t, config.IsBreaking(CompSecurityScheme, "oauth2MetadataUrl", ChangeTypeModified))
	assert.False(t, config.IsBreaking(CompSecurityScheme, "oauth2MetadataUrl", ChangeTypeRemoved))
}
