// Package schema provides utilities to build the schema tree for the entity model.
package schema

import (
	"reflect"

	"github.com/Lyearn/mgod/schema/fieldopt"
	"github.com/Lyearn/mgod/schema/metafield"
	"github.com/Lyearn/mgod/schema/schemaopt"
	"github.com/Lyearn/mgod/schema/transformer"
	"github.com/samber/lo"
)

// EntityModelSchema holds the schema tree for the entity model.
type EntityModelSchema struct {
	// Root is the root node of the schema tree.
	// Root node is a dummy node, it's not a real field in the model. Actual doc parsing starts from the children of root node.
	Root TreeNode

	// Nodes is a map used to quickly access the [TreeNode] by path.
	Nodes map[string]*TreeNode
}

// TreeNode is a node in the schema tree.
type TreeNode struct {
	// Path is used to identify the ancestor chain. Used for debugging purposes.
	Path string
	// BSONKey is the translated bson key.
	BSONKey string
	// Key is the struct field name. Used for debugging purposes.
	Key string
	// Props contains the field properties.
	Props SchemaFieldProps
	// Children contains the child nodes.
	// Array is used instead of map to preserve the order of fields. Fields in bson doc should always match with the schema tree order.
	Children []TreeNode
}

// SchemaFieldProps are the possible field properties.
type SchemaFieldProps struct {
	// Type holds the struct field type or the underlying type in case of pointer.
	Type reflect.Kind
	// IsPointer is used to identify pointer type of fields.
	IsPointer bool
	// Transformers are the transformations that needs to be applied on the field while building the bson doc.
	Transformers []transformer.Transformer
	// Options are the schema options for the field.
	Options fieldopt.SchemaFieldOptions
}

// BuildSchemaForModel builds the schema tree for the given model.
func BuildSchemaForModel[T any](model T, schemaOpts schemaopt.SchemaOptions) (*EntityModelSchema, error) {
	schemaTree := make([]TreeNode, 0)
	rootNode := GetDefaultSchemaTreeRootNode()

	nodes := make(map[string]*TreeNode)
	nodes[rootNode.Path] = &rootNode

	opts := NewEntityModelSchemaOptions().SetXIDRequired(rootNode.Props.Options.XID)
	err := buildSchema(model, &schemaTree, nodes, rootNode.BSONKey, *opts)
	if err != nil {
		return nil, err
	}

	addMetaFields(model, schemaOpts, &schemaTree, nodes, rootNode.BSONKey)

	rootNode.Children = schemaTree

	schema := &EntityModelSchema{
		Root:  rootNode,
		Nodes: nodes,
	}

	return schema, nil
}

// buildSchema is a recursive function that actually builds the schema tree for the given model.
func buildSchema[T any](model T, treeRef *[]TreeNode, nodes map[string]*TreeNode, parent string, opts EntityModelSchemaOptions) error {
	v := reflect.ValueOf(model)

	// converting pointer to value
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	currentLevelBSONFields := getCurrentLevelBSONFields(v)
	xidFound := false

	for i := 0; i < v.NumField(); i++ {
		structField := v.Type().Field(i)
		fieldName := getBSONFieldName(structField)
		if fieldName == "" {
			// skipping unexported fields
			continue
		}

		// if we are processing an inline struct, then need to check if parent already contains
		// the current field. if yes, then we can skip processing the field because bson util
		// internally merges only unique child fields into parent in case of inline structs.
		if opts.bsonInlineParent && lo.Contains(opts.parentBSONFields, fieldName) {
			continue
		}

		// Custom _id field handling
		if fieldName == "_id" {
			xidFound = true
		}

		// Field level changes starts here

		isPointerTypeField := false
		if structField.Type.Kind() == reflect.Ptr {
			isPointerTypeField = true

			// setting the struct field type to the underlying type of the pointer
			elem := structField.Type.Elem()
			structField.Type = elem
		}

		transformers := transformer.GetRequiredTransformersForField(structField)
		options, err := fieldopt.GetSchemaOptionsForField(structField)
		if err != nil {
			return err
		}

		path := GetPathForField(fieldName, parent)

		treeNode := TreeNode{
			Path:    path,
			BSONKey: fieldName,
			Key:     structField.Name,
			Props: SchemaFieldProps{
				Type:         structField.Type.Kind(),
				IsPointer:    isPointerTypeField,
				Transformers: transformers,
				Options:      options,
			},
		}

		// Child level changes starts here

		var recurseErr error

		//nolint:exhaustive // need to handle only complex object types
		switch structField.Type.Kind() {
		case reflect.Struct:
			var field reflect.Value

			if !isPointerTypeField {
				field = v.Field(i)
			} else {
				// need to create a new struct instance for pointer type fields
				field = reflect.New(structField.Type)
			}

			if isBSONInlineField(structField) {
				toAppendTreeNodes := make([]TreeNode, 0)

				// combining all ancestor fields for current child
				existingBSONFields := currentLevelBSONFields
				if opts.parentBSONFields != nil {
					existingBSONFields = append(existingBSONFields, opts.parentBSONFields...)
				}

				opts := NewEntityModelSchemaOptions().SetBSONInlineParent(true).SetParentBSONFields(existingBSONFields)
				inlineFieldsErr := buildSchema(field.Interface(), &toAppendTreeNodes, nodes, parent, *opts)
				if inlineFieldsErr != nil {
					return inlineFieldsErr
				}

				addTreeNodesToSchema(treeRef, nodes, toAppendTreeNodes...)

				// no need to add inline struct as a child node
				continue
			}

			recurseErr = handleStructTypeField(field, &treeNode, nodes, path)

		case reflect.Slice:
			recurseErr = handleSliceTypeField(structField.Type.Elem(), &treeNode, nodes, path)
		}

		if recurseErr != nil {
			return recurseErr
		}

		addTreeNodesToSchema(treeRef, nodes, treeNode)
	}

	// if _id is not found and is required for model, then insert it at the beginning.
	if opts.xidRequired && !xidFound {
		xidField := reflect.StructField{
			Name: "XID",
			Type: reflect.TypeOf(""),
			Tag:  `bson:"_id" mgoType:"id"`,
		}

		fieldName := getBSONFieldName(xidField)
		transformers := transformer.GetRequiredTransformersForField(xidField)
		options, err := fieldopt.GetSchemaOptionsForField(xidField)
		if err != nil {
			return err
		}

		path := GetPathForField(fieldName, parent)

		xidNode := TreeNode{
			Path:    path,
			BSONKey: fieldName,
			Key:     xidField.Name,
			Props: SchemaFieldProps{
				Type:         reflect.String,
				Transformers: transformers,
				Options:      options,
			},
		}

		addTreeNodesToSchema(treeRef, nodes, xidNode)
	}

	return nil
}

