package transformer

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type dateTransformer struct{}

func newDateTransformer() Transformer {
	return &dateTransformer{}
}

func (t dateTransformer) isTransformationRequired(field reflect.StructField) bool {
	return field.Tag.Get("mgoType") == "date"
}

func (t dateTransformer) TransformForMongoDoc(value interface{}) (interface{}, error) {
	primitiveDates, err := convertStringToDateTime(value.(string))
	if err != nil {
		return nil, err
	}

	return primitiveDates[0], nil
}

func (t dateTransformer) TransformForEntityModelDoc(value interface{}) (interface{}, error) {
	dates, err := convertDateTimeToString(value.(primitive.DateTime))
	if err != nil {
		return nil, err
	}

	return dates[0], nil
}

var DateTransformerInstance = newDateTransformer()
