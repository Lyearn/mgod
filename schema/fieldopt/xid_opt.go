package fieldopt

import "reflect"

type xidOption struct{}

func newXIDOption() FieldOption {
	return &xidOption{}
}

func (o xidOption) GetOptName() string {
	return "XID"
}

func (o xidOption) GetBSONTagName() string {
	return "mgoID"
}

func (o xidOption) IsApplicable(field reflect.StructField) bool {
	return field.Type.Kind() == reflect.Struct
}

func (o xidOption) GetDefaultValue(field reflect.StructField) interface{} {
	// if the field is not applicable, then the default value should be false
	defaultValue := true

	if !o.IsApplicable(field) {
		defaultValue = false
	}

	return defaultValue
}

func (o xidOption) GetValue(field reflect.StructField) (interface{}, error) {
	tagVal := field.Tag.Get(o.GetBSONTagName())
	isXIDRequired := true

	if tagVal == "false" {
		isXIDRequired = false
	}

	return isXIDRequired, nil
}

var xidOptionInstance = newXIDOption()
