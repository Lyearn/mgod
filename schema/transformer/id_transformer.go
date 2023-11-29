package transformer

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IDTransformer struct{}

func newIDTransformer() Transformer {
	return &IDTransformer{}
}

func (t IDTransformer) isTransformationRequired(field reflect.StructField) bool {
	return field.Tag.Get("mgoType") == "id"
}

func (t IDTransformer) TransformForMongoDoc(value interface{}) (interface{}, error) {
	objectIDs, err := convertStringToObjectID(value.(string))
	if err != nil {
		return nil, err
	}

	return objectIDs[0], nil
}

func (t IDTransformer) TransformForEntityModelDoc(value interface{}) (interface{}, error) {
	ids, err := convertObjectIDToString(value.(primitive.ObjectID))
	if err != nil {
		return nil, err
	}

	return ids[0], nil
}

var IDTransformerInstance = newIDTransformer()
