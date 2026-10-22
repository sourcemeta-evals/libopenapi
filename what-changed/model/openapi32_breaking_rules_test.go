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
	require.NotNil(t, config.PathItem.Query)
	assert.False(t, *config.PathItem.Query.Added, "Adding query should not be breaking by default")
	assert.False(t, *config.PathItem.Query.Modified, "Modifying query should not be breaking by default")
	assert.True(t, *config.PathItem.Query.Removed, "Removing query should be breaking by default")
}

// TestDefaultBreakingRules_AdditionalOperations verifies the default breaking rules for AdditionalOperations
func TestDefaultBreakingRules_AdditionalOperations(t *testing.T) {
	ResetDefaultBreakingRules()
	config := GenerateDefaultBreakingRules()

	require.NotNil(t, config.PathItem.AdditionalOperations)
	assert.False(t, *config.PathItem.AdditionalOperations.Added, "Adding additionalOperations should not be breaking by default")
	assert.False(t, *config.PathItem.AdditionalOperations.Modified, "Modifying additionalOperations should not be breaking by default")
	assert.True(t, *config.PathItem.AdditionalOperations.Removed, "Removing additionalOperations should be breaking by default")
}

// TestDefaultBreakingRules_ItemSchema verifies the default breaking rules for ItemSchema
func TestDefaultBreakingRules_ItemSchema(t *testing.T) {
	ResetDefaultBreakingRules()
	config := GenerateDefaultBreakingRules()

	require.NotNil(t, config.MediaType.ItemSchema)
	assert.True(t, *config.MediaType.ItemSchema.Added, "Adding itemSchema should be breaking by default")
	assert.False(t, *config.MediaType.ItemSchema.Modified, "Modifying itemSchema should not be breaking by default")
	assert.True(t, *config.MediaType.ItemSchema.Removed, "Removing itemSchema should be breaking by default")
}

// TestDefaultBreakingRules_ItemEncoding verifies the default breaking rules for ItemEncoding
func TestDefaultBreakingRules_ItemEncoding(t *testing.T) {
	ResetDefaultBreakingRules()
	config := GenerateDefaultBreakingRules()

	require.NotNil(t, config.MediaType.ItemEncoding)
	assert.False(t, *config.MediaType.ItemEncoding.Added, "Adding itemEncoding should not be breaking by default")
	assert.False(t, *config.MediaType.ItemEncoding.Modified, "Modifying itemEncoding should not be breaking by default")
	assert.True(t, *config.MediaType.ItemEncoding.Removed, "Removing itemEncoding should be breaking by default")
}

