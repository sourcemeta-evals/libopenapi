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

func TestDefaultBreakingRules_PathItem_Query(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.PathItem)
	assert.NotNil(t, config.PathItem.Query)

	assert.False(t, *config.PathItem.Query.Added)
	assert.False(t, *config.PathItem.Query.Modified)
	assert.True(t, *config.PathItem.Query.Removed)
}

func TestDefaultBreakingRules_PathItem_AdditionalOperations(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.PathItem)
	assert.NotNil(t, config.PathItem.AdditionalOperations)

	assert.False(t, *config.PathItem.AdditionalOperations.Added)
	assert.False(t, *config.PathItem.AdditionalOperations.Modified)
	assert.True(t, *config.PathItem.AdditionalOperations.Removed)
}

func TestDefaultBreakingRules_MediaType_ItemSchema(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.MediaType)
	assert.NotNil(t, config.MediaType.ItemSchema)

	assert.True(t, *config.MediaType.ItemSchema.Added)
	assert.False(t, *config.MediaType.ItemSchema.Modified)
	assert.True(t, *config.MediaType.ItemSchema.Removed)
}

func TestDefaultBreakingRules_MediaType_ItemEncoding(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.MediaType)
	assert.NotNil(t, config.MediaType.ItemEncoding)

	assert.True(t, *config.MediaType.ItemEncoding.Added)
	assert.False(t, *config.MediaType.ItemEncoding.Modified)
	assert.True(t, *config.MediaType.ItemEncoding.Removed)
}

func TestDefaultBreakingRules_SecurityScheme_OAuth2MetadataUrl(t *testing.T) {
	config := GenerateDefaultBreakingRules()
	assert.NotNil(t, config.SecurityScheme)
	assert.NotNil(t, config.SecurityScheme.OAuth2MetadataUrl)

	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Added)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Modified)
	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Removed)
}

func TestGetRule_OpenAPI32_PathItemProperties(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	queryRule := config.GetRule("pathItem", "query")
	assert.NotNil(t, queryRule, "expected rule for pathItem/query")
	assert.False(t, *queryRule.Added)
	assert.True(t, *queryRule.Removed)

	additionalOpsRule := config.GetRule("pathItem", "additionalOperations")
	assert.NotNil(t, additionalOpsRule, "expected rule for pathItem/additionalOperations")
	assert.False(t, *additionalOpsRule.Added)
	assert.True(t, *additionalOpsRule.Removed)
}

func TestGetRule_OpenAPI32_MediaTypeProperties(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	itemSchemaRule := config.GetRule("mediaType", "itemSchema")
	assert.NotNil(t, itemSchemaRule, "expected rule for mediaType/itemSchema")
	assert.True(t, *itemSchemaRule.Added)
	assert.True(t, *itemSchemaRule.Removed)

	itemEncodingRule := config.GetRule("mediaType", "itemEncoding")
	assert.NotNil(t, itemEncodingRule, "expected rule for mediaType/itemEncoding")
	assert.True(t, *itemEncodingRule.Added)
	assert.True(t, *itemEncodingRule.Removed)
}

func TestGetRule_OpenAPI32_SecuritySchemeProperties(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	oauth2MetadataUrlRule := config.GetRule("securityScheme", "oauth2MetadataUrl")
	assert.NotNil(t, oauth2MetadataUrlRule, "expected rule for securityScheme/oauth2MetadataUrl")
	assert.False(t, *oauth2MetadataUrlRule.Added)
	assert.False(t, *oauth2MetadataUrlRule.Modified)
	assert.False(t, *oauth2MetadataUrlRule.Removed)
}

func TestIsBreaking_OpenAPI32_PathItemQuery(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	assert.False(t, config.IsBreaking("pathItem", "query", ChangeTypeAdded))
	assert.False(t, config.IsBreaking("pathItem", "query", ChangeTypeModified))
	assert.True(t, config.IsBreaking("pathItem", "query", ChangeTypeRemoved))
}

