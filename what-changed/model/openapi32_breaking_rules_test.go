// Copyright 2024-2025 Princess Beef Heavy Industries, LLC / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"context"
	"testing"

	"github.com/pb33f/libopenapi/datamodel/low"
	v3 "github.com/pb33f/libopenapi/datamodel/low/v3"
	"github.com/pb33f/libopenapi/index"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v4"
)

// TestQueryMethodConfigurableBreakingRules tests that the Query method's breaking change
// behavior can be configured via the breaking rules system.
func TestQueryMethodConfigurableBreakingRules(t *testing.T) {
	// ensure clean state
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	// Test that removing a Query operation is breaking by default
	left := `query:
  summary: Query resources
  operationId: queryResources`

	right := `get:
  summary: Get resources`

	var lNode, rNode yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(left), &lNode), "left YAML fixture must parse")
	require.NoError(t, yaml.Unmarshal([]byte(right), &rNode), "right YAML fixture must parse")

	lIdx := index.NewSpecIndexWithConfig(&lNode, index.CreateOpenAPIIndexConfig())
	rIdx := index.NewSpecIndexWithConfig(&rNode, index.CreateOpenAPIIndexConfig())
	ctx := context.Background()

	var lPath, rPath v3.PathItem
	require.NoError(t, low.BuildModel(&lNode, &lPath), "l model BuildModel must succeed")
	require.NoError(t, low.BuildModel(&rNode, &rPath), "r model BuildModel must succeed")

	require.NoError(t, lPath.Build(ctx, nil, lNode.Content[0], lIdx), "lPath Build must succeed")
	require.NoError(t, rPath.Build(ctx, nil, rNode.Content[0], rIdx), "rPath Build must succeed")

	// Default behavior: removing Query should be breaking
	changes := ComparePathItems(&lPath, &rPath)
	assert.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalBreakingChanges(), "Removing query should be breaking by default")

	// Now configure Query removal as non-breaking
	customConfig := &BreakingRulesConfig{
		PathItem: &PathItemRules{
			Query: &BreakingChangeRule{
				Added:    boolPtr(false),
				Modified: boolPtr(false),
				Removed:  boolPtr(false), // Override: removal is not breaking
			},
		},
	}
	SetActiveBreakingRulesConfig(customConfig)

	// Re-run comparison with custom config
	changes2 := ComparePathItems(&lPath, &rPath)
	assert.NotNil(t, changes2)
	assert.Equal(t, 0, changes2.TotalBreakingChanges(), "With custom config, removing query should not be breaking")
}

// TestQueryMethodAddedConfigurable tests that adding a Query operation can be configured
// to be breaking or non-breaking.
func TestQueryMethodAddedConfigurable(t *testing.T) {
	// ensure clean state
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `get:
  summary: Get resources`

	right := `get:
  summary: Get resources
query:
  summary: Query resources
  operationId: queryResources`

	var lNode, rNode yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(left), &lNode), "left YAML fixture must parse")
	require.NoError(t, yaml.Unmarshal([]byte(right), &rNode), "right YAML fixture must parse")

	lIdx := index.NewSpecIndexWithConfig(&lNode, index.CreateOpenAPIIndexConfig())
	rIdx := index.NewSpecIndexWithConfig(&rNode, index.CreateOpenAPIIndexConfig())
	ctx := context.Background()

	var lPath, rPath v3.PathItem
	require.NoError(t, low.BuildModel(&lNode, &lPath), "l model BuildModel must succeed")
	require.NoError(t, low.BuildModel(&rNode, &rPath), "r model BuildModel must succeed")

	require.NoError(t, lPath.Build(ctx, nil, lNode.Content[0], lIdx), "lPath Build must succeed")
	require.NoError(t, rPath.Build(ctx, nil, rNode.Content[0], rIdx), "rPath Build must succeed")

	// Default behavior: adding Query should not be breaking
	changes := ComparePathItems(&lPath, &rPath)
	assert.NotNil(t, changes)
	assert.Equal(t, 0, changes.TotalBreakingChanges(), "Adding query should not be breaking by default")

	// Now configure Query addition as breaking
	customConfig := &BreakingRulesConfig{
		PathItem: &PathItemRules{
			Query: &BreakingChangeRule{
				Added:    boolPtr(true), // Override: addition is breaking
				Modified: boolPtr(false),
				Removed:  boolPtr(true),
			},
		},
	}
	SetActiveBreakingRulesConfig(customConfig)

	// Re-run comparison with custom config
	changes2 := ComparePathItems(&lPath, &rPath)
	assert.NotNil(t, changes2)
	assert.Equal(t, 1, changes2.TotalBreakingChanges(), "With custom config, adding query should be breaking")
}

// TestAdditionalOperationsConfigurableBreakingRules tests that AdditionalOperations
// breaking change behavior can be configured.
func TestAdditionalOperationsConfigurableBreakingRules(t *testing.T) {
	// ensure clean state
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	// Test with additionalOperations present then removed
	left := `additionalOperations:
  CUSTOM:
    summary: Custom operation`

	right := `get:
  summary: Get resources`

	var lNode, rNode yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(left), &lNode), "left YAML fixture must parse")
	require.NoError(t, yaml.Unmarshal([]byte(right), &rNode), "right YAML fixture must parse")

	lIdx := index.NewSpecIndexWithConfig(&lNode, index.CreateOpenAPIIndexConfig())
	rIdx := index.NewSpecIndexWithConfig(&rNode, index.CreateOpenAPIIndexConfig())
	ctx := context.Background()

	var lPath, rPath v3.PathItem
	require.NoError(t, low.BuildModel(&lNode, &lPath), "l model BuildModel must succeed")
	require.NoError(t, low.BuildModel(&rNode, &rPath), "r model BuildModel must succeed")

	require.NoError(t, lPath.Build(ctx, nil, lNode.Content[0], lIdx), "lPath Build must succeed")
	require.NoError(t, rPath.Build(ctx, nil, rNode.Content[0], rIdx), "rPath Build must succeed")

	// Default behavior: removing additionalOperations should be breaking
	changes := ComparePathItems(&lPath, &rPath)
	assert.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalBreakingChanges(), "Removing additionalOperations should be breaking by default")

	// Now configure as non-breaking
	customConfig := &BreakingRulesConfig{
		PathItem: &PathItemRules{
			AdditionalOperations: &BreakingChangeRule{
				Added:    boolPtr(false),
				Modified: boolPtr(false),
				Removed:  boolPtr(false), // Override: removal is not breaking
			},
		},
	}
	SetActiveBreakingRulesConfig(customConfig)

	// Re-run comparison with custom config
	changes2 := ComparePathItems(&lPath, &rPath)
	assert.NotNil(t, changes2)
	assert.Equal(t, 0, changes2.TotalBreakingChanges(), "With custom config, removing additionalOperations should not be breaking")
}

