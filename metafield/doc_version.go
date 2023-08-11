package metafield

import (
	"reflect"

	"github.com/Lyearn/backend-universe/packages/store/acl/model"
	"github.com/Lyearn/backend-universe/packages/store/mongomodel/transformer"
	"go.mongodb.org/mongo-driver/bson"
)

type DocVersionMetaField struct{}

func newDocVersionMetaField() MetaField {
	return &DocVersionMetaField{}
}

var docVersionMetaFieldInstance = newDocVersionMetaField()

func (m DocVersionMetaField) GetKey() MetaFieldKey {
	return MetaFieldKeyDocVersion
}

func (m DocVersionMetaField) GetReflectKind() reflect.Kind {
	return reflect.Int
}

func (m DocVersionMetaField) IsApplicable(schemaOptions model.SchemaOptions) bool {
	if schemaOptions.VersionKey == nil {
		// doc versioning is enabled by default.
		return true
	}

	return *schemaOptions.VersionKey
}

func (m DocVersionMetaField) GetApplicableTransformers() []transformer.Transformer {
	return []transformer.Transformer{}
}

func (m DocVersionMetaField) CheckIfValidValue(val interface{}) bool {
	_, ok := val.(int)
	return ok
}

func (m DocVersionMetaField) FieldAlreadyPresent(doc *bson.D, index int) {
	// field is already present. hence, incrementing the value.
	(*doc)[index].Value = (*doc)[index].Value.(int) + 1
}

func (m DocVersionMetaField) FieldPresentWithIncorrectVal(doc *bson.D, index int) error {
	(*doc)[index].Value = 0

	return nil
}

func (m DocVersionMetaField) FieldNotPresent(doc *bson.D) {
	*doc = append(*doc, bson.E{
		Key:   string(m.GetKey()),
		Value: 0,
	})
}
