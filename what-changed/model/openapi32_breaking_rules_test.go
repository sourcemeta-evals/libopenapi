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

func TestComparePathItems_QueryAdded_OpenAPI32(t *testing.T) {
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
	assert.Nil(t, changes.QueryChanges)
	assert.Equal(t, 1, changes.TotalChanges())

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

func TestComparePathItems_QueryRemoved_OpenAPI32(t *testing.T) {
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
	assert.Nil(t, changes.QueryChanges)
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

func TestComparePathItems_QueryModified_OpenAPI32(t *testing.T) {
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

func TestComparePathItems_AdditionalOperationsAdded_OpenAPI32(t *testing.T) {
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

	assert.NotNil(t, changes)
	assert.True(t, changes.TotalChanges() >= 1)
}

func TestComparePathItems_AdditionalOperationsRemoved_OpenAPI32(t *testing.T) {
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

	assert.NotNil(t, changes)
	assert.True(t, changes.TotalChanges() >= 1)
}

func TestComparePathItems_AdditionalOperationsModified_OpenAPI32(t *testing.T) {
	low.ClearHashCache()

	left := `x-custom-op:
  summary: Custom operation v1
  operationId: customOpV1`

	right := `x-custom-op:
  summary: Custom operation v2
  operationId: customOpV2`

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
	assert.True(t, changes.TotalChanges() >= 1)
}

func TestCompareMediaTypes_ItemSchemaAdded_OpenAPI32(t *testing.T) {
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

func TestCompareMediaTypes_ItemSchemaRemoved_OpenAPI32(t *testing.T) {
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

func TestCompareMediaTypes_ItemSchemaModified_OpenAPI32(t *testing.T) {
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

func TestCompareMediaTypes_ItemEncodingAdded_OpenAPI32(t *testing.T) {
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

func TestCompareMediaTypes_ItemEncodingRemoved_OpenAPI32(t *testing.T) {
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

func TestCompareMediaTypes_ItemEncodingModified_OpenAPI32(t *testing.T) {
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

func TestCompareSecuritySchemes_OAuth2MetadataUrlAdded_OpenAPI32(t *testing.T) {
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
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 0, changes.TotalBreakingChanges())

	foundOAuth2MetadataUrlAdded := false
	for _, change := range changes.GetAllChanges() {
		if change.Property == v3.OAuth2MetadataUrlLabel && change.ChangeType == PropertyAdded {
			foundOAuth2MetadataUrlAdded = true
			assert.False(t, change.Breaking)
			break
		}
	}
	assert.True(t, foundOAuth2MetadataUrlAdded, "oauth2MetadataUrl addition should be detected")
}

func TestCompareSecuritySchemes_OAuth2MetadataUrlRemoved_OpenAPI32(t *testing.T) {
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
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 0, changes.TotalBreakingChanges())

	foundOAuth2MetadataUrlRemoved := false
	for _, change := range changes.GetAllChanges() {
		if change.Property == v3.OAuth2MetadataUrlLabel && change.ChangeType == PropertyRemoved {
			foundOAuth2MetadataUrlRemoved = true
			assert.False(t, change.Breaking)
			break
		}
	}
	assert.True(t, foundOAuth2MetadataUrlRemoved, "oauth2MetadataUrl removal should be detected")
}

func TestCompareSecuritySchemes_OAuth2MetadataUrlModified_OpenAPI32(t *testing.T) {
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
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 0, changes.TotalBreakingChanges())

	foundOAuth2MetadataUrlModified := false
	for _, change := range changes.GetAllChanges() {
		if change.Property == v3.OAuth2MetadataUrlLabel && change.ChangeType == Modified {
			foundOAuth2MetadataUrlModified = true
			assert.False(t, change.Breaking)
			assert.Equal(t, "https://example.com/.well-known/oauth-authorization-server", change.Original)
			assert.Equal(t, "https://example.com/.well-known/oauth-authorization-server-v2", change.New)
			break
		}
	}
	assert.True(t, foundOAuth2MetadataUrlModified, "oauth2MetadataUrl modification should be detected")
}

func TestPathItemChanges_AdditionalOperationChanges_GetAllChanges(t *testing.T) {
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
			"x-custom-op": {
				PropertyChanges: &PropertyChanges{
					Changes: []*Change{
						{
							ChangeType: Modified,
							Property:   "summary",
							Breaking:   false,
						},
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

	allChanges := pc.GetAllChanges()
	assert.NotNil(t, allChanges)
	assert.Equal(t, 3, len(allChanges))

	additionalOpChangeCount := 0
	for _, change := range allChanges {
		if (change.Property == "summary" && change.ChangeType == Modified) ||
			(change.Property == "operationId" && change.ChangeType == PropertyRemoved) {
			additionalOpChangeCount++
		}
	}
	assert.Equal(t, 2, additionalOpChangeCount, "Both AdditionalOperationChanges should be in the result")

	total := pc.TotalChanges()
	assert.Equal(t, 3, total)

	breaking := pc.TotalBreakingChanges()
	assert.Equal(t, 1, breaking)
}

func TestPathItemChanges_AdditionalOperationChanges_Nil(t *testing.T) {
	pc := &PathItemChanges{
		PropertyChanges: &PropertyChanges{
			Changes: []*Change{
				{
					ChangeType: PropertyAdded,
					Property:   "summary",
					Breaking:   false,
				},
			},
		},
		AdditionalOperationChanges: nil,
	}

	allChanges := pc.GetAllChanges()
	assert.NotNil(t, allChanges)
	assert.Equal(t, 1, len(allChanges))

	total := pc.TotalChanges()
	assert.Equal(t, 1, total)

	breaking := pc.TotalBreakingChanges()
	assert.Equal(t, 0, breaking)
}

func TestMediaTypeChanges_ItemSchemaChanges_GetAllChanges(t *testing.T) {
	mc := &MediaTypeChanges{
		PropertyChanges: &PropertyChanges{
			Changes: []*Change{
				{
					ChangeType: PropertyAdded,
					Property:   "example",
					Breaking:   false,
				},
			},
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

	allChanges := mc.GetAllChanges()
	assert.NotNil(t, allChanges)
	assert.Equal(t, 2, len(allChanges))

	total := mc.TotalChanges()
	assert.Equal(t, 2, total)

	breaking := mc.TotalBreakingChanges()
	assert.Equal(t, 1, breaking)
}

func TestMediaTypeChanges_ItemEncodingChanges_GetAllChanges(t *testing.T) {
	mc := &MediaTypeChanges{
		PropertyChanges: &PropertyChanges{
			Changes: []*Change{
				{
					ChangeType: PropertyAdded,
					Property:   "example",
					Breaking:   false,
				},
			},
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

	allChanges := mc.GetAllChanges()
	assert.NotNil(t, allChanges)
	assert.Equal(t, 2, len(allChanges))

	total := mc.TotalChanges()
	assert.Equal(t, 2, total)

	breaking := mc.TotalBreakingChanges()
	assert.Equal(t, 1, breaking)
}

func TestMediaTypeChanges_BothItemSchemaAndItemEncoding(t *testing.T) {
	mc := &MediaTypeChanges{
		PropertyChanges: &PropertyChanges{
			Changes: []*Change{
				{
					ChangeType: PropertyAdded,
					Property:   "example",
					Breaking:   false,
				},
			},
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

	allChanges := mc.GetAllChanges()
	assert.NotNil(t, allChanges)
	assert.Equal(t, 3, len(allChanges))

	total := mc.TotalChanges()
	assert.Equal(t, 3, total)

	breaking := mc.TotalBreakingChanges()
	assert.Equal(t, 2, breaking)
}
