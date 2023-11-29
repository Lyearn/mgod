package schema

import (
	"reflect"

	"github.com/Lyearn/mgod/schema/fieldopt"
)

func GetSchemaNameForModel[T any](model T) string {
	return reflect.TypeOf(model).Name()
}

func GetDefaultSchemaTreeRootNode() TreeNode {
	rootNode := TreeNode{
		Path:    "$root",
		Key:     "$root",
		BSONKey: "$root",
		Props: SchemaFieldProps{
			Type: reflect.Struct,
			Options: fieldopt.SchemaFieldOptions{
				// _id is required by default at root of the doc
				XID: true,
			},
		},
	}

	return rootNode
}

func GetPathForField(field, parent string) string {
	path := field
	if parent != "" {
		path = parent + "." + field
	}

	return path
}