// TestDefaultBreakingRules_OAuth2MetadataUrl verifies the default breaking rules for OAuth2MetadataUrl
func TestDefaultBreakingRules_OAuth2MetadataUrl(t *testing.T) {
	ResetDefaultBreakingRules()
	config := GenerateDefaultBreakingRules()

	require.NotNil(t, config.SecurityScheme.OAuth2MetadataUrl)
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
// CompareSecuritySchemes honors the active-config polarity of each
// OAuth2MetadataUrl transition independently. A previous implementation reused
// a single boolean for all three transitions, so overriding Added or Removed
// without touching Modified silently had no effect.
//
// The table walks every combination of the three booleans rather than one
// hand-picked configuration. A boolean carries only two values, so no single
// configuration can tell all three lookups apart: one pair always shares a
// value, and a transition wired to the wrong accessor of that pair goes
// unnoticed. Across all eight combinations every pair of transitions differs
// in four of them, so any transition reading another transition's rule, or
// left on a hardcoded literal, fails here.
func TestEvalonGolden_OAuth2MetadataUrlPerTransitionOverrides(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	const (
		schemeWithoutURL = `type: oauth2
description: OAuth2 auth`
		metadataURL = "https://example.com/.well-known/oauth-authorization-server"
		changedURL  = "https://example.com/v2/.well-known/oauth-authorization-server"
	)

	schemeWithURL := func(url string) string {
		return schemeWithoutURL + "\noauth2MetadataUrl: " + url
	}

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

	findRecord := func(t *testing.T, changes *SecuritySchemeChanges, changeType int) *Change {
		require.NotNil(t, changes, "the comparison must report changes")
		var found *Change
		for _, ch := range changes.GetAllChanges() {
			if ch.Property != v3.OAuth2MetadataUrlLabel || ch.ChangeType != changeType {
				continue
			}
			require.Nil(t, found, "oauth2MetadataUrl must emit exactly one record per transition")
			found = ch
		}
		require.NotNil(t, found, "expected one oauth2MetadataUrl record for this transition")
		return found
	}

	transitions := []struct {
		name       string
		left       string
		right      string
		changeType int
		original   string
		updated    string
		polarity   func(added, modified, removed bool) bool
	}{
		{
			name:       "addition",
			left:       schemeWithoutURL,
			right:      schemeWithURL(metadataURL),
			changeType: PropertyAdded,
			original:   "",
			updated:    metadataURL,
			polarity:   func(added, _, _ bool) bool { return added },
		},
		{
			name:       "modification",
			left:       schemeWithURL(metadataURL),
			right:      schemeWithURL(changedURL),
			changeType: Modified,
			original:   metadataURL,
			updated:    changedURL,
			polarity:   func(_, modified, _ bool) bool { return modified },
		},
		{
			name:       "removal",
			left:       schemeWithURL(metadataURL),
			right:      schemeWithoutURL,
			changeType: PropertyRemoved,
			original:   metadataURL,
			updated:    "",
			polarity:   func(_, _, removed bool) bool { return removed },
		},
	}

	combinations := []struct{ added, modified, removed bool }{
		{false, false, false},
		{false, false, true},
		{false, true, false},
		{false, true, true},
		{true, false, false},
		{true, false, true},
		{true, true, false},
		{true, true, true},
	}

	for _, combination := range combinations {
		for _, transition := range transitions {
			SetActiveBreakingRulesConfig(&BreakingRulesConfig{
				SecurityScheme: &SecuritySchemeRules{
					OAuth2MetadataUrl: &BreakingChangeRule{
						Added:    boolPtr(combination.added),
						Modified: boolPtr(combination.modified),
						Removed:  boolPtr(combination.removed),
					},
				},
			})
			low.ClearHashCache()

			record := findRecord(t, compare(t, transition.left, transition.right), transition.changeType)

			assert.Equal(t, transition.polarity(combination.added, combination.modified, combination.removed),
				record.Breaking,
				"the oauth2MetadataUrl %s record must reflect its own configured polarity under Added=%v Modified=%v Removed=%v, not another transition's",
				transition.name, combination.added, combination.modified, combination.removed)
			assert.Equal(t, transition.original, record.Original,
				"the oauth2MetadataUrl %s record must carry the left-side URL in Change.Original", transition.name)
			assert.Equal(t, transition.updated, record.New,
				"the oauth2MetadataUrl %s record must carry the right-side URL in Change.New", transition.name)
		}
	}
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

// TestEvalonGolden_QueryRemovedExactPayload mirrors the added-Query golden
// test for the left-only direction: the PropertyRemoved record for query
// must carry the left-hand Query operation exactly as OriginalObject, not
// nil or an arbitrary substitute, with no NewObject and the default
// breaking polarity (Query.Removed is true).
func TestEvalonGolden_QueryRemovedExactPayload(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `get:
  summary: Get resources
query:
  summary: Query resources
  operationId: queryResources`

	right := `get:
  summary: Get resources`

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

	var queryRemoved *Change
	for _, ch := range changes.GetAllChanges() {
		if ch.Property == v3.QueryLabel && ch.ChangeType == PropertyRemoved {
			queryRemoved = ch
			break
		}
	}
	assert.NotNil(t, queryRemoved, "expected a PropertyRemoved change for v3.QueryLabel")
	if queryRemoved == nil {
		return
	}
	assert.Same(t, lPath.Query.Value, queryRemoved.OriginalObject,
		"OriginalObject must be the left-hand Query operation exactly, not an arbitrary non-nil value")
	assert.Nil(t, queryRemoved.NewObject,
		"NewObject must be nil when query appears only on the left side")
	assert.True(t, queryRemoved.Breaking,
		"removing query is breaking by default (Query.Removed is true in the default polarity table)")
}

// TestEvalonGolden_AdditionalOperationsSharedKeyModificationDetected hardens
// the shared-key matching path with keys whose operations actually differ.
// Left and right share two keys in swapped positions: UNCHANGED is identical
// on both sides and must produce no nested record, while CHANGED differs in
// summary and must produce a nested OperationChanges entry under its own key
// whose modification reflects the real left/right values. Rejects
// implementations that skip shared-key comparison entirely or match a shared
// key against the wrong counterpart from the other side.
func TestEvalonGolden_AdditionalOperationsSharedKeyModificationDetected(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	left := `additionalOperations:
  UNCHANGED:
    summary: Stable operation
    operationId: stableOp
  CHANGED:
    summary: Original summary
    operationId: changedOp`

	right := `additionalOperations:
  CHANGED:
    summary: Updated summary
    operationId: changedOp
  UNCHANGED:
    summary: Stable operation
    operationId: stableOp`

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

	for _, ch := range changes.GetAllChanges() {
		assert.NotEqual(t, v3.AdditionalOperationsLabel, ch.Property,
			"no add/remove record expected: both keys exist on both sides")
	}

	_, unchangedPresent := changes.AdditionalOperationChanges["UNCHANGED"]
	assert.False(t, unchangedPresent,
		"the identical UNCHANGED operation must not produce a nested OperationChanges entry")

	changedOperation, changedPresent := changes.AdditionalOperationChanges["CHANGED"]
	require.True(t, changedPresent,
		"the differing CHANGED operation must produce a nested OperationChanges entry under its own key")
	require.NotNil(t, changedOperation)

	sawSummaryModification := false
	for _, ch := range changedOperation.GetAllChanges() {
		if ch.ChangeType == Modified &&
			ch.Original == "Original summary" && ch.New == "Updated summary" {
			sawSummaryModification = true
		}
	}
	assert.True(t, sawSummaryModification,
		"the CHANGED entry must record the real summary modification (left value against right value), not a comparison against the wrong counterpart")
}

// TestEvalonGolden_ItemEncodingNestedChangeReferencePreserved asserts that a
// nested ItemEncoding modification keeps the reference metadata carried by the
// left entry. The generic map helper that the itemEncoding comparison replaced
// called SetReferenceIfExists after storing the nested changes, so a $ref-backed
// entry surfaced its reference through GetChangeReference(). A hand-written map
// loop that stores the CompareEncoding result without that call silently drops
// the reference, and consumers rendering a diff lose the pointer back to the
// component the entry came from.
func TestEvalonGolden_ItemEncodingNestedChangeReferencePreserved(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	const encodingReference = "#/components/encodings/FileEncoding"

	buildMT := func(key, contentType, reference string) *v3.MediaType {
		valueNode := &yaml.Node{Kind: yaml.MappingNode}
		contentTypeNode := &yaml.Node{Kind: yaml.ScalarNode, Value: contentType}
		entry := low.ValueReference[*v3.Encoding]{
			Value: &v3.Encoding{
				ContentType: low.NodeReference[string]{
					Value:     contentType,
					KeyNode:   contentTypeNode,
					ValueNode: contentTypeNode,
				},
				Reference: new(low.Reference),
			},
			ValueNode: valueNode,
		}
		if reference != "" {
			entry.SetReference(reference, valueNode)
		}
		entries := orderedmap.New[low.KeyReference[string], low.ValueReference[*v3.Encoding]]()
		entries.Set(low.KeyReference[string]{Value: key}, entry)
		return &v3.MediaType{
			ItemEncoding: low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*v3.Encoding]]]{
				Value: entries,
			},
		}
	}

	lMT := buildMT("file", "application/json", encodingReference)
	rMT := buildMT("file", "application/xml", "")

	changes := CompareMediaTypes(lMT, rMT)
	require.NotNil(t, changes,
		"a shared ItemEncoding key whose content types differ must report changes")

	nested := changes.ItemEncodingChanges["file"]
	require.NotNil(t, nested,
		"the shared ItemEncoding key must carry a nested EncodingChanges entry")
	assert.Equal(t, 1, nested.TotalChanges(),
		"the nested entry must record the content type modification")
	assert.Equal(t, encodingReference, nested.GetChangeReference(),
		"a nested ItemEncoding modification must preserve the reference carried by the left entry")
}

