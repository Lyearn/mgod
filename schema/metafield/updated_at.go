package metafield

import (
	"reflect"
	"time"

	"github.com/Lyearn/mgod/dateformatter"
	"github.com/Lyearn/mgod/schema/schemaopt"
	"github.com/Lyearn/mgod/schema/transformer"
	"go.mongodb.org/mongo-driver/bson"
)

type updatedAtMetaField struct{}

func newUpdatedAtMetaField() MetaField {
	return &updatedAtMetaField{}
}

var updatedAtMetaFieldInstance = newUpdatedAtMetaField()

func (m updatedAtMetaField) GetKey() MetaFieldKey {
	return MetaFieldKeyUpdatedAt
}

func (m updatedAtMetaField) GetReflectKind() reflect.Kind {
	return reflect.String
}

func (m updatedAtMetaField) GetApplicableTransformers() []transformer.Transformer {
	return []transformer.Transformer{transformer.DateTransformerInstance}
}

func (m updatedAtMetaField) IsApplicable(schemaOptions schemaopt.SchemaOptions) bool {
	return schemaOptions.Timestamps
}

func (m updatedAtMetaField) CheckIfValidValue(val interface{}) bool {
	if val, ok := val.(string); ok && val != "" {
		return true
	}

	return false
}

func (m updatedAtMetaField) FieldAlreadyPresent(doc *bson.D, index int) {
	// field is already present. hence, updating the value.
	isoString, _ := dateformatter.New(time.Now().UTC()).GetISOString()
	(*doc)[index].Value = isoString
}

func (m updatedAtMetaField) FieldPresentWithIncorrectVal(doc *bson.D, index int) error {
	isoString, err := dateformatter.New(time.Now().UTC()).GetISOString()
	if err != nil {
		return err
	}

	(*doc)[index].Value = isoString

	return nil
}

func (m updatedAtMetaField) FieldNotPresent(doc *bson.D) {
	isoString, _ := dateformatter.New(time.Now().UTC()).GetISOString()
	*doc = append(*doc, bson.E{
		Key:   string(m.GetKey()),
		Value: isoString,
	})
}
