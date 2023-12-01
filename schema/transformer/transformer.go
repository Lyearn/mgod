package transformer

import "reflect"

// Transformer can transform fields in both directions i.e. from entity model to mongo doc and vice versa.
type Transformer interface {
	isTransformationRequired(field reflect.StructField) bool
	// TransformForMongoDoc transforms the incoming value according to mongo requirements.
	TransformForMongoDoc(value interface{}) (interface{}, error)
	// TransformForEntityModelDoc transforms the incoming value according to entity model requirements.
	TransformForEntityModelDoc(value interface{}) (interface{}, error)
}

var availableTransformers = []Transformer{
	IDTransformerInstance,
	DateTransformerInstance,
}

func GetRequiredTransformersForField(field reflect.StructField) []Transformer {
	transformers := []Transformer{}

	for _, transformer := range availableTransformers {
		if transformer.isTransformationRequired(field) {
			transformers = append(transformers, transformer)
		}
	}

	return transformers
}
