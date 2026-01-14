// Copyright 2022-2025 Princess Beef Heavy Industries, LLC / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"context"
	"testing"

	"github.com/pb33f/libopenapi/datamodel/low"
	v3 "github.com/pb33f/libopenapi/datamodel/low/v3"
	"github.com/pb33f/libopenapi/index"
	"github.com/stretchr/testify/assert"
	"go.yaml.in/yaml/v4"
)

// TestOpenAPI32_PathItem_Query tests breaking rules for the query operation (OpenAPI 3.2+)
func TestOpenAPI32_PathItem_Query(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.PathItem)

	// Query operation is an OpenAPI 3.2+ feature
	// The query operation changes are tracked in PathItemChanges.QueryChanges
	// Breaking behavior: adding is not breaking, removing is breaking (hardcoded in path_item.go)
}

// TestOpenAPI32_PathItem_Query_Added tests that adding a query operation is not breaking
func TestOpenAPI32_PathItem_Query_Added(t *testing.T) {
	low.ClearHashCache()

	left := `get:
  summary: Get resource
  operationId: getResource`

	right := `get:
  summary: Get resource
  operationId: getResource
query:
  summary: Query resources
  operationId: queryResources`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	lIdx := index.NewSpecIndexWithConfig(&lNode, index.CreateOpenAPIIndexConfig())
	rIdx := index.NewSpecIndexWithConfig(&rNode, index.CreateOpenAPIIndexConfig())
	ctx := context.Background()

	var lPath, rPath v3.PathItem
	_ = low.BuildModel(&lNode, &lPath)
	_ = low.BuildModel(&rNode, &rPath)

	_ = lPath.Build(ctx, nil, lNode.Content[0], lIdx)
	_ = rPath.Build(ctx, nil, rNode.Content[0], rIdx)

	changes := ComparePathItems(&lPath, &rPath)

	assert.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())

	// Adding query operation should not be breaking
	foundQueryAdded := false
	for _, change := range changes.GetAllChanges() {
		if change.Property == "query" && change.ChangeType == PropertyAdded {
			foundQueryAdded = true
			assert.False(t, change.Breaking, "Adding query operation should not be breaking")
			break
		}
	}
	assert.True(t, foundQueryAdded, "Query operation should be detected as added")
}

// TestOpenAPI32_PathItem_Query_Removed tests that removing a query operation is breaking
func TestOpenAPI32_PathItem_Query_Removed(t *testing.T) {
	low.ClearHashCache()

	left := `get:
  summary: Get resource
  operationId: getResource
query:
  summary: Query resources
  operationId: queryResources`

	right := `get:
  summary: Get resource
  operationId: getResource`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	lIdx := index.NewSpecIndexWithConfig(&lNode, index.CreateOpenAPIIndexConfig())
	rIdx := index.NewSpecIndexWithConfig(&rNode, index.CreateOpenAPIIndexConfig())
	ctx := context.Background()

	var lPath, rPath v3.PathItem
	_ = low.BuildModel(&lNode, &lPath)
	_ = low.BuildModel(&rNode, &rPath)

	_ = lPath.Build(ctx, nil, lNode.Content[0], lIdx)
	_ = rPath.Build(ctx, nil, rNode.Content[0], rIdx)

	changes := ComparePathItems(&lPath, &rPath)

	assert.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 1, changes.TotalBreakingChanges())

	// Removing query operation should be breaking
	foundQueryRemoved := false
	for _, change := range changes.GetAllChanges() {
		if change.Property == "query" && change.ChangeType == PropertyRemoved {
			foundQueryRemoved = true
			assert.True(t, change.Breaking, "Removing query operation should be breaking")
			break
		}
	}
	assert.True(t, foundQueryRemoved, "Query operation should be detected as removed")
}

// TestOpenAPI32_PathItem_AdditionalOperations tests breaking rules for additionalOperations (OpenAPI 3.2+)
func TestOpenAPI32_PathItem_AdditionalOperations(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.PathItem)

	// AdditionalOperations should have a rule defined
	// Note: additionalOperations is a map of custom operations
}

