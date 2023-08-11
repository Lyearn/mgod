package metafield

import (
	"reflect"
	"time"

	"github.com/Lyearn/backend-universe/packages/common/dateformatter"
	"github.com/Lyearn/backend-universe/packages/store/acl/model"
	"github.com/Lyearn/backend-universe/packages/store/mongomodel/transformer"
	"go.mongodb.org/mongo-driver/bson"
)

type CreatedAtMetaField struct{}

func newCreatedAtMetaField() MetaField {
	return &CreatedAtMetaField{}
}

var createdAtMetaFieldInstance = newCreatedAtMetaField()

func (m CreatedAtMetaField) GetKey() MetaFieldKey {
	return MetaFieldKeyCreatedAt
}

func (m CreatedAtMetaField) GetReflectKind() reflect.Kind {
	return reflect.String
}

func (m CreatedAtMetaField) GetApplicableTransformers() []transformer.Transformer {
	return []transformer.Transformer{transformer.DateTransformerInstance}
}

func (m CreatedAtMetaField) IsApplicable(schemaOptions model.SchemaOptions) bool {
	return schemaOptions.Timestamps
}

func (m CreatedAtMetaField) CheckIfValidValue(val interface{}) bool {
	if val, ok := val.(string); ok && val != "" {
		return true
	}

	return false
}

func (m CreatedAtMetaField) FieldAlreadyPresent(doc *bson.D, index int) {
	// do nothing.
}

func (m CreatedAtMetaField) FieldPresentWithIncorrectVal(doc *bson.D, index int) error {
	isoString, err := dateformatter.New(time.Now().UTC()).GetISOString()
	if err != nil {
		return err
	}

	(*doc)[index].Value = isoString

	return nil
}

func (m CreatedAtMetaField) FieldNotPresent(doc *bson.D) {
	isoString, _ := dateformatter.New(time.Now().UTC()).GetISOString()
	*doc = append(*doc, bson.E{
		Key:   string(m.GetKey()),
		Value: isoString,
	})
}
