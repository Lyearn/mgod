package fieldopt

import (
	"reflect"

	"github.com/samber/lo"
)

type FieldOption interface {
	// GetOptName returns the name of the schema option. This name is used to identify the unique option.
	// NOTE: Make sure to return the same name as the name of the struct field.
	GetOptName() string
	// bson tag name for the option. This is used to identify the option and its flags in the bson tag.
	GetBSONTagName() string
	IsApplicable(field reflect.StructField) bool
	GetDefaultValue(field reflect.StructField) interface{}
	GetValue(field reflect.StructField) (interface{}, error)
}

type SchemaFieldOptions struct {
	Required bool        // defaults to true. can be identified using omitempty bson flag. [FIELD_LEVEL]
	XID      bool        // defaults to true wherever applicable. [STRUCT_LEVEL]
	Default  interface{} // will be populated using reflect. will be of same type as Type in SchemaFieldProps. [FIELD_LEVEL]
	Select   bool
}

var availableSchemaOptions = []FieldOption{
	requiredOptionInstance,
	xidOptionInstance,
	defaultValueOptionInstance,
}

var optNameToSchemaOptionMap = lo.KeyBy(availableSchemaOptions, func(opt FieldOption) string {
	return opt.GetOptName()
})

func GetSchemaOptionsForField(field reflect.StructField) (SchemaFieldOptions, error) {
	options := SchemaFieldOptions{}

	optsVal := reflect.ValueOf(&options).Elem()

	for i := 0; i < optsVal.NumField(); i++ {
		structField := optsVal.Type().Field(i)
		fieldName := structField.Name

		schemaOption, ok := optNameToSchemaOptionMap[fieldName]
		if !ok {
			continue
		}

		var valueToSet reflect.Value

		if schemaOption.IsApplicable(field) {
			fieldVal, err := schemaOption.GetValue(field)
			if err != nil {
				return options, err
			}
			valueToSet = reflect.ValueOf(fieldVal)
		} else {
			valueToSet = reflect.ValueOf(schemaOption.GetDefaultValue(field))
		}

		if valueToSet.IsValid() {
			optsVal.Field(i).Set(valueToSet)
		}
	}

	return options, nil
}