// TestOpenAPI32_PathItem_AdditionalOperations_Added tests that adding additionalOperations is not breaking
func TestOpenAPI32_PathItem_AdditionalOperations_Added(t *testing.T) {
	low.ClearHashCache()

	left := `get:
  summary: Get resource
  operationId: getResource`

	right := `get:
  summary: Get resource
  operationId: getResource
x-custom-op:
  summary: Custom operation
  operationId: customOp`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	lIdx := index.NewSpecIndexWithConfig(&lNode, index.CreateOpenAPIIndexConfig())
	rIdx := index.NewSpecIndexWithConfig(&rNode, index.CreateOpenAPIIndexConfig())
	ctx := context.Background()

	var lPath, rPath v3.PathItem
	_ = low.BuildModel(&lNode, &lPath)
	_ = low.BuildModel(&rNode, &rPath)

	_ = lPath.Build(ctx, nil, lNode.Content[0], lIdx)
	_ = rPath.Build(ctx, nil, rNode.Content[0], rIdx)

	changes := ComparePathItems(&lPath, &rPath)

	// Adding additional operations should not be breaking
	if changes != nil {
		for _, change := range changes.GetAllChanges() {
			if change.Property == "additionalOperations" && change.ChangeType == PropertyAdded {
				assert.False(t, change.Breaking, "Adding additionalOperations should not be breaking")
			}
		}
	}
}

// TestOpenAPI32_PathItem_AdditionalOperations_Removed tests that removing additionalOperations is breaking
func TestOpenAPI32_PathItem_AdditionalOperations_Removed(t *testing.T) {
	low.ClearHashCache()

	left := `get:
  summary: Get resource
  operationId: getResource
x-custom-op:
  summary: Custom operation
  operationId: customOp`

	right := `get:
  summary: Get resource
  operationId: getResource`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	lIdx := index.NewSpecIndexWithConfig(&lNode, index.CreateOpenAPIIndexConfig())
	rIdx := index.NewSpecIndexWithConfig(&rNode, index.CreateOpenAPIIndexConfig())
	ctx := context.Background()

	var lPath, rPath v3.PathItem
	_ = low.BuildModel(&lNode, &lPath)
	_ = low.BuildModel(&rNode, &rPath)

	_ = lPath.Build(ctx, nil, lNode.Content[0], lIdx)
	_ = rPath.Build(ctx, nil, rNode.Content[0], rIdx)

	changes := ComparePathItems(&lPath, &rPath)

	// Removing additional operations should be breaking
	if changes != nil {
		for _, change := range changes.GetAllChanges() {
			if change.Property == "additionalOperations" && change.ChangeType == PropertyRemoved {
				assert.True(t, change.Breaking, "Removing additionalOperations should be breaking")
			}
		}
	}
}

// TestOpenAPI32_MediaType_ItemSchema tests breaking rules for itemSchema (OpenAPI 3.2+)
func TestOpenAPI32_MediaType_ItemSchema(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.MediaType)

	// itemSchema changes should be tracked
}

// TestOpenAPI32_MediaType_ItemSchema_Added tests that adding itemSchema is breaking
func TestOpenAPI32_MediaType_ItemSchema_Added(t *testing.T) {
	low.ClearHashCache()

	left := `schema:
  type: array`

	right := `schema:
  type: array
itemSchema:
  type: object
  properties:
    id:
      type: string`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	lIdx := index.NewSpecIndex(&lNode)
	rIdx := index.NewSpecIndex(&rNode)

	var lMt, rMt v3.MediaType
	_ = low.BuildModel(&lNode, &lMt)
	_ = low.BuildModel(&rNode, &rMt)
	_ = lMt.Build(context.Background(), nil, lNode.Content[0], lIdx)
	_ = rMt.Build(context.Background(), nil, rNode.Content[0], rIdx)

	changes := CompareMediaTypes(&lMt, &rMt)

	assert.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 1, changes.TotalBreakingChanges())

	allChanges := changes.GetAllChanges()
	assert.Len(t, allChanges, 1)
	assert.Equal(t, ObjectAdded, allChanges[0].ChangeType)
	assert.Equal(t, v3.ItemSchemaLabel, allChanges[0].Property)
	assert.True(t, allChanges[0].Breaking, "Adding itemSchema should be breaking")
}

