package metafield

import (
	"time"

	"github.com/Lyearn/backend-universe/packages/store/acl/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreatedAtMetaField struct{}

func newCreatedAtMetaField() MetaField {
	return &CreatedAtMetaField{}
}

var createdAtMetaFieldInstance = newCreatedAtMetaField()

func (m CreatedAtMetaField) GetMetaFieldKey() MetaFieldKey {
	return MetaFieldKeyCreatedAt
}

func (m CreatedAtMetaField) IsApplicable(schemaOptions model.SchemaOptions) bool {
	return schemaOptions.Timestamps
}

func (m CreatedAtMetaField) CheckIfValidValue(val interface{}) bool {
	_, ok := val.(primitive.DateTime)
	return ok
}

func (m CreatedAtMetaField) FieldAlreadyPresent(doc *bson.D, index int) {
	// do nothing.
}

func (m CreatedAtMetaField) FieldPresentWithIncorrectVal(doc *bson.D, index int) {
	(*doc)[index].Value = primitive.NewDateTimeFromTime(time.Now())
}

func (m CreatedAtMetaField) FieldNotPresent(doc *bson.D) {
	*doc = append(*doc, bson.E{
		Key:   string(m.GetMetaFieldKey()),
		Value: primitive.NewDateTimeFromTime(time.Now()),
	})
}
