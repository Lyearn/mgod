package schemaopt

import "reflect"

type XIDOption struct{}

func newXIDOption() SchemaOption {
	return &XIDOption{}
}

func (o XIDOption) GetOptName() string {
	return "XID"
}

func (o XIDOption) GetBSONTagName() string {
	return "mgoID"
}

func (o XIDOption) IsApplicable(field reflect.StructField) bool {
	return field.Type.Kind() == reflect.Struct
}

func (o XIDOption) GetDefaultValue(field reflect.StructField) interface{} {
	// if the field is not applicable, then the default value should be false
	defaultValue := true

	if !o.IsApplicable(field) {
		defaultValue = false
	}

	return defaultValue
}

func (o XIDOption) GetValue(field reflect.StructField) (interface{}, error) {
	tagVal := field.Tag.Get(o.GetBSONTagName())
	isXIDRequired := true

	if tagVal == "false" {
		isXIDRequired = false
	}

	return isXIDRequired, nil
}

var xidOptionInstance = newXIDOption()
