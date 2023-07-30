package mongomodel

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Lyearn/backend-universe/packages/common/logger"
	"github.com/Lyearn/backend-universe/packages/observability/errorhandler"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BSONDocTranslateToEnum string

const (
	BSONDocTranslateToEnumMongo       BSONDocTranslateToEnum = "mongo"
	BSONDocTranslateToEnumEntityModel BSONDocTranslateToEnum = "entity_model"
)

func BuildBSONDoc(ctx context.Context, bsonDoc *bson.D, entityModelSchema *EntityModelSchema, translateTo BSONDocTranslateToEnum) error {
	if entityModelSchema == nil {
		return nil
	}

	if bsonDoc == nil && len(entityModelSchema.Root.Children) != 0 {
		return errors.New("BSON doc is nil but entity model schema is not empty")
	}

	if len(*bsonDoc) == 0 && len(entityModelSchema.Root.Children) == 0 {
		return nil
	}

	return buildBSONDoc(ctx, bsonDoc, &entityModelSchema.Root, translateTo)
}

func buildBSONDoc(ctx context.Context, bsonDocRef interface{}, schemaTreeNode *TreeNode, translateTo BSONDocTranslateToEnum) error {
	if schemaTreeNode == nil {
		return nil
	}

	switch schemaTreeNode.Props.Type {
	case reflect.Struct:
		bsonDoc, ok := bsonDocRef.(*bson.D)
		if !ok {
			logger.Error(ctx, fmt.Sprintf(
				"schema node is of type struct but bson node is not of type bson.D. Path: %s, Model key: %s, Schema key: %s",
				schemaTreeNode.Path, schemaTreeNode.Key, schemaTreeNode.BSONKey,
			))
			return fmt.Errorf("invalid bson node type for key %s", schemaTreeNode.Key)
		}

		if bsonDoc == nil {
			return errors.New("bson doc is nil")
		}

		for index, c := range schemaTreeNode.Children {
			schemaNode := c

			nodesMatch := true
			// nodes do not match if bson doc is shorter than schema tree node or
			// if bson doc key does not match schema tree node key
			if index >= len(*bsonDoc) || schemaNode.BSONKey != (*bsonDoc)[index].Key {
				nodesMatch = false
			}

			// Schema Options related logic starts here

			if !nodesMatch {
				isIDField := schemaNode.BSONKey == "_id"

				// throw error if schema node is not _id field (special field) and is required but has no default value.
				if !isIDField && schemaNode.Props.Options.Required && schemaNode.Props.Options.Default == nil {
					logger.Error(ctx, fmt.Sprintf(
						"schema and bson nodes are not equal. Path: %s, Model key: %s, Schema key: %s",
						schemaNode.Path, schemaNode.Key, schemaNode.BSONKey,
					))

					return fmt.Errorf("key %s not found at path %s bson doc", schemaNode.Key, schemaNode.Path)
				}

				var bsonNodeToAppend bson.E

				// add bson node with default value if value is available. else skip this schema node as it is not compulsory.
				// But _id is a special field and it needs to be populated with ObjectID if not available.
				if isIDField {
					var valueToAppend interface{}
					// populate _id field only if translating this doc to mongo doc.
					// in other cases, inserting a dummy _id node which will throw error in service logic.
					// reason being, while translating to entity model, if _id is not available, then this
					// logic will populate new objectId every time for the same object and cause FE unique key issue.
					// expectation here is if mgoID property of any field is changed to true in schema,
					// then it should be populated via script beforehand.
					if translateTo == BSONDocTranslateToEnumMongo {
						valueToAppend = primitive.NewObjectID()
					} else {
						valueToAppend = ""
					}
					bsonNodeToAppend = bson.E{Key: "_id", Value: valueToAppend}
				} else {
					bsonNodeToAppend = bson.E{Key: schemaNode.BSONKey, Value: schemaNode.Props.Options.Default}
				}

				if index >= len(*bsonDoc) {
					// adding bson node with default value at the end of bson doc
					*bsonDoc = append(*bsonDoc, bsonNodeToAppend)
					continue
				} else {
					// adding bson node with default value at current index
					*bsonDoc = append((*bsonDoc)[:index+1], (*bsonDoc)[index:]...)
					(*bsonDoc)[index] = bsonNodeToAppend
				}

				continue
			}

			bsonNode := (*bsonDoc)[index]

			convertedVal, err := getConvertedValueForNode(ctx, bsonNode.Value, &schemaNode, translateTo)
			if err != nil {
				return err
			}

			bsonNode.Value = convertedVal
			(*bsonDoc)[index] = bsonNode
		}

	case reflect.Slice:
		bsonDoc, ok := bsonDocRef.(*bson.A)
		if !ok || bsonDoc == nil {
			return errors.New("bson doc is nil but entity model schema is not empty")
		}

		schemaNodeSliceElem := schemaTreeNode.Children[0]

		for index, elem := range *bsonDoc {
			bsonNodeSliceElem, err := getConvertedValueForNode(ctx, elem, &schemaNodeSliceElem, translateTo)
			if err != nil {
				return err
			}

			(*bsonDoc)[index] = bsonNodeSliceElem
		}

	// default case handles all primitive types i.e. all leaf nodes of schema tree
	default:
		// Transformations related logic starts here

		if len(schemaTreeNode.Props.Transformers) != 0 {
			for _, transformer := range schemaTreeNode.Props.Transformers {
				if transformer == nil {
					continue
				}

				var elemVal interface{}
				if _, ok := bsonDocRef.(*interface{}); !ok {
					// this case handles only elements of array which are passed as reference from the above *bson.A case.
					// hence, reject any other type.
					return nil
				} else {
					elemVal = *(bsonDocRef.(*interface{}))
				}

				var modifiedBSONNodeVal interface{}
				var err error

				switch translateTo {
				case BSONDocTranslateToEnumMongo:
					modifiedBSONNodeVal, err = transformer.TransformForMongoDoc(elemVal)
				case BSONDocTranslateToEnumEntityModel:
					modifiedBSONNodeVal, err = transformer.TransformForEntityModelDoc(elemVal)
				default:
					err = fmt.Errorf("unknown translate from enum value %s", translateTo)
				}

				if err != nil {
					return err
				}

				*(bsonDocRef.(*interface{})) = modifiedBSONNodeVal
			}
		}
	}

	return nil
}

