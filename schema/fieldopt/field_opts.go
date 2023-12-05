// Package fieldopt provides custom options for schema fields.
package fieldopt

import (
	"reflect"

	"github.com/samber/lo"
)

// FieldOption is the interface that needs to be implemented by all [SchemaFieldOptions].
type FieldOption interface {
	// GetOptName returns the name of the schema option. This name is used to identify the unique option.
	// NOTE: Make sure to return the same name as the name of the field in [SchemaFieldOptions] struct.
	GetOptName() string
	// GetBSONTagName returns the bson tag name for the option. This is used to identify the option and its flags in the bson tag.
	GetBSONTagName() string
	// IsApplicable returns true if the option is applicable for the field.
	IsApplicable(field reflect.StructField) bool
	// GetDefaultValue returns the default value for the field (if available).
	GetDefaultValue(field reflect.StructField) interface{}
	// GetValue returns the provided value for the field.
	GetValue(field reflect.StructField) (interface{}, error)
}

// SchemaFieldOptions are custom schema options available for struct fields.
// These options either modifies the schema based on the field or adds validations to the field.
type SchemaFieldOptions struct {
	// Required suggests whether the field is required or not. [FIELD_LEVEL]
	// Defaults to true. Can be identified using omitempty bson flag.
	Required bool
	// XID suggests whether "_id" field needs to be added in the bson doc for the following object type field. [STRUCT_LEVEL]
	// Defaults to true. This option is applicable for fields holding structs only.
	XID bool
	// Default is the default value for the field. [FIELD_LEVEL]
	// Defaults to nil. Will be populated using reflect and will be of the same type as Type in SchemaFieldProps.
	Default interface{}
	// not implemented yet
	Select bool
}

var availableSchemaOptions = []FieldOption{
	RequiredOption,
	XIDOption,
	DefaultValueOption,
}

var optNameToSchemaOptionMap = lo.KeyBy(availableSchemaOptions, func(opt FieldOption) string {
	return opt.GetOptName()
})

// GetSchemaOptionsForField returns all the applicable schema field options for the provided field.
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
