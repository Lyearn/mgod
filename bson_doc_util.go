package mgod

import (
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func getFieldValueFromBSONRootDoc(doc *bson.D, field string) interface{} {
	if doc == nil {
		return nil
	}

	for _, elem := range *doc {
		if elem.Key == field {
			return elem.Value
		}
	}

	return nil
}

func getCurrentLevelBSONFields(s reflect.Value) []string {
	currentLevelBSONFields := make([]string, 0)

	for i := 0; i < s.NumField(); i++ {
		structField := s.Type().Field(i)
		fieldName := getBSONFieldName(structField)

		currentLevelBSONFields = append(currentLevelBSONFields, fieldName)
	}

	return currentLevelBSONFields
}

func isBSONInlineField(field reflect.StructField) bool {
	bsonTagVal := field.Tag.Get("bson")
	if bsonTagVal == "" || bsonTagVal == "-" {
		return false
	}

	tagValues := strings.Split(bsonTagVal, ",")
	flags := tagValues[1:]

	for _, flag := range flags {
		if flag == "inline" {
			return true
		}
	}

	return false
}

func getBSONFieldName(field reflect.StructField) string {
	bsonTag := field.Tag.Get("bson")
	if bsonTag == "" {
		return getDefaultBSONFieldName(field)
	}

	parts := strings.Split(bsonTag, ",")
	bsonFieldName := parts[0]

	// the case when bson field needs to be ignored
	if bsonFieldName == "-" {
		return ""
	}

	// the case when only bson flags are specified and name is expected to be default bson field name
	if bsonFieldName == "" {
		return getDefaultBSONFieldName(field)
	}

	return bsonFieldName
}

func getDefaultBSONFieldName(field reflect.StructField) string {
	fieldName := field.Name
	if len(fieldName) == 0 {
		return fieldName
	}

	return strings.ToLower(fieldName)
}
