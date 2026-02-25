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
		{"pathItem", "query"},
		{"pathItem", "additionalOperations"},
		{"mediaType", "itemSchema"},
		{"mediaType", "itemEncoding"},
		{"securityScheme", "oauth2MetadataUrl"},
	}

	for _, tt := range tests {
		t.Run(tt.component+"/"+tt.property, func(t *testing.T) {
			rule := config.GetRule(tt.component, tt.property)
			assert.NotNil(t, rule, "expected breaking rule for %s/%s", tt.component, tt.property)
		})
	}
}

func TestIsBreaking_OpenAPI32Properties(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	tests := []struct {
		component  string
		property   string
		changeType string
		expected   bool
	}{
		// Path Item query operation behaves like other HTTP operations.
		{"pathItem", "query", ChangeTypeAdded, false},
		{"pathItem", "query", ChangeTypeModified, false},
		{"pathItem", "query", ChangeTypeRemoved, true},

		// Path Item additionalOperations behaves like other operation collections.
		{"pathItem", "additionalOperations", ChangeTypeAdded, false},
		{"pathItem", "additionalOperations", ChangeTypeModified, false},
		{"pathItem", "additionalOperations", ChangeTypeRemoved, true},

		// Media Type itemSchema behaves like schema.
		{"mediaType", "itemSchema", ChangeTypeAdded, true},
		{"mediaType", "itemSchema", ChangeTypeModified, false},
		{"mediaType", "itemSchema", ChangeTypeRemoved, true},

		// Media Type itemEncoding behaves like encoding maps.
		{"mediaType", "itemEncoding", ChangeTypeAdded, false},
		{"mediaType", "itemEncoding", ChangeTypeModified, false},
		{"mediaType", "itemEncoding", ChangeTypeRemoved, true},

		// Security Scheme oauth2MetadataUrl behaves like other informational URLs.
		{"securityScheme", "oauth2MetadataUrl", ChangeTypeAdded, false},
		{"securityScheme", "oauth2MetadataUrl", ChangeTypeModified, false},
		{"securityScheme", "oauth2MetadataUrl", ChangeTypeRemoved, false},
	}

	for _, tt := range tests {
		name := tt.component + "/" + tt.property + "/" + tt.changeType
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, config.IsBreaking(tt.component, tt.property, tt.changeType))
		})
	}
}
