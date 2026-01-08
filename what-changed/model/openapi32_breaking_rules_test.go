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
		if change.Property == v3.QueryLabel && change.ChangeType == PropertyAdded {
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
		if change.Property == v3.QueryLabel && change.ChangeType == PropertyRemoved {
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
  operationId: queryResources`

	right := `query:
  summary: Query resources v2
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
	assert.Greater(t, changes.TotalChanges(), 0)
}

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

	foundAdditionalOpsAdded := false
	for _, change := range changes.GetAllChanges() {
		if change.Property == v3.AdditionalOperationsLabel && change.ChangeType == PropertyAdded {
			foundAdditionalOpsAdded = true
			break
		}
	}
	assert.True(t, foundAdditionalOpsAdded, "AdditionalOperations should be detected as added")
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

	foundAdditionalOpsRemoved := false
	for _, change := range changes.GetAllChanges() {
		if change.Property == v3.AdditionalOperationsLabel && change.ChangeType == PropertyRemoved {
			foundAdditionalOpsRemoved = true
			assert.True(t, change.Breaking)
			break
		}
	}
	assert.True(t, foundAdditionalOpsRemoved, "AdditionalOperations should be detected as removed")
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
	assert.Greater(t, len(changes.AdditionalOperationChanges), 0)
}

func TestOpenAPI32_PathItem_AdditionalOperationAddedToExisting(t *testing.T) {
	low.ClearHashCache()

	left := `additionalOperations:
  customOp1:
    summary: Custom operation 1`

	right := `additionalOperations:
  customOp1:
    summary: Custom operation 1
  customOp2:
    summary: Custom operation 2`

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

func TestOpenAPI32_PathItem_AdditionalOperationRemovedFromExisting(t *testing.T) {
	low.ClearHashCache()

	left := `additionalOperations:
  customOp1:
    summary: Custom operation 1
  customOp2:
    summary: Custom operation 2`

	right := `additionalOperations:
  customOp1:
    summary: Custom operation 1`

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
	assert.Greater(t, changes.TotalChanges(), 0)
}

func TestOpenAPI32_SecurityScheme_OAuth2MetadataUrlAdded(t *testing.T) {
	low.ClearHashCache()

	left := `type: oauth2
flows:
  implicit:
    authorizationUrl: https://example.com/oauth/authorize
    scopes:
      read: Read access`

	right := `type: oauth2
oauth2MetadataUrl: https://example.com/.well-known/oauth-authorization-server
flows:
  implicit:
    authorizationUrl: https://example.com/oauth/authorize
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

func TestOpenAPI32_SecurityScheme_OAuth2MetadataUrlRemoved(t *testing.T) {
	low.ClearHashCache()

	left := `type: oauth2
oauth2MetadataUrl: https://example.com/.well-known/oauth-authorization-server
flows:
  implicit:
    authorizationUrl: https://example.com/oauth/authorize
    scopes:
      read: Read access`

	right := `type: oauth2
flows:
  implicit:
    authorizationUrl: https://example.com/oauth/authorize
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
    authorizationUrl: https://example.com/oauth/authorize
    scopes:
      read: Read access`

	right := `type: oauth2
oauth2MetadataUrl: https://example.com/.well-known/oauth-authorization-server-v2
flows:
  implicit:
    authorizationUrl: https://example.com/oauth/authorize
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

func TestOpenAPI32_AllPropertiesNoChanges(t *testing.T) {
	low.ClearHashCache()

	pathItemYaml := `get:
  summary: Get resource
query:
  summary: Query resources
additionalOperations:
  customOp:
    summary: Custom operation`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(pathItemYaml), &lNode)
	_ = yaml.Unmarshal([]byte(pathItemYaml), &rNode)

	lIdx := index.NewSpecIndexWithConfig(&lNode, index.CreateOpenAPIIndexConfig())
	rIdx := index.NewSpecIndexWithConfig(&rNode, index.CreateOpenAPIIndexConfig())
	ctx := context.Background()

	var lPath, rPath v3.PathItem
	_ = low.BuildModel(&lNode, &lPath)
	_ = low.BuildModel(&rNode, &rPath)

	_ = lPath.Build(ctx, nil, lNode.Content[0], lIdx)
	_ = rPath.Build(ctx, nil, rNode.Content[0], rIdx)

	changes := ComparePathItems(&lPath, &rPath)

	assert.Nil(t, changes)
}

func TestOpenAPI32_MediaType_NoChanges(t *testing.T) {
	low.ClearHashCache()

	mediaTypeYaml := `schema:
  type: array
itemSchema:
  type: object
itemEncoding:
  file:
    contentType: image/jpeg`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(mediaTypeYaml), &lNode)
	_ = yaml.Unmarshal([]byte(mediaTypeYaml), &rNode)

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

func TestOpenAPI32_SecurityScheme_NoChanges(t *testing.T) {
	low.ClearHashCache()

	securitySchemeYaml := `type: oauth2
oauth2MetadataUrl: https://example.com/.well-known/oauth-authorization-server
flows:
  implicit:
    authorizationUrl: https://example.com/oauth/authorize
    scopes:
      read: Read access`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(securitySchemeYaml), &lNode)
	_ = yaml.Unmarshal([]byte(securitySchemeYaml), &rNode)

	var lDoc v3.SecurityScheme
	var rDoc v3.SecurityScheme
	_ = low.BuildModel(lNode.Content[0], &lDoc)
	_ = low.BuildModel(rNode.Content[0], &rDoc)
	_ = lDoc.Build(context.Background(), nil, lNode.Content[0], nil)
	_ = rDoc.Build(context.Background(), nil, rNode.Content[0], nil)

	changes := CompareSecuritySchemes(&lDoc, &rDoc)

	assert.Nil(t, changes)
}