func TestIsBreaking_OpenAPI32_PathItemAdditionalOperations(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	assert.False(t, config.IsBreaking("pathItem", "additionalOperations", ChangeTypeAdded))
	assert.False(t, config.IsBreaking("pathItem", "additionalOperations", ChangeTypeModified))
	assert.True(t, config.IsBreaking("pathItem", "additionalOperations", ChangeTypeRemoved))
}

func TestIsBreaking_OpenAPI32_MediaTypeItemSchema(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	assert.True(t, config.IsBreaking("mediaType", "itemSchema", ChangeTypeAdded))
	assert.False(t, config.IsBreaking("mediaType", "itemSchema", ChangeTypeModified))
	assert.True(t, config.IsBreaking("mediaType", "itemSchema", ChangeTypeRemoved))
}

func TestIsBreaking_OpenAPI32_MediaTypeItemEncoding(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	assert.True(t, config.IsBreaking("mediaType", "itemEncoding", ChangeTypeAdded))
	assert.False(t, config.IsBreaking("mediaType", "itemEncoding", ChangeTypeModified))
	assert.True(t, config.IsBreaking("mediaType", "itemEncoding", ChangeTypeRemoved))
}

func TestIsBreaking_OpenAPI32_SecuritySchemeOAuth2MetadataUrl(t *testing.T) {
	config := GenerateDefaultBreakingRules()

	assert.False(t, config.IsBreaking("securityScheme", "oauth2MetadataUrl", ChangeTypeAdded))
	assert.False(t, config.IsBreaking("securityScheme", "oauth2MetadataUrl", ChangeTypeModified))
	assert.False(t, config.IsBreaking("securityScheme", "oauth2MetadataUrl", ChangeTypeRemoved))
}

func TestMerge_OpenAPI32_PathItemQueryOverride(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	assert.False(t, *config.PathItem.Query.Added)
	assert.True(t, *config.PathItem.Query.Removed)

	override := &BreakingRulesConfig{
		PathItem: &PathItemRules{
			Query: &BreakingChangeRule{
				Added:   boolPtr(true),
				Removed: boolPtr(false),
			},
		},
	}
	config.Merge(override)

	assert.True(t, *config.PathItem.Query.Added)
	assert.False(t, *config.PathItem.Query.Removed)
}

func TestMerge_OpenAPI32_MediaTypeItemSchemaOverride(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	assert.True(t, *config.MediaType.ItemSchema.Added)
	assert.True(t, *config.MediaType.ItemSchema.Removed)

	override := &BreakingRulesConfig{
		MediaType: &MediaTypeRules{
			ItemSchema: &BreakingChangeRule{
				Added:   boolPtr(false),
				Removed: boolPtr(false),
			},
		},
	}
	config.Merge(override)

	assert.False(t, *config.MediaType.ItemSchema.Added)
	assert.False(t, *config.MediaType.ItemSchema.Removed)
}

func TestMerge_OpenAPI32_SecuritySchemeOAuth2MetadataUrlOverride(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	assert.False(t, *config.SecurityScheme.OAuth2MetadataUrl.Modified)

	override := &BreakingRulesConfig{
		SecurityScheme: &SecuritySchemeRules{
			OAuth2MetadataUrl: &BreakingChangeRule{
				Modified: boolPtr(true),
			},
		},
	}
	config.Merge(override)

	assert.True(t, *config.SecurityScheme.OAuth2MetadataUrl.Modified)
}

func TestComparePathItems_QueryBreakingChange(t *testing.T) {
	low.ClearHashCache()

	left := `get:
  summary: Get resource
query:
  summary: Query resources
  operationId: queryResources`

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
	assert.Equal(t, 1, changes.TotalBreakingChanges())

	foundQueryRemoved := false
	for _, change := range changes.GetAllChanges() {
		if change.Property == "query" && change.ChangeType == PropertyRemoved {
			foundQueryRemoved = true
			assert.True(t, change.Breaking)
		}
	}
	assert.True(t, foundQueryRemoved, "Query operation removal should be detected as breaking")
}