// TestAdditionalOperationsAddedConfigurable verifies that ComparePathItems honours
// SetActiveBreakingRulesConfig for PathItem.AdditionalOperations additions. The
// scenario is deliberately isolated to an addition (right introduces the map,
// left has none) so that flipping Added to true changes the total breaking-change
// count, discriminating a correct BreakingAdded(CompPathItem, PropAdditionalOperations)
// wiring from an implementation that hard-codes the added flag to false.
func TestAdditionalOperationsAddedConfigurable(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	// left has no additionalOperations, right introduces one.
	left := `get:
  summary: Get resources`

	right := `get:
  summary: Get resources
additionalOperations:
  LINK:
    summary: Link operation`

	var lNode, rNode yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(left), &lNode), "left YAML fixture must parse")
	require.NoError(t, yaml.Unmarshal([]byte(right), &rNode), "right YAML fixture must parse")

	lIdx := index.NewSpecIndexWithConfig(&lNode, index.CreateOpenAPIIndexConfig())
	rIdx := index.NewSpecIndexWithConfig(&rNode, index.CreateOpenAPIIndexConfig())
	ctx := context.Background()

	var lPath, rPath v3.PathItem
	require.NoError(t, low.BuildModel(&lNode, &lPath), "l model BuildModel must succeed")
	require.NoError(t, low.BuildModel(&rNode, &rPath), "r model BuildModel must succeed")

	require.NoError(t, lPath.Build(ctx, nil, lNode.Content[0], lIdx), "lPath Build must succeed")
	require.NoError(t, rPath.Build(ctx, nil, rNode.Content[0], rIdx), "rPath Build must succeed")

	// Default behavior: adding additionalOperations should not be breaking.
	changes := ComparePathItems(&lPath, &rPath)
	assert.NotNil(t, changes)
	assert.Equal(t, 0, changes.TotalBreakingChanges(),
		"default additionalOperations rules should classify the addition as non-breaking")

	// Override Added=true: the addition should now be breaking.
	SetActiveBreakingRulesConfig(&BreakingRulesConfig{
		PathItem: &PathItemRules{
			AdditionalOperations: &BreakingChangeRule{
				Added:    boolPtr(true),
				Modified: boolPtr(false),
				Removed:  boolPtr(true),
			},
		},
	})

	changes2 := ComparePathItems(&lPath, &rPath)
	assert.NotNil(t, changes2)
	assert.Equal(t, 1, changes2.TotalBreakingChanges(),
		"overridden additionalOperations rules should classify the addition as breaking")
}

// TestItemSchemaConfigurableBreakingRules tests that the ItemSchema field in MediaType
// has configurable breaking change behavior.
func TestItemSchemaConfigurableBreakingRules(t *testing.T) {
	// ensure clean state
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	// Test with itemSchema present then removed
	left := `schema:
  type: array
itemSchema:
  type: string`

	right := `schema:
  type: array`

	var lNode, rNode yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(left), &lNode), "left YAML fixture must parse")
	require.NoError(t, yaml.Unmarshal([]byte(right), &rNode), "right YAML fixture must parse")

	lIdx := index.NewSpecIndexWithConfig(&lNode, index.CreateOpenAPIIndexConfig())
	rIdx := index.NewSpecIndexWithConfig(&rNode, index.CreateOpenAPIIndexConfig())
	ctx := context.Background()

	var lMT, rMT v3.MediaType
	require.NoError(t, low.BuildModel(&lNode, &lMT), "l model BuildModel must succeed")
	require.NoError(t, low.BuildModel(&rNode, &rMT), "r model BuildModel must succeed")

	require.NoError(t, lMT.Build(ctx, nil, lNode.Content[0], lIdx), "lMT Build must succeed")
	require.NoError(t, rMT.Build(ctx, nil, rNode.Content[0], rIdx), "rMT Build must succeed")

	// Default behavior: removing itemSchema should be breaking
	changes := CompareMediaTypes(&lMT, &rMT)
	assert.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalBreakingChanges(), "Removing itemSchema should be breaking by default")

	// Now configure as non-breaking
	customConfig := &BreakingRulesConfig{
		MediaType: &MediaTypeRules{
			ItemSchema: &BreakingChangeRule{
				Added:    boolPtr(false),
				Modified: boolPtr(false),
				Removed:  boolPtr(false), // Override: removal is not breaking
			},
		},
	}
	SetActiveBreakingRulesConfig(customConfig)

	// Re-run comparison with custom config
	changes2 := CompareMediaTypes(&lMT, &rMT)
	assert.NotNil(t, changes2)
	assert.Equal(t, 0, changes2.TotalBreakingChanges(), "With custom config, removing itemSchema should not be breaking")
}

// TestItemSchemaAddedConfigurable tests that adding itemSchema can be configured.
func TestItemSchemaAddedConfigurable(t *testing.T) {
	// ensure clean state
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `schema:
  type: array`

	right := `schema:
  type: array
itemSchema:
  type: string`

	var lNode, rNode yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(left), &lNode), "left YAML fixture must parse")
	require.NoError(t, yaml.Unmarshal([]byte(right), &rNode), "right YAML fixture must parse")

	lIdx := index.NewSpecIndexWithConfig(&lNode, index.CreateOpenAPIIndexConfig())
	rIdx := index.NewSpecIndexWithConfig(&rNode, index.CreateOpenAPIIndexConfig())
	ctx := context.Background()

	var lMT, rMT v3.MediaType
	require.NoError(t, low.BuildModel(&lNode, &lMT), "l model BuildModel must succeed")
	require.NoError(t, low.BuildModel(&rNode, &rMT), "r model BuildModel must succeed")

	require.NoError(t, lMT.Build(ctx, nil, lNode.Content[0], lIdx), "lMT Build must succeed")
	require.NoError(t, rMT.Build(ctx, nil, rNode.Content[0], rIdx), "rMT Build must succeed")

	// Default behavior: adding itemSchema should be breaking
	changes := CompareMediaTypes(&lMT, &rMT)
	assert.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalBreakingChanges(), "Adding itemSchema should be breaking by default")

	// Now configure as non-breaking
	customConfig := &BreakingRulesConfig{
		MediaType: &MediaTypeRules{
			ItemSchema: &BreakingChangeRule{
				Added:    boolPtr(false), // Override: addition is not breaking
				Modified: boolPtr(false),
				Removed:  boolPtr(true),
			},
		},
	}
	SetActiveBreakingRulesConfig(customConfig)

	// Re-run comparison with custom config
	changes2 := CompareMediaTypes(&lMT, &rMT)
	assert.NotNil(t, changes2)
	assert.Equal(t, 0, changes2.TotalBreakingChanges(), "With custom config, adding itemSchema should not be breaking")
}

