package metafield

import (
	"reflect"
	"time"

	"github.com/Lyearn/mgod/dateformatter"
	"github.com/Lyearn/mgod/schema/schemaopt"
	"github.com/Lyearn/mgod/schema/transformer"
	"go.mongodb.org/mongo-driver/bson"
)

type createdAtMetaField struct{}

func newCreatedAtMetaField() MetaField {
	return &createdAtMetaField{}
}

var createdAtMetaFieldInstance = newCreatedAtMetaField()

func (m createdAtMetaField) GetKey() MetaFieldKey {
	return MetaFieldKeyCreatedAt
}

func (m createdAtMetaField) GetReflectKind() reflect.Kind {
	return reflect.String
}

func (m createdAtMetaField) GetApplicableTransformers() []transformer.Transformer {
	return []transformer.Transformer{transformer.DateTransformerInstance}
}

func (m createdAtMetaField) IsApplicable(schemaOptions schemaopt.SchemaOptions) bool {
	return schemaOptions.Timestamps
}

func (m createdAtMetaField) CheckIfValidValue(val interface{}) bool {
	if val, ok := val.(string); ok && val != "" {
		return true
	}

	return false
}

func (m createdAtMetaField) FieldAlreadyPresent(doc *bson.D, index int) {
	// do nothing.
}

func (m createdAtMetaField) FieldPresentWithIncorrectVal(doc *bson.D, index int) error {
	isoString, err := dateformatter.New(time.Now().UTC()).GetISOString()
	if err != nil {
		return err
	}

	(*doc)[index].Value = isoString

	return nil
}

func (m createdAtMetaField) FieldNotPresent(doc *bson.D) {
	isoString, _ := dateformatter.New(time.Now().UTC()).GetISOString()
	*doc = append(*doc, bson.E{
		Key:   string(m.GetKey()),
		Value: isoString,
	})
}
