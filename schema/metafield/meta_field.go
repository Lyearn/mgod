package metafield

import (
	"reflect"

	"github.com/Lyearn/mgod/schema/schemaopt"
	"github.com/Lyearn/mgod/schema/transformer"
	"go.mongodb.org/mongo-driver/bson"
)

type MetaField interface {
	// GetKey returns the unique key of the meta field.
	GetKey() MetaFieldKey

	// GetReflectKind returns the reflect kind of the meta field.
	GetReflectKind() reflect.Kind

	// GetApplicableTransformers returns the list of transformers applicable for the meta field.
	// Meta fields are added to the bson doc before calling the BuildBSONDoc method which transforms the doc.
	GetApplicableTransformers() []transformer.Transformer

	// IsApplicable returns true if the meta field is applicable for the given schema options.
	// Meta field is processed against the doc only if it is applicable.
	IsApplicable(schemaOptions schemaopt.SchemaOptions) bool

	// CheckIfValidValue validates the type of the provided value against the expected type.
	CheckIfValidValue(val interface{}) bool

	// FieldAlreadyPresent modifies the doc at the provided index if the field is already present in the doc
	// and is of the expected type.
	FieldAlreadyPresent(doc *bson.D, index int)

	// FieldPresentWithIncorrectVal modifies the doc at the provided index if the field is already present in the doc
	// but is not of the expected type.
	FieldPresentWithIncorrectVal(doc *bson.D, index int) error

	// FieldNotPresent appends the missing field in the doc.
	FieldNotPresent(doc *bson.D)
}

var AvailableMetaFields = []MetaField{
	createdAtMetaFieldInstance,
	updatedAtMetaFieldInstance,
	docVersionMetaFieldInstance,
}

func AddMetaFields(bsonDoc *bson.D, schemaOptions schemaopt.SchemaOptions) error {
	for _, metaField := range AvailableMetaFields {
		if !metaField.IsApplicable(schemaOptions) {
			continue
		}

		if err := validatedAndAddFieldValue(bsonDoc, metaField); err != nil {
			return err
		}
	}

	return nil
}
