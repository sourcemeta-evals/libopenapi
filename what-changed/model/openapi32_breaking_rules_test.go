package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultBreakingRules_OpenAPI32PathItemProperties(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	if assert.NotNil(t, config.PathItem) {
		if assert.NotNil(t, config.PathItem.Query) {
			assert.False(t, *config.PathItem.Query.Added)
			assert.False(t, *config.PathItem.Query.Modified)
			assert.True(t, *config.PathItem.Query.Removed)
		}
		if assert.NotNil(t, config.PathItem.AdditionalOperations) {
			assert.False(t, *config.PathItem.AdditionalOperations.Added)
			assert.False(t, *config.PathItem.AdditionalOperations.Modified)
			assert.True(t, *config.PathItem.AdditionalOperations.Removed)
		}
	}
}

func TestDefaultBreakingRules_OpenAPI32MediaTypeProperties(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	if assert.NotNil(t, config.MediaType) {
		if assert.NotNil(t, config.MediaType.ItemSchema) {
			assert.True(t, *config.MediaType.ItemSchema.Added)
			assert.False(t, *config.MediaType.ItemSchema.Modified)
			assert.True(t, *config.MediaType.ItemSchema.Removed)
		}
		if assert.NotNil(t, config.MediaType.ItemEncoding) {
			assert.False(t, *config.MediaType.ItemEncoding.Added)
			assert.False(t, *config.MediaType.ItemEncoding.Modified)
			assert.True(t, *config.MediaType.ItemEncoding.Removed)
		}
	}
}

func TestDefaultBreakingRules_OpenAPI32SecuritySchemeProperties(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	if assert.NotNil(t, config.SecurityScheme) && assert.NotNil(t, config.SecurityScheme.OAuth2MetadataURL) {
		assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Added)
		assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Modified)
		assert.False(t, *config.SecurityScheme.OAuth2MetadataURL.Removed)
	}
}

func TestGetRule_OpenAPI32Properties(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	tests := []struct {
		name      string
		component string
		property  string
		rule      *BreakingChangeRule
	}{
		{"pathItem/query", CompPathItem, PropQuery, config.PathItem.Query},
		{"pathItem/additionalOperations", CompPathItem, PropAdditionalOperations, config.PathItem.AdditionalOperations},
		{"mediaType/itemSchema", CompMediaType, PropItemSchema, config.MediaType.ItemSchema},
		{"mediaType/itemEncoding", CompMediaType, PropItemEncoding, config.MediaType.ItemEncoding},
		{"securityScheme/oauth2MetadataUrl", CompSecurityScheme, PropOAuth2MetadataURL, config.SecurityScheme.OAuth2MetadataURL},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Same(t, tt.rule, config.GetRule(tt.component, tt.property))
			assert.NotNil(t, config.GetRule(tt.component, tt.property))
		})
	}
}