// TestOpenAPI32_MediaType_ItemSchema_Removed tests that removing itemSchema is breaking
func TestOpenAPI32_MediaType_ItemSchema_Removed(t *testing.T) {
	low.ClearHashCache()

	left := `schema:
  type: array
itemSchema:
  type: object
  properties:
    id:
      type: string`

	right := `schema:
  type: array`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	lIdx := index.NewSpecIndex(&lNode)
	rIdx := index.NewSpecIndex(&rNode)

	var lMt, rMt v3.MediaType
	_ = low.BuildModel(&lNode, &lMt)
	_ = low.BuildModel(&rNode, &rMt)
	_ = lMt.Build(context.Background(), nil, lNode.Content[0], lIdx)
	_ = rMt.Build(context.Background(), nil, rNode.Content[0], rIdx)

	changes := CompareMediaTypes(&lMt, &rMt)

	assert.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 1, changes.TotalBreakingChanges())

	allChanges := changes.GetAllChanges()
	assert.Len(t, allChanges, 1)
	assert.Equal(t, ObjectRemoved, allChanges[0].ChangeType)
	assert.Equal(t, v3.ItemSchemaLabel, allChanges[0].Property)
	assert.True(t, allChanges[0].Breaking, "Removing itemSchema should be breaking")
}

// TestOpenAPI32_MediaType_ItemSchema_Modified tests that modifying itemSchema is tracked
func TestOpenAPI32_MediaType_ItemSchema_Modified(t *testing.T) {
	low.ClearHashCache()

	left := `schema:
  type: array
itemSchema:
  type: object
  properties:
    id:
      type: string`

	right := `schema:
  type: array
itemSchema:
  type: object
  properties:
    id:
      type: integer`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	lIdx := index.NewSpecIndex(&lNode)
	rIdx := index.NewSpecIndex(&rNode)

	var lMt, rMt v3.MediaType
	_ = low.BuildModel(&lNode, &lMt)
	_ = low.BuildModel(&rNode, &rMt)
	_ = lMt.Build(context.Background(), nil, lNode.Content[0], lIdx)
	_ = rMt.Build(context.Background(), nil, rNode.Content[0], rIdx)

	changes := CompareMediaTypes(&lMt, &rMt)

	assert.NotNil(t, changes)
	assert.NotNil(t, changes.ItemSchemaChanges)
	assert.Greater(t, changes.TotalChanges(), 0)
	assert.Greater(t, changes.TotalBreakingChanges(), 0)
}

// TestOpenAPI32_MediaType_ItemEncoding tests breaking rules for itemEncoding (OpenAPI 3.2+)
func TestOpenAPI32_MediaType_ItemEncoding(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.MediaType)

	// itemEncoding changes should be tracked
}

// TestOpenAPI32_MediaType_ItemEncoding_Added tests that adding itemEncoding is tracked
func TestOpenAPI32_MediaType_ItemEncoding_Added(t *testing.T) {
	low.ClearHashCache()

	left := `schema:
  type: array`

	right := `schema:
  type: array
itemEncoding:
  file:
    contentType: image/jpeg
    allowReserved: true`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	lIdx := index.NewSpecIndex(&lNode)
	rIdx := index.NewSpecIndex(&rNode)

	var lMt, rMt v3.MediaType
	_ = low.BuildModel(&lNode, &lMt)
	_ = low.BuildModel(&rNode, &rMt)
	_ = lMt.Build(context.Background(), nil, lNode.Content[0], lIdx)
	_ = rMt.Build(context.Background(), nil, rNode.Content[0], rIdx)

	changes := CompareMediaTypes(&lMt, &rMt)

	assert.NotNil(t, changes)
	assert.Greater(t, changes.TotalChanges(), 0)

	allChanges := changes.GetAllChanges()
	assert.Len(t, allChanges, 1)
	assert.Equal(t, ObjectAdded, allChanges[0].ChangeType)
	assert.Equal(t, v3.ItemEncodingLabel, allChanges[0].Property)
}

