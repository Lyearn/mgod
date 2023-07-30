package metafield

import (
	"github.com/Lyearn/backend-universe/packages/store/acl/model"
	"go.mongodb.org/mongo-driver/bson"
)

type MetaField interface {
	GetMetaFieldKey() MetaFieldKey

	// IsApplicable returns true if the meta field is applicable for the given schema options.
	// Meta field is processed against the doc only if it is applicable.
	IsApplicable(schemaOptions model.SchemaOptions) bool

	// CheckIfValidValue validates the type of the provided value against the expected type.
	CheckIfValidValue(val interface{}) bool

	// FieldAlreadyPresent modifies the doc at the provided index if the field is already present in the doc
	// and is of the expected type.
	FieldAlreadyPresent(doc *bson.D, index int)

	// FieldPresentWithIncorrectVal modifies the doc at the provided index if the field is already present in the doc
	// but is not of the expected type.
	FieldPresentWithIncorrectVal(doc *bson.D, index int)

	// FieldNotPresent appends the missing field in the doc.
	FieldNotPresent(doc *bson.D)
}

var availableMetaFields = []MetaField{
	createdAtMetaFieldInstance,
	updatedAtMetaFieldInstance,
	docVersionMetaFieldInstance,
}

func AddMetaFields(bsonDoc *bson.D, schemaOptions model.SchemaOptions) {
	for _, metaField := range availableMetaFields {
		if !metaField.IsApplicable(schemaOptions) {
			continue
		}

		ValidatedAndAddFieldValue(bsonDoc, metaField)
	}
}
