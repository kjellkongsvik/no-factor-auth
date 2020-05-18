package controllers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"sort"
	"testing"
)

func TestShouldReturnParseError(t *testing.T) {
	jsonString := []byte ("malformed json")

	_, err := ParseExtraClaims(jsonString)
	assert.Error(t, err, "Malformed json should return error")

}

func TestParseJson(t *testing.T) {
	oidValue := "b213-61024b63a7ea"
	arrValue := "some-value"

	jsonString := []byte (`{"oid":"` + oidValue + `", "arr_elem_key":["` + arrValue + `","foo"]}`)

	result, err := ParseExtraClaims(jsonString)

	if assert.NoError(t, err) {
		assert.Equal(t, oidValue, result["oid"], "Simple json elements should be parsed")

		arrElements := extractArrayElement(result)
		assert.True(t, contains(arrElements, arrValue), "Json array should be parsed")
	}

}

func contains(s []string, searchterm string) bool {
	sort.Strings(s)
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}

func extractArrayElement(result map[string]interface{}) []string {
	var items []string

	object := reflect.ValueOf(result["arr_elem_key"])
	for i := 0; i < object.Len(); i++ {
		s := fmt.Sprint(object.Index(i).Interface())
		items = append(items, s)
	}

	return items
}
