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

// Tests for OpenAPI 3.2 properties in the what-changed module:
// - query (Path Item)
// - additionalOperations (Path Item)
// - itemSchema (Media Type)
// - itemEncoding (Media Type)
// - oauth2MetadataUrl (Security Scheme)

// ============================================================================
// Query (Path Item) Tests
// ============================================================================

func TestOpenAPI32_PathItem_QueryAdded(t *testing.T) {
	low.ClearHashCache()

	left := `get:
  summary: Get resource
  operationId: getResource`

	right := `get:
  summary: Get resource
  operationId: getResource
query:
  summary: Query resources
  operationId: queryResources
  requestBody:
    required: true
    content:
      application/json:
        schema:
          type: object`

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
			break
		}
	}
	assert.True(t, foundQueryAdded, "Query operation should be detected as added")
}

func TestOpenAPI32_PathItem_QueryRemoved(t *testing.T) {
	low.ClearHashCache()

	left := `get:
  summary: Get resource
  operationId: getResource
query:
  summary: Query resources
  operationId: queryResources
  requestBody:
    required: true
    content:
      application/json:
        schema:
          type: object`

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

func TestOpenAPI32_PathItem_QueryModified(t *testing.T) {
	low.ClearHashCache()

	left := `query:
  summary: Query resources
  operationId: queryResources
  responses:
    '200':
      description: OK`

	right := `query:
  summary: Query resources with filters
  operationId: queryResourcesV2
  requestBody:
    required: true
    content:
      application/json:
        schema:
          type: object
  responses:
    '200':
      description: Query results`

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
	assert.True(t, changes.QueryChanges.TotalChanges() > 0)

	foundSummaryChange := false
	foundOperationIdChange := false
	foundRequestBodyAdded := false

	for _, change := range changes.QueryChanges.GetAllChanges() {
		if change.Property == "summary" {
			foundSummaryChange = true
			assert.Equal(t, "Query resources", change.Original)
			assert.Equal(t, "Query resources with filters", change.New)
		}
		if change.Property == "operationId" {
			foundOperationIdChange = true
			assert.Equal(t, "queryResources", change.Original)
			assert.Equal(t, "queryResourcesV2", change.New)
		}
		if change.Property == "requestBody" && change.ChangeType == PropertyAdded {
			foundRequestBodyAdded = true
		}
	}

	assert.True(t, foundSummaryChange, "Summary change should be detected")
	assert.True(t, foundOperationIdChange, "OperationId change should be detected")
	assert.True(t, foundRequestBodyAdded, "RequestBody addition should be detected")
}

// ============================================================================
// AdditionalOperations (Path Item) Tests
// ============================================================================

func TestOpenAPI32_PathItem_AdditionalOperationsAdded(t *testing.T) {
	low.ClearHashCache()

	left := `get:
  summary: Get resource`

	right := `get:
  summary: Get resource
additionalOperations:
  customOp:
    summary: Custom operation
    operationId: customOperation`

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

func TestOpenAPI32_PathItem_AdditionalOperationsRemoved(t *testing.T) {
	low.ClearHashCache()

	left := `get:
  summary: Get resource
additionalOperations:
  customOp:
    summary: Custom operation
    operationId: customOperation`

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

func TestOpenAPI32_PathItem_AdditionalOperationsModified(t *testing.T) {
	low.ClearHashCache()

	left := `additionalOperations:
  customOp:
    summary: Custom operation v1
    operationId: customOperationV1`

	right := `additionalOperations:
  customOp:
    summary: Custom operation v2
    operationId: customOperationV2`

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
	assert.NotNil(t, changes.AdditionalOperationChanges)
	assert.NotNil(t, changes.AdditionalOperationChanges["customOp"])
	assert.True(t, changes.TotalChanges() > 0)
}

func TestOpenAPI32_PathItem_MultipleAdditionalOperations(t *testing.T) {
	low.ClearHashCache()

	left := `additionalOperations:
  op1:
    summary: Operation 1
  op2:
    summary: Operation 2`

	right := `additionalOperations:
  op1:
    summary: Operation 1 modified
  op3:
    summary: Operation 3`

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
	assert.True(t, changes.TotalChanges() >= 2)
	assert.True(t, changes.TotalBreakingChanges() >= 1)
}

// ============================================================================
// ItemSchema (Media Type) Tests
// ============================================================================

func TestOpenAPI32_MediaType_ItemSchemaAdded(t *testing.T) {
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
	assert.Nil(t, changes.ItemSchemaChanges)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 1, changes.TotalBreakingChanges())

	allChanges := changes.GetAllChanges()
	assert.Len(t, allChanges, 1)
	assert.Equal(t, ObjectAdded, allChanges[0].ChangeType)
	assert.Equal(t, v3.ItemSchemaLabel, allChanges[0].Property)
}

func TestOpenAPI32_MediaType_ItemSchemaRemoved(t *testing.T) {
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
	assert.Nil(t, changes.ItemSchemaChanges)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 1, changes.TotalBreakingChanges())

	allChanges := changes.GetAllChanges()
	assert.Len(t, allChanges, 1)
	assert.Equal(t, ObjectRemoved, allChanges[0].ChangeType)
	assert.Equal(t, v3.ItemSchemaLabel, allChanges[0].Property)
}

func TestOpenAPI32_MediaType_ItemSchemaModified(t *testing.T) {
	low.ClearHashCache()

	left := `schema:
  type: array
itemSchema:
  type: object
  properties:
    id:
      type: string
    name:
      type: string`

	right := `schema:
  type: array
itemSchema:
  type: object
  properties:
    id:
      type: integer
    name:
      type: string
    age:
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

// ============================================================================
// ItemEncoding (Media Type) Tests
// ============================================================================

func TestOpenAPI32_MediaType_ItemEncodingAdded(t *testing.T) {
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

func TestOpenAPI32_MediaType_ItemEncodingRemoved(t *testing.T) {
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

func TestOpenAPI32_MediaType_ItemEncodingModified(t *testing.T) {
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
	assert.NotNil(t, changes.ItemEncodingChanges["file"])
	assert.Greater(t, changes.TotalChanges(), 0)
}

func TestOpenAPI32_MediaType_BothItemSchemaAndItemEncoding(t *testing.T) {
	low.ClearHashCache()

	left := `schema:
  type: array
itemSchema:
  type: string
itemEncoding:
  data:
    contentType: text/plain`

	right := `schema:
  type: array
itemSchema:
  type: integer
itemEncoding:
  data:
    contentType: application/json`

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
	assert.NotNil(t, changes.ItemEncodingChanges)
	assert.Greater(t, changes.TotalChanges(), 1)
}

// ============================================================================
// OAuth2MetadataUrl (Security Scheme) Tests
// ============================================================================

func TestOpenAPI32_SecurityScheme_OAuth2MetadataUrlAdded(t *testing.T) {
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

	var lDoc v3.SecurityScheme
	var rDoc v3.SecurityScheme
	_ = low.BuildModel(lNode.Content[0], &lDoc)
	_ = low.BuildModel(rNode.Content[0], &rDoc)
	_ = lDoc.Build(context.Background(), nil, lNode.Content[0], nil)
	_ = rDoc.Build(context.Background(), nil, rNode.Content[0], nil)

	changes := CompareSecuritySchemes(&lDoc, &rDoc)

	assert.NotNil(t, changes)
	assert.Greater(t, changes.TotalChanges(), 0)

	foundOAuth2MetadataUrlAdded := false
	for _, change := range changes.GetAllChanges() {
		if change.Property == v3.OAuth2MetadataUrlLabel && change.ChangeType == PropertyAdded {
			foundOAuth2MetadataUrlAdded = true
			break
		}
	}
	assert.True(t, foundOAuth2MetadataUrlAdded, "oauth2MetadataUrl should be detected as added")
}

func TestOpenAPI32_SecurityScheme_OAuth2MetadataUrlRemoved(t *testing.T) {
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

	var lDoc v3.SecurityScheme
	var rDoc v3.SecurityScheme
	_ = low.BuildModel(lNode.Content[0], &lDoc)
	_ = low.BuildModel(rNode.Content[0], &rDoc)
	_ = lDoc.Build(context.Background(), nil, lNode.Content[0], nil)
	_ = rDoc.Build(context.Background(), nil, rNode.Content[0], nil)

	changes := CompareSecuritySchemes(&lDoc, &rDoc)

	assert.NotNil(t, changes)
	assert.Greater(t, changes.TotalChanges(), 0)

	foundOAuth2MetadataUrlRemoved := false
	for _, change := range changes.GetAllChanges() {
		if change.Property == v3.OAuth2MetadataUrlLabel && change.ChangeType == PropertyRemoved {
			foundOAuth2MetadataUrlRemoved = true
			break
		}
	}
	assert.True(t, foundOAuth2MetadataUrlRemoved, "oauth2MetadataUrl should be detected as removed")
}

func TestOpenAPI32_SecurityScheme_OAuth2MetadataUrlModified(t *testing.T) {
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

	var lDoc v3.SecurityScheme
	var rDoc v3.SecurityScheme
	_ = low.BuildModel(lNode.Content[0], &lDoc)
	_ = low.BuildModel(rNode.Content[0], &rDoc)
	_ = lDoc.Build(context.Background(), nil, lNode.Content[0], nil)
	_ = rDoc.Build(context.Background(), nil, rNode.Content[0], nil)

	changes := CompareSecuritySchemes(&lDoc, &rDoc)

	assert.NotNil(t, changes)
	assert.Greater(t, changes.TotalChanges(), 0)

	foundOAuth2MetadataUrlModified := false
	for _, change := range changes.GetAllChanges() {
		if change.Property == v3.OAuth2MetadataUrlLabel && change.ChangeType == Modified {
			foundOAuth2MetadataUrlModified = true
			assert.Equal(t, "https://example.com/.well-known/oauth-authorization-server", change.Original)
			assert.Equal(t, "https://example.com/.well-known/oauth-authorization-server-v2", change.New)
			break
		}
	}
	assert.True(t, foundOAuth2MetadataUrlModified, "oauth2MetadataUrl should be detected as modified")
}

func TestOpenAPI32_SecurityScheme_OAuth2MetadataUrlNoChange(t *testing.T) {
	low.ClearHashCache()

	left := `type: oauth2
oauth2MetadataUrl: https://example.com/.well-known/oauth-authorization-server
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

	var lDoc v3.SecurityScheme
	var rDoc v3.SecurityScheme
	_ = low.BuildModel(lNode.Content[0], &lDoc)
	_ = low.BuildModel(rNode.Content[0], &rDoc)
	_ = lDoc.Build(context.Background(), nil, lNode.Content[0], nil)
	_ = rDoc.Build(context.Background(), nil, rNode.Content[0], nil)

	changes := CompareSecuritySchemes(&lDoc, &rDoc)

	assert.Nil(t, changes)
}

// ============================================================================
// Combined OpenAPI 3.2 Tests
// ============================================================================

func TestOpenAPI32_AllNewProperties_NoChanges(t *testing.T) {
	low.ClearHashCache()

	left := `get:
  summary: Get resource
query:
  summary: Query resources
additionalOperations:
  custom:
    summary: Custom op`

	right := `get:
  summary: Get resource
query:
  summary: Query resources
additionalOperations:
  custom:
    summary: Custom op`

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

	assert.Equal(t, 0, len(changes.GetAllChanges()))
}

func TestOpenAPI32_MediaType_NoChangesWithItemFields(t *testing.T) {
	low.ClearHashCache()

	left := `schema:
  type: array
itemSchema:
  type: string
itemEncoding:
  data:
    contentType: text/plain`

	right := `schema:
  type: array
itemSchema:
  type: string
itemEncoding:
  data:
    contentType: text/plain`

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

	assert.Nil(t, changes)
}