// TestOAuth2MetadataUrlConfigurableBreakingRules tests that the oauth2MetadataUrl field
// in SecurityScheme has configurable breaking change behavior.
func TestOAuth2MetadataUrlConfigurableBreakingRules(t *testing.T) {
	// ensure clean state
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	// Include type field which is required for the comparison to work properly
	left := `type: oauth2
description: OAuth2 auth
oauth2MetadataUrl: https://example.com/.well-known/oauth-authorization-server`

	right := `type: oauth2
description: OAuth2 auth
oauth2MetadataUrl: https://example.com/v2/.well-known/oauth-authorization-server`

	var lNode, rNode yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(left), &lNode), "left YAML fixture must parse")
	require.NoError(t, yaml.Unmarshal([]byte(right), &rNode), "right YAML fixture must parse")

	ctx := context.Background()

	var lSS, rSS v3.SecurityScheme
	require.NoError(t, low.BuildModel(lNode.Content[0], &lSS), "l model BuildModel must succeed")
	require.NoError(t, low.BuildModel(rNode.Content[0], &rSS), "r model BuildModel must succeed")

	require.NoError(t, lSS.Build(ctx, nil, lNode.Content[0], nil), "lSS Build must succeed")
	require.NoError(t, rSS.Build(ctx, nil, rNode.Content[0], nil), "rSS Build must succeed")

	// Default behavior: modifying oauth2MetadataUrl should not be breaking
	changes := CompareSecuritySchemes(&lSS, &rSS)
	assert.NotNil(t, changes, "Changes should not be nil when oauth2MetadataUrl differs")
	assert.Equal(t, 1, changes.TotalChanges(), "Should detect 1 change for oauth2MetadataUrl modification")
	assert.Equal(t, 0, changes.TotalBreakingChanges(), "Modifying oauth2MetadataUrl should not be breaking by default")

	// Now configure as breaking
	customConfig := &BreakingRulesConfig{
		SecurityScheme: &SecuritySchemeRules{
			OAuth2MetadataUrl: &BreakingChangeRule{
				Added:    boolPtr(true),
				Modified: boolPtr(true), // Override: modification is breaking
				Removed:  boolPtr(true),
			},
		},
	}
	SetActiveBreakingRulesConfig(customConfig)

	// Re-run comparison with custom config
	changes2 := CompareSecuritySchemes(&lSS, &rSS)
	assert.NotNil(t, changes2, "Changes should not be nil with custom config")
	assert.Equal(t, 1, changes2.TotalBreakingChanges(), "With custom config, modifying oauth2MetadataUrl should be breaking")
}

// TestDefaultBreakingRules_Query verifies the default breaking rules for Query method
func TestDefaultBreakingRules_Query(t *testing.T) {
	ResetDefaultBreakingRules()
	config := GenerateDefaultBreakingRules()

	// Query should have same defaults as other HTTP methods
	assert.NotNil(t, config.PathItem.Query)
	assert.False(t, *config.PathItem.Query.Added, "Adding query should not be breaking by default")
	assert.False(t, *config.PathItem.Query.Modified, "Modifying query should not be breaking by default")
	assert.True(t, *config.PathItem.Query.Removed, "Removing query should be breaking by default")
}

// TestDefaultBreakingRules_AdditionalOperations verifies the default breaking rules for AdditionalOperations
func TestDefaultBreakingRules_AdditionalOperations(t *testing.T) {
	ResetDefaultBreakingRules()
	config := GenerateDefaultBreakingRules()

	assert.NotNil(t, config.PathItem.AdditionalOperations)
	assert.False(t, *config.PathItem.AdditionalOperations.Added, "Adding additionalOperations should not be breaking by default")
	assert.False(t, *config.PathItem.AdditionalOperations.Modified, "Modifying additionalOperations should not be breaking by default")
	assert.True(t, *config.PathItem.AdditionalOperations.Removed, "Removing additionalOperations should be breaking by default")
}

// TestDefaultBreakingRules_ItemSchema verifies the default breaking rules for ItemSchema
func TestDefaultBreakingRules_ItemSchema(t *testing.T) {
	ResetDefaultBreakingRules()
	config := GenerateDefaultBreakingRules()

	assert.NotNil(t, config.MediaType.ItemSchema)
	assert.True(t, *config.MediaType.ItemSchema.Added, "Adding itemSchema should be breaking by default")
	assert.False(t, *config.MediaType.ItemSchema.Modified, "Modifying itemSchema should not be breaking by default")
	assert.True(t, *config.MediaType.ItemSchema.Removed, "Removing itemSchema should be breaking by default")
}

// TestDefaultBreakingRules_ItemEncoding verifies the default breaking rules for ItemEncoding
func TestDefaultBreakingRules_ItemEncoding(t *testing.T) {
	ResetDefaultBreakingRules()
	config := GenerateDefaultBreakingRules()

	assert.NotNil(t, config.MediaType.ItemEncoding)
	assert.False(t, *config.MediaType.ItemEncoding.Added, "Adding itemEncoding should not be breaking by default")
	assert.False(t, *config.MediaType.ItemEncoding.Modified, "Modifying itemEncoding should not be breaking by default")
	assert.True(t, *config.MediaType.ItemEncoding.Removed, "Removing itemEncoding should be breaking by default")
}

// TestDefaultBreakingRules_OAuth2MetadataUrl verifies the default breaking rules for OAuth2MetadataUrl
func TestDefaultBreakingRules_OAuth2MetadataUrl(t *testing.T) {
	ResetDefaultBreakingRules()
	config := GenerateDefaultBreakingRules()

	assert.NotNil(t, config.SecurityScheme.OAuth2MetadataUrl)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Added, "Adding oauth2MetadataUrl should not be breaking by default")
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Modified, "Modifying oauth2MetadataUrl should not be breaking by default")
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Removed, "Removing oauth2MetadataUrl should not be breaking by default")
}

// TestBreakingHelperFunctions_OpenAPI32 tests the BreakingAdded/Modified/Removed helper functions
// work correctly with the new OpenAPI 3.2 properties.
func TestBreakingHelperFunctions_OpenAPI32(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	// Test Query
	assert.False(t, BreakingAdded(CompPathItem, "query"))
	assert.False(t, BreakingModified(CompPathItem, "query"))
	assert.True(t, BreakingRemoved(CompPathItem, "query"))

	// Test AdditionalOperations
	assert.False(t, BreakingAdded(CompPathItem, "additionalOperations"))
	assert.False(t, BreakingModified(CompPathItem, "additionalOperations"))
	assert.True(t, BreakingRemoved(CompPathItem, "additionalOperations"))

	// Test ItemSchema
	assert.True(t, BreakingAdded(CompMediaType, "itemSchema"))
	assert.False(t, BreakingModified(CompMediaType, "itemSchema"))
	assert.True(t, BreakingRemoved(CompMediaType, "itemSchema"))

	// Test ItemEncoding
	assert.False(t, BreakingAdded(CompMediaType, "itemEncoding"))
	assert.False(t, BreakingModified(CompMediaType, "itemEncoding"))
	assert.True(t, BreakingRemoved(CompMediaType, "itemEncoding"))

	// Test OAuth2MetadataUrl
	assert.False(t, BreakingAdded(CompSecurityScheme, PropOAuth2MetadataUrl))
	assert.False(t, BreakingModified(CompSecurityScheme, PropOAuth2MetadataUrl))
	assert.False(t, BreakingRemoved(CompSecurityScheme, PropOAuth2MetadataUrl))
}

