package transformer

import "reflect"

// Transformer can transform fields in both directions i.e. from entity model to mongo doc and vice versa.
type Transformer interface {
	// IsTransformationRequired reports whether the transformer is required for the given field.
	IsTransformationRequired(field reflect.StructField) bool
	// TransformForMongoDoc transforms the incoming value according to mongo requirements.
	TransformForMongoDoc(value interface{}) (interface{}, error)
	// TransformForEntityModelDoc transforms the incoming value according to entity model requirements.
	TransformForEntityModelDoc(value interface{}) (interface{}, error)
}

var availableTransformers = []Transformer{
	IDTransformer,
	DateTransformer,
}

// GetRequiredTransformersForField returns the transformers required for the given field.
func GetRequiredTransformersForField(field reflect.StructField) []Transformer {
	transformers := []Transformer{}

	for _, transformer := range availableTransformers {
		if transformer.IsTransformationRequired(field) {
			transformers = append(transformers, transformer)
		}
	}

	return transformers
}
