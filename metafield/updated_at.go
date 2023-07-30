package metafield

import (
	"time"

	"github.com/Lyearn/backend-universe/packages/store/acl/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdatedAtMetaField struct{}

func newUpdatedAtMetaField() MetaField {
	return &UpdatedAtMetaField{}
}

var updatedAtMetaFieldInstance = newUpdatedAtMetaField()

func (m UpdatedAtMetaField) GetMetaFieldKey() MetaFieldKey {
	return MetaFieldKeyUpdatedAt
}

func (m UpdatedAtMetaField) IsApplicable(schemaOptions model.SchemaOptions) bool {
	return schemaOptions.Timestamps
}

func (m UpdatedAtMetaField) CheckIfValidValue(val interface{}) bool {
	_, ok := val.(primitive.DateTime)
	return ok
}

func (m UpdatedAtMetaField) FieldAlreadyPresent(doc *bson.D, index int) {
	// field is already present. hence, updating the value.
	(*doc)[index].Value = primitive.NewDateTimeFromTime(time.Now())
}

func (m UpdatedAtMetaField) FieldPresentWithIncorrectVal(doc *bson.D, index int) {
	(*doc)[index].Value = primitive.NewDateTimeFromTime(time.Now())
}

func (m UpdatedAtMetaField) FieldNotPresent(doc *bson.D) {
	*doc = append(*doc, bson.E{
		Key:   string(m.GetMetaFieldKey()),
		Value: primitive.NewDateTimeFromTime(time.Now()),
	})
}