func handleStructTypeField(field reflect.Value, treeNode *TreeNode, nodes map[string]*TreeNode, path string) error {
	opts := NewEntityModelSchemaOptions().SetXIDRequired(treeNode.Props.Options.XID)

	return buildSchema(field.Interface(), &treeNode.Children, nodes, path, *opts)
}

func handleSliceTypeField(sliceElemType reflect.Type, treeNode *TreeNode, nodes map[string]*TreeNode, path string) error {
	// In case of slice, transformations are applicable on the slice elements only,
	// whereas options are applicable on the slice itself.
	parentTransformers := treeNode.Props.Transformers
	treeNode.Props.Transformers = []transformer.Transformer{}

	// $ is used to denote the slice elements.
	path += ".$"

	// if slice element is a pointer, then we need to get the underlying type first.
	if sliceElemType.Kind() == reflect.Pointer {
		sliceElemType = sliceElemType.Elem()
	}

	// slice will only have one child, which will be the slice element.
	childNode := TreeNode{
		Path:    path,
		BSONKey: "$",
		Key:     "$",
		Props: SchemaFieldProps{
			Type:         sliceElemType.Kind(),
			Transformers: parentTransformers,
		},
	}

	addTreeNodesToSchema(&treeNode.Children, nodes, childNode)

	// if slice element is a struct, then we need to recurse.
	if sliceElemType.Kind() == reflect.Struct {
		// creating a new instance of slice element type to pass to buildSchema
		sliceElemModel := reflect.New(sliceElemType).Interface()

		opts := NewEntityModelSchemaOptions().SetXIDRequired(treeNode.Props.Options.XID)
		err := buildSchema(sliceElemModel, &treeNode.Children[0].Children, nodes, path, *opts)
		if err != nil {
			return err
		}
	}

	return nil
}

// addMetaFields adds meta type fields to the schema tree so that the bson doc can be built without any errors
// of fields not found in the tree (Meta fields are appended to the bson doc based on the schema options dynamically).
func addMetaFields[T any](model T, schemaOptions schemaopt.SchemaOptions, treeRef *[]TreeNode, nodes map[string]*TreeNode, parent string) {
	v := reflect.ValueOf(model)

	// converting pointer to value
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	rootStructFields := getCurrentLevelBSONFields(v)

	for _, metaField := range metafield.GetAvailableMetaFields() {
		if !metaField.IsApplicable(schemaOptions) {
			continue
		}

		if lo.Contains(rootStructFields, string(metaField.GetKey())) {
			// meta field is already present in the model, so no need to add it.
			continue
		}

		field := string(metaField.GetKey())
		path := GetPathForField(field, parent)

		// append meta field in the schema tree.
		toAppendTreeNode := TreeNode{
			Path:    path,
			BSONKey: field,
			Key:     field,
			Props: SchemaFieldProps{
				Type: metaField.GetReflectKind(),
				Options: fieldopt.SchemaFieldOptions{
					// not keeping meta fields as required to prevent any errors while building the bson doc.
					// meta fields are always added if enabled in schema options and not present in the bson doc.
					Required: false,
				},
				Transformers: metaField.GetApplicableTransformers(),
			},
		}

		addTreeNodesToSchema(treeRef, nodes, toAppendTreeNode)
	}
}

// addTreeNodesToSchema adds the given tree nodes to the schema tree as well as to the nodes map.
func addTreeNodesToSchema(treeRef *[]TreeNode, nodes map[string]*TreeNode, toAddTreeNodes ...TreeNode) {
	*treeRef = append(*treeRef, make([]TreeNode, len(toAddTreeNodes))...)

	// assigning the address of the new nodes to the nodes map.
	// need to do this after appending to the treeRef, because the addresses of the appended nodes will change.
	for i := range toAddTreeNodes {
		parentIdx := len(*treeRef) - len(toAddTreeNodes) + i

		(*treeRef)[parentIdx] = toAddTreeNodes[i]
		nodes[toAddTreeNodes[i].Path] = &(*treeRef)[parentIdx]
	}
}