func TestComparePathItems_AdditionalOperationsBreakingChange(t *testing.T) {
	low.ClearHashCache()

	left := `additionalOperations:
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
	assert.Equal(t, 1, changes.TotalBreakingChanges())

	foundAdditionalOpsRemoved := false
	for _, change := range changes.GetAllChanges() {
		if change.Property == v3.AdditionalOperationsLabel && change.ChangeType == PropertyRemoved {
			foundAdditionalOpsRemoved = true
			assert.True(t, change.Breaking)
		}
	}
	assert.True(t, foundAdditionalOpsRemoved, "AdditionalOperations removal should be detected as breaking")
}

func TestCompareMediaTypes_ItemSchemaBreakingChange(t *testing.T) {
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
	assert.Equal(t, 1, changes.TotalBreakingChanges())

	allChanges := changes.GetAllChanges()
	assert.Len(t, allChanges, 1)
	assert.Equal(t, ObjectRemoved, allChanges[0].ChangeType)
	assert.Equal(t, v3.ItemSchemaLabel, allChanges[0].Property)
	assert.True(t, allChanges[0].Breaking)
}

func TestCompareMediaTypes_ItemEncodingBreakingChange(t *testing.T) {
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

	foundItemEncodingRemoved := false
	for _, change := range changes.GetAllChanges() {
		if change.Property == v3.ItemEncodingLabel && change.ChangeType == ObjectRemoved {
			foundItemEncodingRemoved = true
		}
	}
	assert.True(t, foundItemEncodingRemoved, "ItemEncoding removal should be detected")
}

func TestCompareSecuritySchemes_OAuth2MetadataUrlChange(t *testing.T) {
	low.ClearHashCache()

	left := `type: oauth2
oauth2MetadataUrl: https://auth.example.com/.well-known/oauth_authorization_server`

	right := `type: oauth2
oauth2MetadataUrl: https://auth.different.com/.well-known/oauth_authorization_server`

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
	assert.Equal(t, 0, changes.TotalBreakingChanges())
}

func TestCompareSecuritySchemes_OAuth2MetadataUrlAdded(t *testing.T) {
	low.ClearHashCache()

	left := `type: oauth2`

	right := `type: oauth2
oauth2MetadataUrl: https://auth.example.com/.well-known/oauth_authorization_server`

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
	assert.Equal(t, 0, changes.TotalBreakingChanges())
}

func TestMerge_AllOpenAPI32Properties(t *testing.T) {
	ResetDefaultBreakingRules()
	defer ResetDefaultBreakingRules()

	config := GenerateDefaultBreakingRules()

	override := &BreakingRulesConfig{
		PathItem: &PathItemRules{
			Query:                rule(true, true, false),
			AdditionalOperations: rule(true, true, false),
		},
		MediaType: &MediaTypeRules{
			ItemSchema:   rule(false, true, false),
			ItemEncoding: rule(false, true, false),
		},
		SecurityScheme: &SecuritySchemeRules{
			OAuth2MetadataUrl: rule(true, true, true),
		},
	}
	config.Merge(override)

	assert.True(t, *config.PathItem.Query.Added)
	assert.True(t, *config.PathItem.Query.Modified)
	assert.False(t, *config.PathItem.Query.Removed)

	assert.True(t, *config.PathItem.AdditionalOperations.Added)
	assert.True(t, *config.PathItem.AdditionalOperations.Modified)
	assert.False(t, *config.PathItem.AdditionalOperations.Removed)

	assert.False(t, *config.MediaType.ItemSchema.Added)
	assert.True(t, *config.MediaType.ItemSchema.Modified)
	assert.False(t, *config.MediaType.ItemSchema.Removed)

	assert.False(t, *config.MediaType.ItemEncoding.Added)
	assert.True(t, *config.MediaType.ItemEncoding.Modified)
	assert.False(t, *config.MediaType.ItemEncoding.Removed)

	assert.True(t, *config.SecurityScheme.OAuth2MetadataUrl.Added)
	assert.True(t, *config.SecurityScheme.OAuth2MetadataUrl.Modified)
	assert.True(t, *config.SecurityScheme.OAuth2MetadataUrl.Removed)
}
