package metafield

import (
	"github.com/Lyearn/backend-universe/packages/store/acl/model"
	"go.mongodb.org/mongo-driver/bson"
)

type DocVersionMetaField struct{}

func newDocVersionMetaField() MetaField {
	return &DocVersionMetaField{}
}

var docVersionMetaFieldInstance = newDocVersionMetaField()

func (m DocVersionMetaField) GetMetaFieldKey() MetaFieldKey {
	return MetaFieldKeyDocVersion
}

func (m DocVersionMetaField) IsApplicable(schemaOptions model.SchemaOptions) bool {
	// doc versioning is enabled by default.
	return true
}

func (m DocVersionMetaField) CheckIfValidValue(val interface{}) bool {
	_, ok := val.(int)
	return ok
}

func (m DocVersionMetaField) FieldAlreadyPresent(doc *bson.D, index int) {
	// field is already present. hence, incrementing the value.
	(*doc)[index].Value = (*doc)[index].Value.(int) + 1
}

func (m DocVersionMetaField) FieldPresentWithIncorrectVal(doc *bson.D, index int) {
	(*doc)[index].Value = 0
}

func (m DocVersionMetaField) FieldNotPresent(doc *bson.D) {
	*doc = append(*doc, bson.E{
		Key:   string(m.GetMetaFieldKey()),
		Value: 0,
	})
}
