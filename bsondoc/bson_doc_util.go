package bsondoc

import (
	"go.mongodb.org/mongo-driver/bson"
)

// GetFieldValueFromRootDoc returns the value of the provided field (nil if not found) from the root of a bson doc.
func GetFieldValueFromRootDoc(doc *bson.D, field string) interface{} {
	if doc == nil {
		return nil
	}

	for _, elem := range *doc {
		if elem.Key == field {
			return elem.Value
		}
	}

	return nil
}