// TestItemEncodingConfigurableBreakingRules verifies that CompareMediaTypes honours
// SetActiveBreakingRulesConfig for MediaType.ItemEncoding on removal-only diffs.
// The scenario is deliberately isolated to a removal so that flipping Removed to
// false changes the total breaking-change count, discriminating a correct
// BreakingRemoved(CompMediaType, PropItemEncoding) wiring from the old generic
// map-helper behavior which hard-codes removals as breaking regardless of the
// active config.
func TestItemEncodingConfigurableBreakingRules(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	// left has one itemEncoding entry, right has none — removal only.
	left := `schema:
  type: array
itemEncoding:
  removed-entry:
    contentType: application/json`

	right := `schema:
  type: array`

	var lNode, rNode yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(left), &lNode), "left YAML fixture must parse")
	require.NoError(t, yaml.Unmarshal([]byte(right), &rNode), "right YAML fixture must parse")

	lIdx := index.NewSpecIndexWithConfig(&lNode, index.CreateOpenAPIIndexConfig())
	rIdx := index.NewSpecIndexWithConfig(&rNode, index.CreateOpenAPIIndexConfig())
	ctx := context.Background()

	var lMT, rMT v3.MediaType
	require.NoError(t, low.BuildModel(&lNode, &lMT), "l model BuildModel must succeed")
	require.NoError(t, low.BuildModel(&rNode, &rMT), "r model BuildModel must succeed")

	require.NoError(t, lMT.Build(ctx, nil, lNode.Content[0], lIdx), "lMT Build must succeed")
	require.NoError(t, rMT.Build(ctx, nil, rNode.Content[0], rIdx), "rMT Build must succeed")

	// Default behavior: removing an itemEncoding entry should be breaking.
	changes := CompareMediaTypes(&lMT, &rMT)
	assert.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalBreakingChanges(),
		"default itemEncoding rules should classify the removal as breaking")

	// Override Removed=false: the removal should no longer be breaking.
	SetActiveBreakingRulesConfig(&BreakingRulesConfig{
		MediaType: &MediaTypeRules{
			ItemEncoding: &BreakingChangeRule{
				Added:    boolPtr(false),
				Modified: boolPtr(false),
				Removed:  boolPtr(false),
			},
		},
	})

	changes2 := CompareMediaTypes(&lMT, &rMT)
	assert.NotNil(t, changes2)
	assert.Equal(t, 0, changes2.TotalBreakingChanges(),
		"overridden itemEncoding rules should classify the removal as non-breaking")
}

// TestItemEncodingAddedConfigurable verifies that CompareMediaTypes honours
// SetActiveBreakingRulesConfig for MediaType.ItemEncoding on addition-only diffs.
// The scenario is deliberately isolated to an addition so that flipping Added to
// true changes the total breaking-change count, discriminating a correct
// BreakingAdded(CompMediaType, PropItemEncoding) wiring from the old generic
// map-helper behavior which hard-codes additions as non-breaking regardless of
// the active config.
func TestItemEncodingAddedConfigurable(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	// left has no itemEncoding, right has one entry — addition only.
	left := `schema:
  type: array`

	right := `schema:
  type: array
itemEncoding:
  added-entry:
    contentType: application/xml`

	var lNode, rNode yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(left), &lNode), "left YAML fixture must parse")
	require.NoError(t, yaml.Unmarshal([]byte(right), &rNode), "right YAML fixture must parse")

	lIdx := index.NewSpecIndexWithConfig(&lNode, index.CreateOpenAPIIndexConfig())
	rIdx := index.NewSpecIndexWithConfig(&rNode, index.CreateOpenAPIIndexConfig())
	ctx := context.Background()

	var lMT, rMT v3.MediaType
	require.NoError(t, low.BuildModel(&lNode, &lMT), "l model BuildModel must succeed")
	require.NoError(t, low.BuildModel(&rNode, &rMT), "r model BuildModel must succeed")

	require.NoError(t, lMT.Build(ctx, nil, lNode.Content[0], lIdx), "lMT Build must succeed")
	require.NoError(t, rMT.Build(ctx, nil, rNode.Content[0], rIdx), "rMT Build must succeed")

	// Default behavior: adding an itemEncoding entry should not be breaking.
	changes := CompareMediaTypes(&lMT, &rMT)
	assert.NotNil(t, changes)
	assert.Equal(t, 0, changes.TotalBreakingChanges(),
		"default itemEncoding rules should classify the addition as non-breaking")

	// Override Added=true: the addition should now be breaking.
	SetActiveBreakingRulesConfig(&BreakingRulesConfig{
		MediaType: &MediaTypeRules{
			ItemEncoding: &BreakingChangeRule{
				Added:    boolPtr(true),
				Modified: boolPtr(false),
				Removed:  boolPtr(true),
			},
		},
	})

	changes2 := CompareMediaTypes(&lMT, &rMT)
	assert.NotNil(t, changes2)
	assert.Equal(t, 1, changes2.TotalBreakingChanges(),
		"overridden itemEncoding rules should classify the addition as breaking")
}

