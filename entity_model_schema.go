package mongomodel

import (
	"reflect"
	"strings"

	"github.com/Lyearn/backend-universe/packages/common/util/typeutil"
	"github.com/Lyearn/backend-universe/packages/store/mongomodel/schemaopt"
	"github.com/Lyearn/backend-universe/packages/store/mongomodel/transformer"
	"github.com/samber/lo"
)

type EntityModelSchema struct {
	// root node is a dummy node, it's not a real field in the model.
	// actual doc parsing starts from the children of root node.
	Root TreeNode

	// nodes map is used to quickly access the schema tree node by path.
	Nodes map[string]*TreeNode
}

type TreeNode struct {
	Path    string // path will be used to identify the ancestor chain. used for debugging purposes.
	BSONKey string // translated bson key.
	Key     string // struct key. used for debugging purposes.
	Props   SchemaFieldProps
	// array is used instead of map to preserve the order of fields.
	// fields in bson doc should always match with the schema tree order.
	Children []TreeNode
}

type SchemaFieldProps struct {
	Type         reflect.Kind              // contains struct field type or the underlying type in case of pointer.
	IsPointer    bool                      // will be used to identify pointer type of fields.
	Transformers []transformer.Transformer // reference to id, date, etc. transformers
	Options      schemaopt.SchemaFieldOptions
}

func BuildSchemaForModel[T any](model T) (*EntityModelSchema, error) {
	schemaTree := make([]TreeNode, 0)
	rootNode := GetDefaultSchemaTreeRootNode()

	nodes := make(map[string]*TreeNode)
	nodes[rootNode.Path] = &rootNode

	opts := NewEntityModelSchemaOptions().SetXIDRequired(rootNode.Props.Options.XID)
	err := buildSchema(model, &schemaTree, nodes, rootNode.BSONKey, *opts)
	if err != nil {
		return nil, err
	}

	rootNode.Children = schemaTree

	schema := &EntityModelSchema{
		Root:  rootNode,
		Nodes: nodes,
	}

	return schema, nil
}

func GetDefaultSchemaTreeRootNode() TreeNode {
	rootNode := TreeNode{
		Path:    "$root",
		Key:     "$root",
		BSONKey: "$root",
		Props: SchemaFieldProps{
			Type: reflect.Struct,
			Options: schemaopt.SchemaFieldOptions{
				// _id is required by default at root of the doc
				XID: true,
			},
		},
	}

	return rootNode
}

func buildSchema[T any](model T, treeRef *[]TreeNode, nodes map[string]*TreeNode, parent string, opts EntityModelSchemaOptions) error {
	v := reflect.ValueOf(model)

	// converting pointer to value
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	currentLevelBSONFields := GetCurrentLevelBSONFields(v)
	xidFound := false

	for i := 0; i < v.NumField(); i++ {
		structField := v.Type().Field(i)
		fieldName := typeutil.GetBSONFieldName(structField)

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
		options, err := schemaopt.GetSchemaOptionsForField(structField)
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

			if IsBSONInlineField(structField) {
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

				AddTreeNodesToSchema(treeRef, nodes, toAppendTreeNodes...)

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

		AddTreeNodesToSchema(treeRef, nodes, treeNode)
	}

	// if _id is not found and is required for model, then insert it at the beginning.
	if opts.xidRequired && !xidFound {
		xidField := reflect.StructField{
			Name: "XID",
			Type: reflect.TypeOf(""),
			Tag:  `bson:"_id" mgoType:"id"`,
		}

		fieldName := typeutil.GetBSONFieldName(xidField)
		transformers := transformer.GetRequiredTransformersForField(xidField)
		options, err := schemaopt.GetSchemaOptionsForField(xidField)
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

		AddTreeNodesToSchema(treeRef, nodes, xidNode)
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

	AddTreeNodesToSchema(&treeNode.Children, nodes, childNode)

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

// AddTreeNodesToSchema adds the given tree nodes to the schema tree as well as to the nodes map.
func AddTreeNodesToSchema(treeRef *[]TreeNode, nodes map[string]*TreeNode, toAddTreeNodes ...TreeNode) {
	*treeRef = append(*treeRef, make([]TreeNode, len(toAddTreeNodes))...)

	// assigning the address of the new nodes to the nodes map.
	// need to do this after appending to the treeRef, because the addresses of the appended nodes will change.
	for i := range toAddTreeNodes {
		parentIdx := len(*treeRef) - len(toAddTreeNodes) + i

		(*treeRef)[parentIdx] = toAddTreeNodes[i]
		nodes[toAddTreeNodes[i].Path] = &(*treeRef)[parentIdx]
	}
}

func GetCurrentLevelBSONFields(s reflect.Value) []string {
	currentLevelBSONFields := make([]string, 0)

	for i := 0; i < s.NumField(); i++ {
		structField := s.Type().Field(i)
		fieldName := typeutil.GetBSONFieldName(structField)

		currentLevelBSONFields = append(currentLevelBSONFields, fieldName)
	}

	return currentLevelBSONFields
}

func IsBSONInlineField(field reflect.StructField) bool {
	bsonTagVal := field.Tag.Get("bson")
	if bsonTagVal == "" || bsonTagVal == "-" {
		return false
	}

	tagValues := strings.Split(bsonTagVal, ",")
	flags := tagValues[1:]

	for _, flag := range flags {
		if flag == "inline" {
			return true
		}
	}

	return false
}

func GetPathForField(field, parent string) string {
	path := field
	if parent != "" {
		path = parent + "." + field
	}

	return path
}