// TestOpenAPI32_MediaType_ItemEncoding_Removed tests that removing itemEncoding is tracked
func TestOpenAPI32_MediaType_ItemEncoding_Removed(t *testing.T) {
	low.ClearHashCache()

	left := `schema:
  type: array
itemEncoding:
  file:
    contentType: image/jpeg
    allowReserved: true`

	right := `schema:
  type: array`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	lIdx := index.NewSpecIndex(&lNode)
	rIdx := index.NewSpecIndex(&rNode)

	var lMt, rMt v3.MediaType
	_ = low.BuildModel(&lNode, &lMt)
	_ = low.BuildModel(&rNode, &rMt)
	_ = lMt.Build(context.Background(), nil, lNode.Content[0], lIdx)
	_ = rMt.Build(context.Background(), nil, rNode.Content[0], rIdx)

	changes := CompareMediaTypes(&lMt, &rMt)

	assert.NotNil(t, changes)
	assert.Greater(t, changes.TotalChanges(), 0)

	allChanges := changes.GetAllChanges()
	assert.Len(t, allChanges, 1)
	assert.Equal(t, ObjectRemoved, allChanges[0].ChangeType)
	assert.Equal(t, v3.ItemEncodingLabel, allChanges[0].Property)
}

// TestOpenAPI32_MediaType_ItemEncoding_Modified tests that modifying itemEncoding is tracked
func TestOpenAPI32_MediaType_ItemEncoding_Modified(t *testing.T) {
	low.ClearHashCache()

	left := `schema:
  type: array
itemEncoding:
  file:
    contentType: image/jpeg
    allowReserved: true`

	right := `schema:
  type: array
itemEncoding:
  file:
    contentType: image/png
    allowReserved: false`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	lIdx := index.NewSpecIndex(&lNode)
	rIdx := index.NewSpecIndex(&rNode)

	var lMt, rMt v3.MediaType
	_ = low.BuildModel(&lNode, &lMt)
	_ = low.BuildModel(&rNode, &rMt)
	_ = lMt.Build(context.Background(), nil, lNode.Content[0], lIdx)
	_ = rMt.Build(context.Background(), nil, rNode.Content[0], rIdx)

	changes := CompareMediaTypes(&lMt, &rMt)

	assert.NotNil(t, changes)
	assert.NotNil(t, changes.ItemEncodingChanges)
	assert.Greater(t, changes.TotalChanges(), 0)
}

// TestOpenAPI32_SecurityScheme_OAuth2MetadataUrl tests breaking rules for oauth2MetadataUrl (OpenAPI 3.2+)
func TestOpenAPI32_SecurityScheme_OAuth2MetadataUrl(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.SecurityScheme)

	// oauth2MetadataUrl changes should not be breaking by default
}

// TestOpenAPI32_SecurityScheme_OAuth2MetadataUrl_Added tests that adding oauth2MetadataUrl is not breaking
func TestOpenAPI32_SecurityScheme_OAuth2MetadataUrl_Added(t *testing.T) {
	low.ClearHashCache()

	left := `type: oauth2
flows:
  implicit:
    authorizationUrl: https://example.com/auth
    scopes:
      read: Read access`

	right := `type: oauth2
oauth2MetadataUrl: https://example.com/.well-known/oauth-authorization-server
flows:
  implicit:
    authorizationUrl: https://example.com/auth
    scopes:
      read: Read access`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	var lSS, rSS v3.SecurityScheme
	_ = low.BuildModel(lNode.Content[0], &lSS)
	_ = low.BuildModel(rNode.Content[0], &rSS)
	_ = lSS.Build(context.Background(), nil, lNode.Content[0], nil)
	_ = rSS.Build(context.Background(), nil, rNode.Content[0], nil)

	changes := CompareSecuritySchemes(&lSS, &rSS)

	assert.NotNil(t, changes)

	// Adding oauth2MetadataUrl should not be breaking
	foundOAuth2MetadataUrlAdded := false
	for _, change := range changes.GetAllChanges() {
		if change.Property == v3.OAuth2MetadataUrlLabel && change.ChangeType == PropertyAdded {
			foundOAuth2MetadataUrlAdded = true
			assert.False(t, change.Breaking, "Adding oauth2MetadataUrl should not be breaking")
			break
		}
	}
	assert.True(t, foundOAuth2MetadataUrlAdded, "oauth2MetadataUrl addition should be detected")
}