// TestQueryAddedChangePayload verifies that when a Query operation is added
// (present on the right, absent on the left), the resulting Change record
// carries the right-hand Query operation as its NewObject rather than the
// left-hand nil value. This guards against a common wiring bug where
// ComparePathItems passes `lPath.Query.Value` (nil) into CreateChange instead
// of `rPath.Query.Value`, which produces a change record with a nil payload
// that downstream consumers cannot introspect.
func TestQueryAddedChangePayload(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `get:
  summary: Get resources`

	right := `get:
  summary: Get resources
query:
  summary: Query resources
  operationId: queryResources`

	var lNode, rNode yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(left), &lNode), "left YAML fixture must parse")
	require.NoError(t, yaml.Unmarshal([]byte(right), &rNode), "right YAML fixture must parse")

	lIdx := index.NewSpecIndexWithConfig(&lNode, index.CreateOpenAPIIndexConfig())
	rIdx := index.NewSpecIndexWithConfig(&rNode, index.CreateOpenAPIIndexConfig())
	ctx := context.Background()

	var lPath, rPath v3.PathItem
	require.NoError(t, low.BuildModel(&lNode, &lPath), "l model BuildModel must succeed")
	require.NoError(t, low.BuildModel(&rNode, &rPath), "r model BuildModel must succeed")

	require.NoError(t, lPath.Build(ctx, nil, lNode.Content[0], lIdx), "lPath Build must succeed")
	require.NoError(t, rPath.Build(ctx, nil, rNode.Content[0], rIdx), "rPath Build must succeed")

	changes := ComparePathItems(&lPath, &rPath)
	assert.NotNil(t, changes)

	// Locate the Query PropertyAdded change and assert it carries a
	// non-nil right-hand payload as NewObject.
	var queryAdded *Change
	for _, ch := range changes.GetAllChanges() {
		if ch.Property == v3.QueryLabel && ch.ChangeType == PropertyAdded {
			queryAdded = ch
			break
		}
	}
	assert.NotNil(t, queryAdded, "expected a PropertyAdded change for v3.QueryLabel")
	if queryAdded != nil {
		assert.NotNil(t, queryAdded.NewObject,
			"the added-Query change must carry the right-hand Query operation as NewObject, not nil")
		assert.Nil(t, queryAdded.OriginalObject,
			"the added-Query change should have no OriginalObject (left side had no query)")
		assert.False(t, queryAdded.Breaking,
			"adding query is non-breaking by default (Query.Added is false in the default polarity table)")
	}
}

// TestAdditionalOperationsReorderedMatching guards against an index bug in
// ComparePathItems where the additionalOperations matching loop looked up the
// left operation with `lKeys[j]` (the right-side index) instead of `lKeys[i]`
// (the left-side index). When left and right contain the same operation name
// in different positions, the buggy path either panics (right slice longer
// than left, `lKeys[j]` out of range) or compares the wrong left operation
// against the matched right operation, reporting spurious modifications.
func TestAdditionalOperationsReorderedMatching(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	// Two additionalOperations on each side. The name SHARED appears in both,
	// but at different positions: left has (SHARED, ONLY_LEFT), right has
	// (ONLY_RIGHT, SHARED). SHARED carries an identical operation body on
	// both sides so a correct implementation reports zero changes for it and
	// only surfaces the removal of ONLY_LEFT and the addition of ONLY_RIGHT.
	left := `additionalOperations:
  SHARED:
    summary: Shared operation
    operationId: sharedOp
  ONLY_LEFT:
    summary: Only-left operation
    operationId: onlyLeftOp`

	right := `additionalOperations:
  ONLY_RIGHT:
    summary: Only-right operation
    operationId: onlyRightOp
  SHARED:
    summary: Shared operation
    operationId: sharedOp`

	var lNode, rNode yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(left), &lNode), "left YAML fixture must parse")
	require.NoError(t, yaml.Unmarshal([]byte(right), &rNode), "right YAML fixture must parse")

	lIdx := index.NewSpecIndexWithConfig(&lNode, index.CreateOpenAPIIndexConfig())
	rIdx := index.NewSpecIndexWithConfig(&rNode, index.CreateOpenAPIIndexConfig())
	ctx := context.Background()

	var lPath, rPath v3.PathItem
	require.NoError(t, low.BuildModel(&lNode, &lPath), "l model BuildModel must succeed")
	require.NoError(t, low.BuildModel(&rNode, &rPath), "r model BuildModel must succeed")

	require.NoError(t, lPath.Build(ctx, nil, lNode.Content[0], lIdx), "lPath Build must succeed")
	require.NoError(t, rPath.Build(ctx, nil, rNode.Content[0], rIdx), "rPath Build must succeed")

	var changes *PathItemChanges
	assert.NotPanics(t, func() {
		changes = ComparePathItems(&lPath, &rPath)
	}, "ComparePathItems must not panic when additionalOperations keys appear in different positions on each side")

	assert.NotNil(t, changes)
	// Expected: ONLY_LEFT removed + ONLY_RIGHT added. SHARED is unchanged.
	// The buggy implementation additionally reports a spurious SHARED
	// modification (or panics) because it compares the wrong left operation
	// against the matched right one, inflating the change count above two.
	assert.Equal(t, 2, changes.TotalChanges(),
		"expected exactly two changes (ONLY_LEFT removed, ONLY_RIGHT added); a higher count indicates the reordering-matching bug is still present")
}

// TestEvalonGolden_QueryAddedExactPayload strengthens the right-only Query
// case by asserting the emitted PropertyAdded record carries the exact
// right-hand operation as NewObject (pointer identity with rPath.Query.Value)
// and no OriginalObject. This rejects implementations that satisfy the
// weaker "NewObject != nil" contract with an arbitrary (e.g. left-hand)
// value.
func TestEvalonGolden_QueryAddedExactPayload(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `get:
  summary: Get resources`

	right := `get:
  summary: Get resources
query:
  summary: Query resources
  operationId: queryResources`

	var lNode, rNode yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(left), &lNode), "left fixture must parse")
	require.NoError(t, yaml.Unmarshal([]byte(right), &rNode), "right fixture must parse")

	lIdx := index.NewSpecIndexWithConfig(&lNode, index.CreateOpenAPIIndexConfig())
	rIdx := index.NewSpecIndexWithConfig(&rNode, index.CreateOpenAPIIndexConfig())
	ctx := context.Background()

	var lPath, rPath v3.PathItem
	require.NoError(t, low.BuildModel(&lNode, &lPath), "left PathItem BuildModel must succeed")
	require.NoError(t, low.BuildModel(&rNode, &rPath), "right PathItem BuildModel must succeed")

	require.NoError(t, lPath.Build(ctx, nil, lNode.Content[0], lIdx), "left PathItem Build must succeed")
	require.NoError(t, rPath.Build(ctx, nil, rNode.Content[0], rIdx), "right PathItem Build must succeed")

	changes := ComparePathItems(&lPath, &rPath)
	assert.NotNil(t, changes)

	var queryAdded *Change
	for _, ch := range changes.GetAllChanges() {
		if ch.Property == v3.QueryLabel && ch.ChangeType == PropertyAdded {
			queryAdded = ch
			break
		}
	}
	assert.NotNil(t, queryAdded, "expected a PropertyAdded change for v3.QueryLabel")
	if queryAdded == nil {
		return
	}
	assert.Same(t, rPath.Query.Value, queryAdded.NewObject,
		"NewObject must be the right-hand Query operation exactly, not an arbitrary non-nil value")
	assert.Nil(t, queryAdded.OriginalObject,
		"OriginalObject must be nil when query appears only on the right side")
	assert.False(t, queryAdded.Breaking,
		"adding query is non-breaking by default (Query.Added is false in the default polarity table)")
}

