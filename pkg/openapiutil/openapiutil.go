// Package openapiutil Services to allow modification of openapi3 file goa generated
package openapiutil

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
)

func scanSchemaType(schema *base.SchemaProxy, refs *[]string) {
	if slices.Contains(schema.Schema().Type, "array") {
		if schema.Schema().Items.IsA() {
			proxy := schema.Schema().Items.A
			if proxy.IsReference() {
				refStr := proxy.GetReference()
				refArr := strings.Split(refStr, "/")
				ref := refArr[len(refArr)-1]
				*refs = append(*refs, ref)
			}
		}
	}

	if !schema.IsReference() {
		return
	}

	// recursive call
	for pair := orderedmap.First(schema.Schema().Properties); pair != nil; pair = pair.Next() {
		data := pair.Value()
		scanSchemaType(data, refs)
	}

	refStr := schema.GetReference()
	refArr := strings.Split(refStr, "/")
	ref := refArr[len(refArr)-1]
	*refs = append(*refs, ref)
}

func scanMediaType(mediaType *v3.MediaType, refs *[]string) {
	if slices.Contains(mediaType.Schema.Schema().Type, "array") {
		if mediaType.Schema.Schema().Items.IsA() {
			proxy := mediaType.Schema.Schema().Items.A
			if proxy.IsReference() {
				refStr := proxy.GetReference()
				refArr := strings.Split(refStr, "/")
				ref := refArr[len(refArr)-1]
				*refs = append(*refs, ref)
			}
		}
	}

	if !mediaType.Schema.IsReference() {
		return
	}

	for pair := orderedmap.First(mediaType.Schema.Schema().Properties); pair != nil; pair = pair.Next() {
		data := pair.Value()
		scanSchemaType(data, refs)
	}

	refStr := mediaType.Schema.GetReference()
	refArr := strings.Split(refStr, "/")
	ref := refArr[len(refArr)-1]
	*refs = append(*refs, ref)
}

func scanSchema(op *v3.Operation, refs *[]string) {
	if op == nil {
		return
	}

	if op.RequestBody != nil && op.RequestBody.Content != nil {
		for pair := orderedmap.First(op.RequestBody.Content); pair != nil; pair = pair.Next() {
			mediaType := pair.Value()
			scanMediaType(mediaType, refs)
		}
	}

	if op.Responses != nil && orderedmap.Len(op.Responses.Codes) > 0 {
		for pairRC := orderedmap.First(op.Responses.Codes); pairRC != nil; pairRC = pairRC.Next() {
			response := pairRC.Value()
			if response.Content != nil {
				for pair := orderedmap.First(response.Content); pair != nil; pair = pair.Next() {
					mediaType := pair.Value()
					scanMediaType(mediaType, refs)
				}
			}
		}
	}
}

func addMissingChildSchemas(docModel *libopenapi.DocumentModel[v3.Document], refsFromPaths []string) []string {
	result := make([]string, 0, 10)

	for len(refsFromPaths) != 0 {

		// get from head
		ref := refsFromPaths[0]
		refsFromPaths = slices.Delete(refsFromPaths, 0, 1)

		if slices.Contains(result, ref) {
			// we already processed this
			continue
		}

		result = append(result, ref)

		// scan for children of the items
		schema, ok := docModel.Model.Components.Schemas.Get(ref)
		if !ok {
			continue
		}

		children := make([]string, 0, 100)
		for pair := orderedmap.First(schema.Schema().Properties); pair != nil; pair = pair.Next() {
			childSchema := pair.Value()
			scanSchemaType(childSchema, &children)
		}

		for _, child := range children {
			// append if we don't have it
			if !slices.Contains(result, child) && !slices.Contains(refsFromPaths, child) {
				refsFromPaths = append(refsFromPaths, child)
			}
		}

	}

	return result
}

func OpenAPI(url string, wantedTag string, removePaths *[]string, onlyPaths *[]string, fileName string) (string, error) {
	// #nosec G304 -- fileName is provided by internal callers and points to generated OpenAPI fixtures.
	bytes, err := os.ReadFile(fileName)

	if err != nil {
		return "", err
	}

	// create a new Document from from the byte slice.
	document, err := libopenapi.NewDocument(bytes)
	if err != nil {
		return "", err
	}

	docModel, err := document.BuildV3Model()

	if err != nil {
		return "", fmt.Errorf("cannot create v3 model from document: %w", err)
	}

	for {
		found := false
		for idx, tag := range docModel.Model.Tags {
			if tag.Name != wantedTag {
				docModel.Model.Tags = slices.Delete(docModel.Model.Tags, idx, idx+1)
				found = true
				break
			}
		}
		if !found {
			break
		}
	}

	for {
		found := false

		for pathPair := orderedmap.First(docModel.Model.Paths.PathItems); pathPair != nil; pathPair = pathPair.Next() {
			key := pathPair.Key()
			path := pathPair.Value()
			operationRemove := false

			if removePaths != nil && len(*removePaths) > 0 && slices.Contains(*removePaths, key) {
				operationRemove = true
			}

			if onlyPaths != nil && len(*onlyPaths) > 0 && !slices.Contains(*onlyPaths, key) {
				operationRemove = true
			}

			for opPair := orderedmap.First(path.GetOperations()); opPair != nil; opPair = opPair.Next() {
				operation := opPair.Value()
				if !slices.Contains(operation.Tags, wantedTag) {
					operationRemove = true
					break
				}
			}

			if operationRemove {
				docModel.Model.Paths.PathItems.Delete(key)
				found = true
				break
			}
		}

		if !found {
			break
		}
	}

	refs := make([]string, 0, 10)

	for pair := orderedmap.First(docModel.Model.Paths.PathItems); pair != nil; pair = pair.Next() {
		path := pair.Value()
		scanSchema(path.Get, &refs)
		scanSchema(path.Put, &refs)
		scanSchema(path.Post, &refs)
		scanSchema(path.Delete, &refs)
		scanSchema(path.Options, &refs)
		scanSchema(path.Head, &refs)
		scanSchema(path.Patch, &refs)
		scanSchema(path.Trace, &refs)
	}

	refs = addMissingChildSchemas(docModel, refs)

	for orderedmap.Len(docModel.Model.Components.Schemas) != 0 {
		found := false
		for pair := orderedmap.First(docModel.Model.Components.Schemas); pair != nil; pair = pair.Next() {
			schema := pair.Key()
			hasUnwantedSchema := false

			if len(refs) == 0 || !slices.Contains(refs, schema) {
				hasUnwantedSchema = true
			}

			if hasUnwantedSchema {
				docModel.Model.Components.Schemas.Delete(schema)
				found = true
				break
			}
		}

		if !found {
			break
		}
	}

	if len(docModel.Model.Servers) == 0 {
		return "", fmt.Errorf("document has no servers")
	}
	if len(docModel.Model.Servers) > 1 {
		docModel.Model.Servers = docModel.Model.Servers[:1]
	}

	docModel.Model.Servers[0].URL = url
	docModel.Model.Servers[0].Description = ""
	docModel.Model.Servers[0].Variables = nil
	docModel.Model.Servers[0].Extensions = nil

	if orderedmap.Len(docModel.Model.Paths.PathItems) == 0 {
		// something crashes, if we removed everything - so return nothing :)
		return "", nil
	}

	data, err := docModel.Model.RenderJSON("")
	if err != nil {
		return "", err
	}
	json := string(data)

	return json, nil
}