// TestEvalonGolden_OpenAPI32RuleFieldsDecodeFromYAML asserts that the five new
// OpenAPI 3.2 rule fields decode from a YAML configuration under their
// documented property names. The rule registry is keyed off each field's json
// tag, so a wrong json tag already surfaces as a failed lookup elsewhere, but
// the yaml tag has no such consumer and a typo in it silently drops the rule
// for anyone loading breaking-rule configuration from YAML.
func TestEvalonGolden_OpenAPI32RuleFieldsDecodeFromYAML(t *testing.T) {
	const config = `pathItem:
  query:
    removed: true
  additionalOperations:
    removed: true
mediaType:
  itemSchema:
    added: true
    removed: true
  itemEncoding:
    removed: true
securityScheme:
  oauth2MetadataUrl:
    added: true
    modified: true
    removed: true`

	var decoded BreakingRulesConfig
	require.NoError(t, yaml.Unmarshal([]byte(config), &decoded), "the rule configuration must decode")

	require.NotNil(t, decoded.PathItem, "pathItem rules must decode")
	require.NotNil(t, decoded.PathItem.Query,
		"the Query rule must decode from the `query` yaml key")
	require.NotNil(t, decoded.PathItem.AdditionalOperations,
		"the AdditionalOperations rule must decode from the `additionalOperations` yaml key")
	require.NotNil(t, decoded.MediaType, "mediaType rules must decode")
	require.NotNil(t, decoded.MediaType.ItemSchema,
		"the ItemSchema rule must decode from the `itemSchema` yaml key")
	require.NotNil(t, decoded.MediaType.ItemEncoding,
		"the ItemEncoding rule must decode from the `itemEncoding` yaml key")
	require.NotNil(t, decoded.SecurityScheme, "securityScheme rules must decode")
	require.NotNil(t, decoded.SecurityScheme.OAuth2MetadataUrl,
		"the OAuth2MetadataUrl rule must decode from the `oauth2MetadataUrl` yaml key")

	require.NotNil(t, decoded.PathItem.Query.Removed)
	assert.True(t, *decoded.PathItem.Query.Removed, "the decoded Query rule must carry removed=true")
	require.NotNil(t, decoded.MediaType.ItemSchema.Added)
	assert.True(t, *decoded.MediaType.ItemSchema.Added, "the decoded ItemSchema rule must carry added=true")
	require.NotNil(t, decoded.SecurityScheme.OAuth2MetadataUrl.Modified)
	assert.True(t, *decoded.SecurityScheme.OAuth2MetadataUrl.Modified,
		"the decoded OAuth2MetadataUrl rule must carry modified=true")
}

