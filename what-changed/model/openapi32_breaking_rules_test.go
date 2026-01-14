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

// Tests for OpenAPI 3.2 specific properties in the what-changed module:
// - query (Path Item)
// - additionalOperations (Path Item)
// - itemSchema (Media Type)
// - itemEncoding (Media Type)
// - oauth2MetadataUrl (Security Scheme)

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
	assert.Equal(t, 0, changes.TotalBreakingChanges())

	foundQueryAdded := false
	for _, change := range changes.GetAllChanges() {
		if change.Property == "query" && change.ChangeType == PropertyAdded {
			foundQueryAdded = true
			assert.False(t, change.Breaking)
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

	foundQueryRemoved := false
	for _, change := range changes.GetAllChanges() {
		if change.Property == "query" && change.ChangeType == PropertyRemoved {
			foundQueryRemoved = true
			assert.True(t, change.Breaking)
			break
		}
	}
	assert.True(t, foundQueryRemoved, "Query operation should be detected as removed")
}

// TestOpenAPI32_PathItem_Query_Modified tests that modifying a query operation tracks changes
func TestOpenAPI32_PathItem_Query_Modified(t *testing.T) {
	low.ClearHashCache()

	left := `query:
  summary: Query resources
  operationId: queryResources`

	right := `query:
  summary: Query resources updated
  operationId: queryResourcesV2`

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
	assert.NotNil(t, changes.QueryChanges)
	assert.Greater(t, changes.QueryChanges.TotalChanges(), 0)
}

// TestOpenAPI32_PathItem_AdditionalOperations_Added tests that adding additional operations is not breaking
func TestOpenAPI32_PathItem_AdditionalOperations_Added(t *testing.T) {
	low.ClearHashCache()

	left := `get:
  summary: Get resource`

	right := `get:
  summary: Get resource
additionalOperations:
  customOp:
    summary: Custom operation`

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
	assert.Equal(t, 0, changes.TotalBreakingChanges())
}

// TestOpenAPI32_PathItem_AdditionalOperations_Removed tests that removing additional operations is breaking
func TestOpenAPI32_PathItem_AdditionalOperations_Removed(t *testing.T) {
	low.ClearHashCache()

	left := `get:
  summary: Get resource
additionalOperations:
  customOp:
    summary: Custom operation`

	right := `get:
  summary: Get resource`

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
}

// TestOpenAPI32_PathItem_AdditionalOperations_Modified tests that modifying additional operations tracks changes
func TestOpenAPI32_PathItem_AdditionalOperations_Modified(t *testing.T) {
	low.ClearHashCache()

	left := `additionalOperations:
  customOp:
    summary: Custom operation`

	right := `additionalOperations:
  customOp:
    summary: Custom operation updated`

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
	assert.Greater(t, changes.TotalChanges(), 0)
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
}

// TestOpenAPI32_MediaType_ItemSchema_Modified tests that modifying itemSchema tracks changes
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

// TestOpenAPI32_MediaType_ItemEncoding_Added tests that adding itemEncoding is tracked
func TestOpenAPI32_MediaType_ItemEncoding_Added(t *testing.T) {
	low.ClearHashCache()

	left := `schema:
  type: array`

	right := `schema:
  type: array
itemEncoding:
  file:
    contentType: image/jpeg`

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
    contentType: image/jpeg`

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

// TestOpenAPI32_MediaType_ItemEncoding_Modified tests that modifying itemEncoding tracks changes
func TestOpenAPI32_MediaType_ItemEncoding_Modified(t *testing.T) {
	low.ClearHashCache()

	left := `schema:
  type: array
itemEncoding:
  file:
    contentType: image/jpeg`

	right := `schema:
  type: array
itemEncoding:
  file:
    contentType: image/png`

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
	assert.NotNil(t, changes.ItemEncodingChanges)
}

// TestOpenAPI32_SecurityScheme_OAuth2MetadataUrl_Added tests that adding oauth2MetadataUrl is not breaking
func TestOpenAPI32_SecurityScheme_OAuth2MetadataUrl_Added(t *testing.T) {
	low.ClearHashCache()

	left := `type: oauth2
flows:
  implicit:
    authorizationUrl: https://example.com/auth`

	right := `type: oauth2
flows:
  implicit:
    authorizationUrl: https://example.com/auth
oauth2MetadataUrl: https://example.com/.well-known/oauth-authorization-server`

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
	assert.Greater(t, changes.TotalChanges(), 0)

	foundOAuth2MetadataUrlAdded := false
	for _, change := range changes.GetAllChanges() {
		if change.Property == v3.OAuth2MetadataUrlLabel && change.ChangeType == PropertyAdded {
			foundOAuth2MetadataUrlAdded = true
			assert.False(t, change.Breaking)
			break
		}
	}
	assert.True(t, foundOAuth2MetadataUrlAdded, "oauth2MetadataUrl should be detected as added")
}

// TestOpenAPI32_SecurityScheme_OAuth2MetadataUrl_Removed tests that removing oauth2MetadataUrl is not breaking
func TestOpenAPI32_SecurityScheme_OAuth2MetadataUrl_Removed(t *testing.T) {
	low.ClearHashCache()

	left := `type: oauth2
flows:
  implicit:
    authorizationUrl: https://example.com/auth
oauth2MetadataUrl: https://example.com/.well-known/oauth-authorization-server`

	right := `type: oauth2
flows:
  implicit:
    authorizationUrl: https://example.com/auth`

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
	assert.Greater(t, changes.TotalChanges(), 0)

	foundOAuth2MetadataUrlRemoved := false
	for _, change := range changes.GetAllChanges() {
		if change.Property == v3.OAuth2MetadataUrlLabel && change.ChangeType == PropertyRemoved {
			foundOAuth2MetadataUrlRemoved = true
			assert.False(t, change.Breaking)
			break
		}
	}
	assert.True(t, foundOAuth2MetadataUrlRemoved, "oauth2MetadataUrl should be detected as removed")
}

// TestOpenAPI32_SecurityScheme_OAuth2MetadataUrl_Modified tests that modifying oauth2MetadataUrl is not breaking
func TestOpenAPI32_SecurityScheme_OAuth2MetadataUrl_Modified(t *testing.T) {
	low.ClearHashCache()

	left := `type: oauth2
flows:
  implicit:
    authorizationUrl: https://example.com/auth
oauth2MetadataUrl: https://example.com/.well-known/oauth-authorization-server`

	right := `type: oauth2
flows:
  implicit:
    authorizationUrl: https://example.com/auth
oauth2MetadataUrl: https://example.com/.well-known/oauth-authorization-server-v2`

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
	assert.Greater(t, changes.TotalChanges(), 0)

	foundOAuth2MetadataUrlModified := false
	for _, change := range changes.GetAllChanges() {
		if change.Property == v3.OAuth2MetadataUrlLabel && change.ChangeType == Modified {
			foundOAuth2MetadataUrlModified = true
			assert.False(t, change.Breaking)
			break
		}
	}
	assert.True(t, foundOAuth2MetadataUrlModified, "oauth2MetadataUrl should be detected as modified")
}

// TestOpenAPI32_AllPropertiesTogether tests all OpenAPI 3.2 properties in a combined scenario
func TestOpenAPI32_AllPropertiesTogether(t *testing.T) {
	low.ClearHashCache()

	left := `get:
  summary: Get resource
query:
  summary: Query resources
additionalOperations:
  customOp:
    summary: Custom operation`

	right := `get:
  summary: Get resource updated
query:
  summary: Query resources updated
additionalOperations:
  customOp:
    summary: Custom operation updated`

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
	assert.Greater(t, changes.TotalChanges(), 0)
	assert.NotNil(t, changes.GetChanges)
	assert.NotNil(t, changes.QueryChanges)
	assert.NotNil(t, changes.AdditionalOperationChanges)
}
