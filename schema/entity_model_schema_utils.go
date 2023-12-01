package schema

import (
	"reflect"

	"github.com/Lyearn/mgod/schema/fieldopt"
)

// GetSchemaNameForModel returns the default schema name for the model.
func GetSchemaNameForModel[T any](model T) string {
	return reflect.TypeOf(model).Name()
}

// GetDefaultSchemaTreeRootNode returns the default root node of a schema tree.
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

// GetPathForField returns the schema tree path for the field.
func GetPathForField(field, parent string) string {
	path := field
	if parent != "" {
		path = parent + "." + field
	}

	return path
}