// TestEvalonGolden_AdditionalOperationsSharedKeyBeyondShorterMap exercises the
// shared-key matching loop when the two additionalOperations maps differ in
// size and the shared key sits past the shorter map's last index. The base
// implementation indexed the LEFT key slice with the RIGHT loop's index, so this
// input did not merely compare the wrong counterpart, it read out of range and
// panicked. Every other fixture here uses equal-sized maps, which cannot catch a
// reintroduced out-of-range access. Both directions are exercised because the
// matching loop walks the left map and the overrun depends on which side is
// longer.
func TestEvalonGolden_AdditionalOperationsSharedKeyBeyondShorterMap(t *testing.T) {
	ResetDefaultBreakingRules()
	ResetActiveBreakingRulesConfig()
	low.ClearHashCache()
	defer func() {
		ResetActiveBreakingRulesConfig()
		ResetDefaultBreakingRules()
	}()

	shorter := `additionalOperations:
  COPY:
    summary: Copy operation
    operationId: copyOp`

	longer := `additionalOperations:
  PURGE:
    summary: Purge operation
    operationId: purgeOp
  COPY:
    summary: Copy operation renamed
    operationId: copyOp`

	build := func(t *testing.T, src string) *v3.PathItem {
		var node yaml.Node
		require.NoError(t, yaml.Unmarshal([]byte(src), &node), "fixture must parse")
		idx := index.NewSpecIndexWithConfig(&node, index.CreateOpenAPIIndexConfig())
		var item v3.PathItem
		require.NoError(t, low.BuildModel(&node, &item), "PathItem BuildModel must succeed")
		require.NoError(t, item.Build(context.Background(), nil, node.Content[0], idx),
			"PathItem Build must succeed")
		return &item
	}

	// Right side longer: the shared COPY key is at right index 1, past the left
	// map's only index, which is where the out-of-range read occurred.
	var rightLongerChanges *PathItemChanges
	require.NotPanics(t, func() {
		low.ClearHashCache()
		rightLongerChanges = ComparePathItems(build(t, shorter), build(t, longer))
	}, "comparing a shared key past the shorter map's last index must not read out of range")
	require.NotNil(t, rightLongerChanges, "the differing shared operation must report changes")
	require.NotNil(t, rightLongerChanges.AdditionalOperationChanges["COPY"],
		"the shared COPY key must be compared against its real counterpart and report a modification")
	require.True(t, hasAdditionalOperationRecord(rightLongerChanges, PropertyAdded),
		"the right-only PURGE key must be reported as an addition")

	// Left side longer: same shape with the sides swapped.
	var leftLongerChanges *PathItemChanges
	require.NotPanics(t, func() {
		low.ClearHashCache()
		leftLongerChanges = ComparePathItems(build(t, longer), build(t, shorter))
	}, "comparing a shared key past the shorter map's last index must not read out of range")
	require.NotNil(t, leftLongerChanges, "the differing shared operation must report changes")
	require.NotNil(t, leftLongerChanges.AdditionalOperationChanges["COPY"],
		"the shared COPY key must be compared against its real counterpart and report a modification")
	require.True(t, hasAdditionalOperationRecord(leftLongerChanges, PropertyRemoved),
		"the left-only PURGE key must be reported as a removal")
}

// hasAdditionalOperationRecord reports whether the comparison emitted a
// top-level additionalOperations record of the given transition type.
func hasAdditionalOperationRecord(changes *PathItemChanges, changeType int) bool {
	for _, ch := range changes.GetAllChanges() {
		if ch.Property == v3.AdditionalOperationsLabel && ch.ChangeType == changeType {
			return true
		}
	}
	return false
}
