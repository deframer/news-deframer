package openapiutil

import (
	"strings"
	"testing"
)

// https://gobyexample.com/testing

func TestFailIfOpenApiJsonFileNotExist(t *testing.T) {
	var arr []string
	_, err := OpenAPI("", "", &arr, nil, "./testfiles/invalid-spec.json")
	if err == nil {
		t.Errorf("no error found for invalid file %s", err)
	}
}

func TestFailOnInvalidOpenApiFile(t *testing.T) {
	var arr []string
	_, err := OpenAPI("", "", &arr, nil, "./testfiles/invalid-spec.json")
	if err == nil {
		t.Errorf("no error found for invalid file %s", err)
	}
}

func TestFilterByTag(t *testing.T) {
	var arr []string
	url := "http://www.example.com"

	files := []string{"./testfiles/filter-by-tag-web.json", "./testfiles/filter-by-tag-checker.json"}
	keepTags := []string{"web", "checker"}
	containsList := []string{"web login", "checker start"}
	notContainsList := []string{"checker start", "web login"}

	for i, file := range files {
		keepTag := keepTags[i]
		contains := containsList[i]
		notContains := notContainsList[i]

		json, err := OpenAPI(url, keepTag, &arr, nil, file)
		if err != nil {
			t.Errorf("no error found for invalid file %s", err)
		}

		if json == "" {
			t.Errorf("got no json data")
		}

		if !strings.Contains(json, url) {
			t.Errorf("url not set")
		}

		if !strings.Contains(json, keepTag) {
			t.Errorf("tag not found")
		}

		if !strings.Contains(json, contains) {
			t.Errorf("elements not removed")
		}

		if strings.Contains(json, notContains) {
			t.Errorf("elements removed")
		}
	}
}

func TestRemoveAllSchemasForUnknowTag(t *testing.T) {
	var arr []string
	url := "http://www.example.com"
	json, err := OpenAPI(url, "unkown", &arr, nil, "./testfiles/filter-by-tag-web.json")
	if err != nil {
		t.Errorf("no error found for invalid file %s", err)
	}

	if strings.Contains(json, "SampleRequest") {
		t.Errorf("schema not removed")
	}
}

func TestFilterPath(t *testing.T) {
	arr := []string{"/resource/service"}
	json, err := OpenAPI("http://www.example.com", "service", &arr, nil, "./testfiles/filter-by-tag-web.json")
	if err != nil {
		t.Errorf("no error found for invalid file %s", err)
	}

	if strings.Contains(json, "/resource/service") {
		t.Errorf("path not filtered")
	}
}

func TestKeepSubSchema(t *testing.T) {
	json, err := OpenAPI("http://www.example.com", "api", nil, nil, "./testfiles/schema-children.json")
	if err != nil {
		t.Errorf("no error found for invalid file %s", err)
	}

	var schemas = []string{"ParentContext", "ChildContext"}

	for _, schema := range schemas {
		const expected = 2
		count := strings.Count(json, schema)

		if count != expected {
			t.Errorf("count mismatch for %s: expected %v - got %v", schema, expected, count)
		}
	}
}

func TestKeepResultSchema(t *testing.T) {
	json, err := OpenAPI("http://www.example.com", "rules", nil, nil, "./testfiles/result-schema.json")
	if err != nil {
		t.Errorf("no error found for invalid file %s", err)
	}

	// var schemas = []string{"GetRuleListResult", "DeleteRuleRequestBody"}
	var schemas = []string{"RuleListResult"}

	for _, schema := range schemas {
		const expected = 2
		count := strings.Count(json, schema)

		if count != expected {
			t.Errorf("count mismatch for %s: expected %v - got %v", schema, expected, count)
		}
	}
}

func TestGetOnlyPaths(t *testing.T) {
	onlyPaths := []string{
		"/resource/delete",
		"/resource/list",
	}

	json, err := OpenAPI("http://www.example.com", "web", nil, &onlyPaths, "./testfiles/path-filter.json")
	if err != nil {
		t.Errorf("no error found for invalid file %s", err)
	}

	if strings.Contains(json, "/resource/show/{id}") {
		t.Errorf("path not filtered")
	}

	if strings.Contains(json, "TokenResult") {
		t.Errorf("content not filtered")
	}

	for _, path := range onlyPaths {
		if !strings.Contains(json, path) {
			t.Errorf("path not kept")
		}
	}
}

func TestGetOnlyPathsWithNestedArrayChildren(t *testing.T) {
	onlyPaths := []string{
		"/resource/items/{group}/{item}/{detail}",
	}

	json, err := OpenAPI("http://www.example.com", "web", nil, &onlyPaths, "./testfiles/nested-array-children.json")
	if err != nil {
		t.Errorf("no error found for invalid file %s", err)
	}

	if strings.Contains(json, "/resource/show/{id}") {
		t.Errorf("path not filtered")
	}

	if !strings.Contains(json, "\"GroupItemEntry\"") {
		t.Errorf("content not filtered")
	}

	if strings.Contains(json, "OpenAPIOnlyForRequestBody") {
		t.Errorf("content filtered")
	}

	for _, path := range onlyPaths {
		if !strings.Contains(json, path) {
			t.Errorf("path not kept")
		}
	}
}
