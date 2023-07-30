package transformer

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DateTransformer struct{}

func newDateTransformer() Transformer {
	return &DateTransformer{}
}

func (t DateTransformer) isTransformationRequired(field reflect.StructField) bool {
	return field.Tag.Get("mgoType") == "date"
}

func (t DateTransformer) TransformForMongoDoc(value interface{}) (interface{}, error) {
	primitiveDates, err := convertStringToDateTime(value.(string))
	if err != nil {
		return nil, err
	}

	return primitiveDates[0], nil
}

func (t DateTransformer) TransformForEntityModelDoc(value interface{}) (interface{}, error) {
	dates, err := convertDateTimeToString(value.(primitive.DateTime))
	if err != nil {
		return nil, err
	}

	return dates[0], nil
}

var DateTransformerInstance = newDateTransformer()