// TestEvalonGolden_AdditionalOperationsReorderedExactRecords strengthens the
// reordered-keys scenario by asserting the two emitted records are exactly
// ONE PropertyRemoved for the ONLY_LEFT key and ONE PropertyAdded for the
// ONLY_RIGHT key, and that NO Change record is emitted for the SHARED key.
// This rejects implementations that produce a compensating pair of Changes
// on the wrong keys but happen to sum to the same total count.
func TestEvalonGolden_AdditionalOperationsReorderedExactRecords(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `additionalOperations:
  SHARED:
    summary: Shared operation
    operationId: sharedOp
  ONLY_LEFT:
    summary: Only-left operation
    operationId: onlyLeftOp`

	right := `additionalOperations:
  ONLY_RIGHT:
    summary: Only-right operation
    operationId: onlyRightOp
  SHARED:
    summary: Shared operation
    operationId: sharedOp`

	var lNode, rNode yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(left), &lNode), "left fixture must parse")
	require.NoError(t, yaml.Unmarshal([]byte(right), &rNode), "right fixture must parse")

	lIdx := index.NewSpecIndexWithConfig(&lNode, index.CreateOpenAPIIndexConfig())
	rIdx := index.NewSpecIndexWithConfig(&rNode, index.CreateOpenAPIIndexConfig())
	ctx := context.Background()

	var lPath, rPath v3.PathItem
	require.NoError(t, low.BuildModel(&lNode, &lPath), "left PathItem BuildModel must succeed")
	require.NoError(t, low.BuildModel(&rNode, &rPath), "right PathItem BuildModel must succeed")

	require.NoError(t, lPath.Build(ctx, nil, lNode.Content[0], lIdx), "left PathItem Build must succeed")
	require.NoError(t, rPath.Build(ctx, nil, rNode.Content[0], rIdx), "right PathItem Build must succeed")

	changes := ComparePathItems(&lPath, &rPath)
	assert.NotNil(t, changes)

	additionalOpsRecords := make([]*Change, 0)
	for _, ch := range changes.GetAllChanges() {
		if ch.Property == v3.AdditionalOperationsLabel {
			additionalOpsRecords = append(additionalOpsRecords, ch)
		}
	}
	assert.Len(t, additionalOpsRecords, 2,
		"exactly two AdditionalOperations records expected: ONLY_LEFT removed and ONLY_RIGHT added")

	sawOnlyLeftRemoved := false
	sawOnlyRightAdded := false
	sharedSpurious := false
	for _, ch := range additionalOpsRecords {
		if ch.ChangeType == PropertyRemoved {
			if op, ok := ch.OriginalObject.(*v3.Operation); ok && op != nil && op.Summary.Value == "Only-left operation" {
				sawOnlyLeftRemoved = true
			}
		}
		if ch.ChangeType == PropertyAdded {
			if op, ok := ch.NewObject.(*v3.Operation); ok && op != nil && op.Summary.Value == "Only-right operation" {
				sawOnlyRightAdded = true
			}
		}
		if op, ok := ch.NewObject.(*v3.Operation); ok && op != nil && op.Summary.Value == "Shared operation" {
			sharedSpurious = true
		}
		if op, ok := ch.OriginalObject.(*v3.Operation); ok && op != nil && op.Summary.Value == "Shared operation" {
			sharedSpurious = true
		}
	}
	assert.True(t, sawOnlyLeftRemoved,
		"expected exactly one PropertyRemoved carrying the ONLY_LEFT operation as OriginalObject")
	assert.True(t, sawOnlyRightAdded,
		"expected exactly one PropertyAdded carrying the ONLY_RIGHT operation as NewObject")
	assert.False(t, sharedSpurious,
		"no AdditionalOperations record should reference the SHARED operation")
}

// TestEvalonGolden_AdditionalOperationsPerKeyOverridesHonored covers the case
// where BOTH sides have additionalOperations with one key added and one key
// removed. Flipping the active AdditionalOperations Added and Removed
// polarities via SetActiveBreakingRulesConfig must observably change the
// Breaking flag on each emitted per-key Change record. Rejects implementations
// where per-key add/remove records ignore the active configuration.
func TestEvalonGolden_AdditionalOperationsPerKeyOverridesHonored(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `additionalOperations:
  KEEP:
    summary: Kept operation
    operationId: keepOp
  DROP:
    summary: Dropped operation
    operationId: dropOp`

	right := `additionalOperations:
  KEEP:
    summary: Kept operation
    operationId: keepOp
  ADD:
    summary: Added operation
    operationId: addOp`

	var lNode, rNode yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(left), &lNode), "left fixture must parse")
	require.NoError(t, yaml.Unmarshal([]byte(right), &rNode), "right fixture must parse")

	lIdx := index.NewSpecIndexWithConfig(&lNode, index.CreateOpenAPIIndexConfig())
	rIdx := index.NewSpecIndexWithConfig(&rNode, index.CreateOpenAPIIndexConfig())
	ctx := context.Background()

	var lPath, rPath v3.PathItem
	require.NoError(t, low.BuildModel(&lNode, &lPath), "left PathItem BuildModel must succeed")
	require.NoError(t, low.BuildModel(&rNode, &rPath), "right PathItem BuildModel must succeed")

	require.NoError(t, lPath.Build(ctx, nil, lNode.Content[0], lIdx), "left PathItem Build must succeed")
	require.NoError(t, rPath.Build(ctx, nil, rNode.Content[0], rIdx), "right PathItem Build must succeed")

	// Default: Added=false Removed=true. Every per-key remove must be breaking;
	// every per-key add must be non-breaking.
	defaultChanges := ComparePathItems(&lPath, &rPath)
	assert.NotNil(t, defaultChanges)
	sawDefaultRemoveBreaking := false
	sawDefaultAddNonBreaking := false
	for _, ch := range defaultChanges.GetAllChanges() {
		if ch.Property != v3.AdditionalOperationsLabel {
			continue
		}
		if ch.ChangeType == PropertyRemoved {
			assert.True(t, ch.Breaking,
				"default AdditionalOperations remove polarity must classify per-key removal as breaking")
			sawDefaultRemoveBreaking = true
		}
		if ch.ChangeType == PropertyAdded {
			assert.False(t, ch.Breaking,
				"default AdditionalOperations add polarity must classify per-key addition as non-breaking")
			sawDefaultAddNonBreaking = true
		}
	}
	assert.True(t, sawDefaultRemoveBreaking, "default run must emit at least one per-key remove")
	assert.True(t, sawDefaultAddNonBreaking, "default run must emit at least one per-key add")

	// Flip polarities: Added=true, Removed=false. Now every per-key add must
	// be breaking, every per-key remove must be non-breaking.
	SetActiveBreakingRulesConfig(&BreakingRulesConfig{
		PathItem: &PathItemRules{
			AdditionalOperations: &BreakingChangeRule{
				Added:    boolPtr(true),
				Modified: boolPtr(false),
				Removed:  boolPtr(false),
			},
		},
	})
	low.ClearHashCache()

	flippedChanges := ComparePathItems(&lPath, &rPath)
	assert.NotNil(t, flippedChanges)
	sawFlippedRemoveNonBreaking := false
	sawFlippedAddBreaking := false
	for _, ch := range flippedChanges.GetAllChanges() {
		if ch.Property != v3.AdditionalOperationsLabel {
			continue
		}
		if ch.ChangeType == PropertyRemoved {
			assert.False(t, ch.Breaking,
				"flipped AdditionalOperations remove polarity must classify per-key removal as non-breaking")
			sawFlippedRemoveNonBreaking = true
		}
		if ch.ChangeType == PropertyAdded {
			assert.True(t, ch.Breaking,
				"flipped AdditionalOperations add polarity must classify per-key addition as breaking")
			sawFlippedAddBreaking = true
		}
	}
	assert.True(t, sawFlippedRemoveNonBreaking, "flipped run must emit at least one per-key remove")
	assert.True(t, sawFlippedAddBreaking, "flipped run must emit at least one per-key add")
}

