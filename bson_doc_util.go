package mongomodel

import "go.mongodb.org/mongo-driver/bson"

func GetFieldValueFromBSONRootDoc(doc *bson.D, field string) interface{} {
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
