package metafield

import (
	"reflect"

	"github.com/Lyearn/mgod/schema/schemaopt"
	"github.com/Lyearn/mgod/schema/transformer"
	"go.mongodb.org/mongo-driver/bson"
)

type docVersionMetaField struct{}

func newDocVersionMetaField() MetaField {
	return &docVersionMetaField{}
}

var docVersionMetaFieldInstance = newDocVersionMetaField()

func (m docVersionMetaField) GetKey() MetaFieldKey {
	return MetaFieldKeyDocVersion
}

func (m docVersionMetaField) GetReflectKind() reflect.Kind {
	return reflect.Int
}

func (m docVersionMetaField) IsApplicable(schemaOptions schemaopt.SchemaOptions) bool {
	if schemaOptions.VersionKey == nil {
		// doc versioning is enabled by default.
		return true
	}

	return *schemaOptions.VersionKey
}

func (m docVersionMetaField) GetApplicableTransformers() []transformer.Transformer {
	return []transformer.Transformer{}
}

func (m docVersionMetaField) CheckIfValidValue(val interface{}) bool {
	_, ok := val.(int)
	return ok
}

func (m docVersionMetaField) FieldAlreadyPresent(doc *bson.D, index int) {
	// field is already present. hence, incrementing the value.
	(*doc)[index].Value = (*doc)[index].Value.(int) + 1
}

func (m docVersionMetaField) FieldPresentWithIncorrectVal(doc *bson.D, index int) error {
	(*doc)[index].Value = 0

	return nil
}

func (m docVersionMetaField) FieldNotPresent(doc *bson.D) {
	*doc = append(*doc, bson.E{
		Key:   string(m.GetKey()),
		Value: 0,
	})
}
