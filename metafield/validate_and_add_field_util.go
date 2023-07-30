package metafield

import (
	"go.mongodb.org/mongo-driver/bson"
)

func ValidatedAndAddFieldValue(doc *bson.D, metaField MetaField) {
	field := string(metaField.GetMetaFieldKey())

	for index, elem := range *doc {
		if elem.Key != field {
			continue
		}

		// field is already present with expected type.
		if metaField.CheckIfValidValue(elem.Value) {
			metaField.FieldAlreadyPresent(doc, index)
			return
		}

		// field is already present but not of the expected type.
		// hence, assigning the expected value to it.
		metaField.FieldPresentWithIncorrectVal(doc, index)

		return
	}

	// field is not present in the existing doc. need to add it.
	metaField.FieldNotPresent(doc)
}