// TestEvalonGolden_ItemEncodingNilValueNodeSafety programmatically constructs
// two MediaTypes whose ItemEncoding entries have nil ValueNode fields, so the
// GetValueNode() lookup in CompareMediaTypes returns nil. Implementations
// that dereference GetValueNode().Value without a nil guard panic here. The
// reference guards `if node != nil` before reading `node.Value`, so the
// comparison succeeds and emits the correct addition and removal records.
func TestEvalonGolden_ItemEncodingNilValueNodeSafety(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	buildMT := func(entries [][2]string) *v3.MediaType {
		mt := &v3.MediaType{}
		if len(entries) == 0 {
			return mt
		}
		m := orderedmap.New[low.KeyReference[string], low.ValueReference[*v3.Encoding]]()
		for _, entry := range entries {
			m.Set(
				low.KeyReference[string]{Value: entry[0]},
				low.ValueReference[*v3.Encoding]{
					Value:     &v3.Encoding{ContentType: low.NodeReference[string]{Value: entry[1]}},
					ValueNode: nil,
				},
			)
		}
		mt.ItemEncoding = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*v3.Encoding]]]{
			Value: m,
		}
		return mt
	}

	lMT := buildMT([][2]string{
		{"stays", "application/json"},
		{"removed-entry", "application/xml"},
	})
	rMT := buildMT([][2]string{
		{"stays", "application/json"},
		{"added-entry", "text/plain"},
	})

	var changes *MediaTypeChanges
	assert.NotPanics(t, func() {
		changes = CompareMediaTypes(lMT, rMT)
	}, "CompareMediaTypes must not panic on ItemEncoding entries with nil ValueNode")
	assert.NotNil(t, changes)

	sawAdd := false
	sawRemove := false
	for _, ch := range changes.GetAllChanges() {
		if ch.Property != v3.ItemEncodingLabel {
			continue
		}
		switch ch.ChangeType {
		case ObjectAdded:
			assert.False(t, ch.Breaking,
				"default ItemEncoding add polarity must classify addition as non-breaking")
			sawAdd = true
		case ObjectRemoved:
			assert.True(t, ch.Breaking,
				"default ItemEncoding remove polarity must classify removal as breaking")
			sawRemove = true
		}
	}
	assert.True(t, sawAdd, "expected one ItemEncoding ObjectAdded record")
	assert.True(t, sawRemove, "expected one ItemEncoding ObjectRemoved record")
}

// TestEvalonGolden_OAuth2MetadataUrlPerTransitionOverrides asserts that
// CompareSecuritySchemes honors DISTINCT active-config polarities for each of
// the three OAuth2MetadataUrl transitions. A previous implementation reused a
// single boolean for all three transitions, so overriding Added or Removed
// without touching Modified silently had no effect. This test flips each of
// the three transitions to a distinct non-default polarity and verifies each
// emitted Change record reflects the correct per-transition override.
func TestEvalonGolden_OAuth2MetadataUrlPerTransitionOverrides(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	// Set distinct polarities per transition: Added=true, Modified=false,
	// Removed=true. Defaults are all false, so any test that only surfaces
	// the Modified transition (the reference's prior bug) cannot distinguish
	// these overrides.
	SetActiveBreakingRulesConfig(&BreakingRulesConfig{
		SecurityScheme: &SecuritySchemeRules{
			OAuth2MetadataUrl: &BreakingChangeRule{
				Added:    boolPtr(true),
				Modified: boolPtr(false),
				Removed:  boolPtr(true),
			},
		},
	})

	compare := func(t *testing.T, left, right string) *SecuritySchemeChanges {
		var lNode, rNode yaml.Node
		require.NoError(t, yaml.Unmarshal([]byte(left), &lNode), "left fixture must parse")
		require.NoError(t, yaml.Unmarshal([]byte(right), &rNode), "right fixture must parse")

		ctx := context.Background()

		var lSS, rSS v3.SecurityScheme
		require.NoError(t, low.BuildModel(lNode.Content[0], &lSS), "left SecurityScheme BuildModel must succeed")
		require.NoError(t, low.BuildModel(rNode.Content[0], &rSS), "right SecurityScheme BuildModel must succeed")

		require.NoError(t, lSS.Build(ctx, nil, lNode.Content[0], nil), "left SecurityScheme Build must succeed")
		require.NoError(t, rSS.Build(ctx, nil, rNode.Content[0], nil), "right SecurityScheme Build must succeed")

		return CompareSecuritySchemes(&lSS, &rSS)
	}

	// Addition: left has no oauth2MetadataUrl, right introduces one.
	// Under overridden Added=true, the emitted addition must be breaking.
	addedChanges := compare(t,
		`type: oauth2
description: OAuth2 auth`,
		`type: oauth2
description: OAuth2 auth
oauth2MetadataUrl: https://example.com/.well-known/oauth-authorization-server`)
	assert.NotNil(t, addedChanges)
	sawAdded := false
	for _, ch := range addedChanges.GetAllChanges() {
		if ch.Property != v3.OAuth2MetadataUrlLabel {
			continue
		}
		if ch.ChangeType == PropertyAdded {
			assert.True(t, ch.Breaking,
				"overridden OAuth2MetadataUrl.Added=true must classify the added record as breaking")
			sawAdded = true
		}
	}
	assert.True(t, sawAdded, "expected one OAuth2MetadataUrl PropertyAdded record")

	// Removal: left has oauth2MetadataUrl, right removes it.
	// Under overridden Removed=true, the emitted removal must be breaking.
	removedChanges := compare(t,
		`type: oauth2
description: OAuth2 auth
oauth2MetadataUrl: https://example.com/.well-known/oauth-authorization-server`,
		`type: oauth2
description: OAuth2 auth`)
	assert.NotNil(t, removedChanges)
	sawRemoved := false
	for _, ch := range removedChanges.GetAllChanges() {
		if ch.Property != v3.OAuth2MetadataUrlLabel {
			continue
		}
		if ch.ChangeType == PropertyRemoved {
			assert.True(t, ch.Breaking,
				"overridden OAuth2MetadataUrl.Removed=true must classify the removed record as breaking")
			sawRemoved = true
		}
	}
	assert.True(t, sawRemoved, "expected one OAuth2MetadataUrl PropertyRemoved record")

	// Modification: both sides have oauth2MetadataUrl with different URLs.
	// Under overridden Modified=false, the emitted modification must be non-breaking.
	modifiedChanges := compare(t,
		`type: oauth2
description: OAuth2 auth
oauth2MetadataUrl: https://example.com/.well-known/oauth-authorization-server`,
		`type: oauth2
description: OAuth2 auth
oauth2MetadataUrl: https://example.com/v2/.well-known/oauth-authorization-server`)
	assert.NotNil(t, modifiedChanges)
	sawModified := false
	for _, ch := range modifiedChanges.GetAllChanges() {
		if ch.Property != v3.OAuth2MetadataUrlLabel {
			continue
		}
		if ch.ChangeType == Modified {
			assert.False(t, ch.Breaking,
				"overridden OAuth2MetadataUrl.Modified=false must classify the modified record as non-breaking")
			sawModified = true
		}
	}
	assert.True(t, sawModified, "expected one OAuth2MetadataUrl Modified record")
}

