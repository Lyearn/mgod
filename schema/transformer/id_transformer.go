package transformer

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type idTransformer struct{}

func newIDTransformer() Transformer {
	return &idTransformer{}
}

func (t idTransformer) isTransformationRequired(field reflect.StructField) bool {
	return field.Tag.Get("mgoType") == "id"
}

func (t idTransformer) TransformForMongoDoc(value interface{}) (interface{}, error) {
	objectIDs, err := convertStringToObjectID(value.(string))
	if err != nil {
		return nil, err
	}

	return objectIDs[0], nil
}

func (t idTransformer) TransformForEntityModelDoc(value interface{}) (interface{}, error) {
	ids, err := convertObjectIDToString(value.(primitive.ObjectID))
	if err != nil {
		return nil, err
	}

	return ids[0], nil
}

var IDTransformerInstance = newIDTransformer()