func getConvertedValueForNode(
	ctx context.Context,
	nodeVal interface{},
	schemaTreeNode *TreeNode,
	translateTo BSONDocTranslateToEnum,
) (interface{}, error) {
	var modifiedVal interface{}
	var err error

	// this switch case provides type support for bson.D and bson.A type of elements.
	// without this, *interface{} type of bsonDoc would be passed in the recursive call,
	// which will then go to the default case and will not be able to handle any nested type.
	switch typedValue := nodeVal.(type) {
	case bson.D:
		err = buildBSONDoc(ctx, &typedValue, schemaTreeNode, translateTo)
		modifiedVal = typedValue

	case bson.A:
		err = buildBSONDoc(ctx, &typedValue, schemaTreeNode, translateTo)
		modifiedVal = typedValue

	case interface{}:
		err = buildBSONDoc(ctx, &typedValue, schemaTreeNode, translateTo)
		modifiedVal = typedValue

	default:
		errorParams := map[string]interface{}{
			"key": schemaTreeNode.BSONKey,
			"val": typedValue,
		}
		return nil, errorhandler.NewBadRequestError(errorhandler.CommonErrorProps{
			Message: "Invalid bson doc type",
			Where:   "mongomodel.getConvertedValueForNode",
			Params:  &errorParams,
		})
	}

	if err != nil {
		return err, err
	}

	return modifiedVal, err
}