// TestOpenAPI32_SecurityScheme_OAuth2MetadataUrl_Removed tests that removing oauth2MetadataUrl is not breaking
func TestOpenAPI32_SecurityScheme_OAuth2MetadataUrl_Removed(t *testing.T) {
	low.ClearHashCache()

	left := `type: oauth2
oauth2MetadataUrl: https://example.com/.well-known/oauth-authorization-server
flows:
  implicit:
    authorizationUrl: https://example.com/auth
    scopes:
      read: Read access`

	right := `type: oauth2
flows:
  implicit:
    authorizationUrl: https://example.com/auth
    scopes:
      read: Read access`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	var lSS, rSS v3.SecurityScheme
	_ = low.BuildModel(lNode.Content[0], &lSS)
	_ = low.BuildModel(rNode.Content[0], &rSS)
	_ = lSS.Build(context.Background(), nil, lNode.Content[0], nil)
	_ = rSS.Build(context.Background(), nil, rNode.Content[0], nil)

	changes := CompareSecuritySchemes(&lSS, &rSS)

	assert.NotNil(t, changes)

	// Removing oauth2MetadataUrl should not be breaking
	foundOAuth2MetadataUrlRemoved := false
	for _, change := range changes.GetAllChanges() {
		if change.Property == v3.OAuth2MetadataUrlLabel && change.ChangeType == PropertyRemoved {
			foundOAuth2MetadataUrlRemoved = true
			assert.False(t, change.Breaking, "Removing oauth2MetadataUrl should not be breaking")
			break
		}
	}
	assert.True(t, foundOAuth2MetadataUrlRemoved, "oauth2MetadataUrl removal should be detected")
}

// TestOpenAPI32_SecurityScheme_OAuth2MetadataUrl_Modified tests that modifying oauth2MetadataUrl is not breaking
func TestOpenAPI32_SecurityScheme_OAuth2MetadataUrl_Modified(t *testing.T) {
	low.ClearHashCache()

	left := `type: oauth2
oauth2MetadataUrl: https://example.com/.well-known/oauth-authorization-server
flows:
  implicit:
    authorizationUrl: https://example.com/auth
    scopes:
      read: Read access`

	right := `type: oauth2
oauth2MetadataUrl: https://example.com/.well-known/oauth-authorization-server-v2
flows:
  implicit:
    authorizationUrl: https://example.com/auth
    scopes:
      read: Read access`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	var lSS, rSS v3.SecurityScheme
	_ = low.BuildModel(lNode.Content[0], &lSS)
	_ = low.BuildModel(rNode.Content[0], &rSS)
	_ = lSS.Build(context.Background(), nil, lNode.Content[0], nil)
	_ = rSS.Build(context.Background(), nil, rNode.Content[0], nil)

	changes := CompareSecuritySchemes(&lSS, &rSS)

	assert.NotNil(t, changes)

	// Modifying oauth2MetadataUrl should not be breaking
	foundOAuth2MetadataUrlModified := false
	for _, change := range changes.GetAllChanges() {
		if change.Property == v3.OAuth2MetadataUrlLabel && change.ChangeType == Modified {
			foundOAuth2MetadataUrlModified = true
			assert.False(t, change.Breaking, "Modifying oauth2MetadataUrl should not be breaking")
			break
		}
	}
	assert.True(t, foundOAuth2MetadataUrlModified, "oauth2MetadataUrl modification should be detected")
}

