package schemaopt

import (
	"reflect"
	"strings"

	"github.com/samber/lo"
)

type RequiredOption struct{}

func newRequiredOption() SchemaOption {
	return &RequiredOption{}
}

func (o RequiredOption) GetOptName() string {
	return "Required"
}

func (o RequiredOption) GetBSONTagName() string {
	return "bson"
}

func (o RequiredOption) IsApplicable(field reflect.StructField) bool {
	return true
}

func (o RequiredOption) GetDefaultValue(field reflect.StructField) interface{} {
	return true
}

func (o RequiredOption) GetValue(field reflect.StructField) (interface{}, error) {
	tagVal := field.Tag.Get(o.GetBSONTagName())
	splittedTagValues := strings.Split(tagVal, ",")

	_, found := lo.Find(splittedTagValues[1:], func(tagVal string) bool {
		return tagVal == "omitempty"
	})

	// if omitempty is found => required is false and vice versa
	return !found, nil
}

var requiredOptionInstance = newRequiredOption()
