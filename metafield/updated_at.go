package metafield

import (
	"time"

	"github.com/Lyearn/backend-universe/packages/common/dateformatter"
	"github.com/Lyearn/backend-universe/packages/store/acl/model"
	"go.mongodb.org/mongo-driver/bson"
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
	if val, ok := val.(string); ok && val != "" {
		return true
	}

	return false
}

func (m UpdatedAtMetaField) FieldAlreadyPresent(doc *bson.D, index int) {
	// field is already present. hence, updating the value.
	isoString, _ := dateformatter.New(time.Now().UTC()).GetISOString()
	(*doc)[index].Value = isoString
}

func (m UpdatedAtMetaField) FieldPresentWithIncorrectVal(doc *bson.D, index int) {
	isoString, _ := dateformatter.New(time.Now().UTC()).GetISOString()
	(*doc)[index].Value = isoString
}

func (m UpdatedAtMetaField) FieldNotPresent(doc *bson.D) {
	isoString, _ := dateformatter.New(time.Now().UTC()).GetISOString()
	*doc = append(*doc, bson.E{
		Key:   string(m.GetMetaFieldKey()),
		Value: isoString,
	})
}
