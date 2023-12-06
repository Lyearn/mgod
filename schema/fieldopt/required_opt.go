package fieldopt

import (
	"reflect"
	"strings"

	"github.com/samber/lo"
)

type requiredOption struct{}

func newRequiredOption() FieldOption {
	return &requiredOption{}
}

// RequiredOption defines if a field is required or not.
// The option can be invalidated using `omitempty` property of `bson` tag.
// Defaults to true for all fields.
var RequiredOption = newRequiredOption()

func (o requiredOption) GetOptName() string {
	return "Required"
}

func (o requiredOption) GetBSONTagName() string {
	return string(FieldOptionTagRequired)
}

func (o requiredOption) IsApplicable(field reflect.StructField) bool {
	return true
}

func (o requiredOption) GetDefaultValue(field reflect.StructField) interface{} {
	return true
}

func (o requiredOption) GetValue(field reflect.StructField) (interface{}, error) {
	tagVal := field.Tag.Get(o.GetBSONTagName())
	splittedTagValues := strings.Split(tagVal, ",")

	_, found := lo.Find(splittedTagValues[1:], func(tagVal string) bool {
		return tagVal == "omitempty"
	})

	// if omitempty is found => required is false and vice versa
	return !found, nil
}