// TestEvalonGolden_ItemEncodingChangeKeyPreservation asserts that when
// itemEncoding entries added or removed via CompareMediaTypes originate from
// YAML text, the emitted Change records preserve the affected map key in
// Change.New / Change.Original for downstream diagnostics. A replacement for
// the generic map-comparison path that passes an empty YAML mapping node's
// Value directly into CreateChange without populating it from the map key can
// lose the affected key on ordinary YAML-backed fixtures while still passing
// polarity, override, and nil-node safety checks.
func TestEvalonGolden_ItemEncodingChangeKeyPreservation(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `schema:
  type: array
itemEncoding:
  kept-entry:
    contentType: application/json
  removed-entry:
    contentType: application/xml`

	right := `schema:
  type: array
itemEncoding:
  kept-entry:
    contentType: application/json
  added-entry:
    contentType: text/plain`

	var lNode, rNode yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(left), &lNode), "left YAML fixture must parse")
	require.NoError(t, yaml.Unmarshal([]byte(right), &rNode), "right YAML fixture must parse")

	lIdx := index.NewSpecIndexWithConfig(&lNode, index.CreateOpenAPIIndexConfig())
	rIdx := index.NewSpecIndexWithConfig(&rNode, index.CreateOpenAPIIndexConfig())
	ctx := context.Background()

	var lMT, rMT v3.MediaType
	require.NoError(t, low.BuildModel(&lNode, &lMT), "left MediaType BuildModel must succeed")
	require.NoError(t, low.BuildModel(&rNode, &rMT), "right MediaType BuildModel must succeed")

	require.NoError(t, lMT.Build(ctx, nil, lNode.Content[0], lIdx), "left MediaType Build must succeed")
	require.NoError(t, rMT.Build(ctx, nil, rNode.Content[0], rIdx), "right MediaType Build must succeed")

	changes := CompareMediaTypes(&lMT, &rMT)
	assert.NotNil(t, changes)

	sawRemovedKey := false
	sawAddedKey := false
	for _, ch := range changes.GetAllChanges() {
		if ch.Property != v3.ItemEncodingLabel {
			continue
		}
		switch ch.ChangeType {
		case ObjectRemoved:
			assert.Equal(t, "removed-entry", ch.Original,
				"ObjectRemoved record for itemEncoding must preserve the removed map key in Change.Original")
			assert.Empty(t, ch.New,
				"ObjectRemoved record must not carry a right-side key in Change.New")
			sawRemovedKey = true
		case ObjectAdded:
			assert.Equal(t, "added-entry", ch.New,
				"ObjectAdded record for itemEncoding must preserve the added map key in Change.New")
			assert.Empty(t, ch.Original,
				"ObjectAdded record must not carry a left-side key in Change.Original")
			sawAddedKey = true
		}
	}
	assert.True(t, sawRemovedKey, "expected one ObjectRemoved record for itemEncoding key 'removed-entry'")
	assert.True(t, sawAddedKey, "expected one ObjectAdded record for itemEncoding key 'added-entry'")
}

// TestEvalonGolden_ItemEncodingSameKeyDifferentValueDetected exercises the
// manual ItemEncoding map-diff loop's MODIFICATION branch: two MediaTypes
// with the SAME encoding key but DIFFERENT encoding values must produce a
// nested EncodingChanges record under ItemEncodingChanges for that key. This
// closes coverage for the low.GenerateHashString-based equality check and the
// mc.ItemEncodingChanges[k] = CompareEncoding(...) assignment in media_type.go
// that add/remove tests alone cannot exercise.
func TestEvalonGolden_ItemEncodingSameKeyDifferentValueDetected(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `schema:
  type: array
itemEncoding:
  shared-key:
    contentType: application/json`

	right := `schema:
  type: array
itemEncoding:
  shared-key:
    contentType: application/xml`

	var lNode, rNode yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(left), &lNode))
	require.NoError(t, yaml.Unmarshal([]byte(right), &rNode))

	lIdx := index.NewSpecIndexWithConfig(&lNode, index.CreateOpenAPIIndexConfig())
	rIdx := index.NewSpecIndexWithConfig(&rNode, index.CreateOpenAPIIndexConfig())
	ctx := context.Background()

	var lMT, rMT v3.MediaType
	require.NoError(t, low.BuildModel(&lNode, &lMT))
	require.NoError(t, low.BuildModel(&rNode, &rMT))
	require.NoError(t, lMT.Build(ctx, nil, lNode.Content[0], lIdx))
	require.NoError(t, rMT.Build(ctx, nil, rNode.Content[0], rIdx))

	changes := CompareMediaTypes(&lMT, &rMT)
	require.NotNil(t, changes, "same-key differing-value ItemEncoding must produce a non-nil MediaTypeChanges")
	require.NotNil(t, changes.ItemEncodingChanges, "ItemEncodingChanges must be populated when the shared key's encoding value differs")

	nested, ok := changes.ItemEncodingChanges["shared-key"]
	require.True(t, ok, "expected a nested EncodingChanges record under the shared key")
	require.NotNil(t, nested)
	assert.Greater(t, nested.TotalChanges(), 0,
		"nested EncodingChanges must contain at least one change reflecting the ContentType difference")
}
