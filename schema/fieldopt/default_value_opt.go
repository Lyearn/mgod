package fieldopt

import (
	"fmt"
	"reflect"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

type defaultValueOption struct{}

func newDefaultValueOption() FieldOption {
	return &defaultValueOption{}
}

// DefaultValueOption provides the default value for a field.
// This value of this option is used when the field is not present in the input document.
// This option is applicable only for fields that are not of type struct.
// Defaults to nil for all fields.
var DefaultValueOption = newDefaultValueOption()

func (o defaultValueOption) GetOptName() string {
	return "Default"
}

func (o defaultValueOption) GetBSONTagName() string {
	return string(FieldOptionTagDefault)
}

func (o defaultValueOption) IsApplicable(field reflect.StructField) bool {
	// not available on struct fields
	if field.Type.Kind() == reflect.Struct {
		return false
	}

	// check if the field has a default value
	tagVal := field.Tag.Get(o.GetBSONTagName())

	return tagVal != ""
}

func (o defaultValueOption) GetDefaultValue(field reflect.StructField) interface{} {
	return nil
}

func (o defaultValueOption) GetValue(field reflect.StructField) (interface{}, error) {
	tagVal := field.Tag.Get(o.GetBSONTagName())

	fieldType := field.Type.Kind()
	if fieldType == reflect.Ptr {
		elem := field.Type.Elem()
		fieldType = elem.Kind()
	}

	switch fieldType {
	case reflect.String:
		return tagVal, nil

	case reflect.Int:
		return strconv.Atoi(tagVal)

	case reflect.Float32:
		return strconv.ParseFloat(tagVal, 32)

	case reflect.Float64:
		return strconv.ParseFloat(tagVal, 64)

	case reflect.Bool:
		return strconv.ParseBool(tagVal)

	case reflect.Slice, reflect.Array:
		return bson.A{}, nil

	default:
		return nil, fmt.Errorf("unsupported type %v", fieldType)
	}
}