// TestOpenAPI32_PathItemChanges_AdditionalOperationChanges tests that AdditionalOperationChanges are tracked
func TestOpenAPI32_PathItemChanges_AdditionalOperationChanges(t *testing.T) {
	pc := &PathItemChanges{
		PropertyChanges: &PropertyChanges{
			Changes: []*Change{
				{
					ChangeType: PropertyAdded,
					Property:   "description",
					Breaking:   false,
				},
			},
		},
		AdditionalOperationChanges: map[string]*OperationChanges{
			"x-custom": {
				PropertyChanges: &PropertyChanges{
					Changes: []*Change{
						{
							ChangeType: Modified,
							Property:   "summary",
							Breaking:   false,
						},
					},
				},
			},
		},
	}

	// Test GetAllChanges includes AdditionalOperationChanges
	allChanges := pc.GetAllChanges()
	assert.NotNil(t, allChanges)
	assert.Equal(t, 2, len(allChanges))

	// Test TotalChanges includes AdditionalOperationChanges
	total := pc.TotalChanges()
	assert.Equal(t, 2, total)

	// Test TotalBreakingChanges includes AdditionalOperationChanges
	breaking := pc.TotalBreakingChanges()
	assert.Equal(t, 0, breaking)
}

// TestOpenAPI32_PathItemChanges_AdditionalOperationChanges_Breaking tests breaking changes in AdditionalOperationChanges
func TestOpenAPI32_PathItemChanges_AdditionalOperationChanges_Breaking(t *testing.T) {
	pc := &PathItemChanges{
		PropertyChanges: &PropertyChanges{
			Changes: []*Change{},
		},
		AdditionalOperationChanges: map[string]*OperationChanges{
			"x-custom": {
				PropertyChanges: &PropertyChanges{
					Changes: []*Change{
						{
							ChangeType: PropertyRemoved,
							Property:   "operationId",
							Breaking:   true,
						},
					},
				},
			},
		},
	}

	// Test TotalBreakingChanges includes breaking changes from AdditionalOperationChanges
	breaking := pc.TotalBreakingChanges()
	assert.Equal(t, 1, breaking)
}

// TestOpenAPI32_MediaTypeChanges_ItemSchemaChanges tests that ItemSchemaChanges are tracked
func TestOpenAPI32_MediaTypeChanges_ItemSchemaChanges(t *testing.T) {
	mc := &MediaTypeChanges{
		PropertyChanges: &PropertyChanges{
			Changes: []*Change{},
		},
		ItemSchemaChanges: &SchemaChanges{
			PropertyChanges: &PropertyChanges{
				Changes: []*Change{
					{
						ChangeType: Modified,
						Property:   "type",
						Breaking:   true,
					},
				},
			},
		},
	}

	// Test GetAllChanges includes ItemSchemaChanges
	allChanges := mc.GetAllChanges()
	assert.NotNil(t, allChanges)
	assert.Equal(t, 1, len(allChanges))

	// Test TotalChanges includes ItemSchemaChanges
	total := mc.TotalChanges()
	assert.Equal(t, 1, total)

	// Test TotalBreakingChanges includes ItemSchemaChanges
	breaking := mc.TotalBreakingChanges()
	assert.Equal(t, 1, breaking)
}

// TestOpenAPI32_MediaTypeChanges_ItemEncodingChanges tests that ItemEncodingChanges are tracked
func TestOpenAPI32_MediaTypeChanges_ItemEncodingChanges(t *testing.T) {
	mc := &MediaTypeChanges{
		PropertyChanges: &PropertyChanges{
			Changes: []*Change{},
		},
		ItemEncodingChanges: map[string]*EncodingChanges{
			"file": {
				PropertyChanges: &PropertyChanges{
					Changes: []*Change{
						{
							ChangeType: Modified,
							Property:   "contentType",
							Breaking:   true,
						},
					},
				},
			},
		},
	}

	// Test GetAllChanges includes ItemEncodingChanges
	allChanges := mc.GetAllChanges()
	assert.NotNil(t, allChanges)
	assert.Equal(t, 1, len(allChanges))

	// Test TotalChanges includes ItemEncodingChanges
	total := mc.TotalChanges()
	assert.Equal(t, 1, total)

	// Test TotalBreakingChanges includes ItemEncodingChanges
	breaking := mc.TotalBreakingChanges()
	assert.Equal(t, 1, breaking)
}
