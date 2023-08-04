package metafield

import (
	"go.mongodb.org/mongo-driver/bson"
)

func ValidatedAndAddFieldValue(doc *bson.D, metaField MetaField) error {
	field := string(metaField.GetMetaFieldKey())

	for index, elem := range *doc {
		if elem.Key != field {
			continue
		}

		// field is already present with expected type.
		if metaField.CheckIfValidValue(elem.Value) {
			metaField.FieldAlreadyPresent(doc, index)
			return nil
		}

		// field is already present but not of the expected type.
		// hence, assigning the expected value to it.
		if err := metaField.FieldPresentWithIncorrectVal(doc, index); err != nil {
			return err
		}

		return nil
	}

	// field is not present in the existing doc. need to add it.
	metaField.FieldNotPresent(doc)

	return nil
}
