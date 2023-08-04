package mongomodel

import (
	"context"
	"errors"
	"fmt"

	"github.com/Lyearn/backend-universe/packages/common/logger"
	"github.com/Lyearn/backend-universe/packages/observability/errorhandler"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BSONDocTranslateToEnum string

const (
	BSONDocTranslateToEnumMongo       BSONDocTranslateToEnum = "mongo"
	BSONDocTranslateToEnumEntityModel BSONDocTranslateToEnum = "entity_model"
)

func BuildBSONDoc(
	ctx context.Context,
	bsonDoc *bson.D,
	entityModelSchema *EntityModelSchema,
	translateTo BSONDocTranslateToEnum,
) error {
	if entityModelSchema == nil {
		return nil
	}

	if bsonDoc == nil && len(entityModelSchema.Root.Children) != 0 {
		return errors.New("BSON doc is nil but entity model schema is not empty")
	}

	if len(*bsonDoc) == 0 && len(entityModelSchema.Root.Children) == 0 {
		return nil
	}

	err := buildBSONDoc(ctx, bsonDoc, entityModelSchema.Nodes, entityModelSchema.Root.Path, translateTo)
	if err != nil {
		return err
	}

	return nil
}

func buildBSONDoc(
	ctx context.Context,
	bsonDocRef interface{},
	schemaNodes map[string]*TreeNode,
	parent string,
	translateTo BSONDocTranslateToEnum,
) error {
	if bsonDocRef == nil {
		return nil
	}

	schemaNode, err := getSchemaNodeForPath(ctx, parent, schemaNodes, translateTo)
	if err != nil {
		return err
	} else if schemaNode == nil {
		return nil
	}

	switch bsonElem := bsonDocRef.(type) {
	case *bson.D:
		if bsonElem == nil {
			return nil
		}

		visitedSchemaNodes := make([]string, 0)

		for bsonIdx, bsonNode := range *bsonElem {
			nodePath := GetPathForField(bsonNode.Key, parent)
			visitedSchemaNodes = append(visitedSchemaNodes, nodePath)

			convertedValue, err := getConvertedValueForNode(ctx, bsonNode.Value, schemaNodes, nodePath, translateTo)
			if err != nil {
				return err
			}

			bsonNode.Value = convertedValue

			(*bsonElem)[bsonIdx] = bsonNode
		}

		// check if there are any missing nodes in the bson doc at the current level as compared to the schema.
		immediateChildren := lo.Map(schemaNode.Children, func(child TreeNode, _ int) string {
			return child.Path
		})
		uniqVisitedSchemaNodes := lo.Uniq(visitedSchemaNodes)

		if len(uniqVisitedSchemaNodes) != len(immediateChildren) {
			err := addMissingNodes(ctx, bsonElem, immediateChildren, uniqVisitedSchemaNodes, schemaNodes, translateTo)
			if err != nil {
				return err
			}
		}

	case *bson.A:
		if bsonElem == nil || len(*bsonElem) == 0 {
			return nil
		}

		// array elements are represented as $.
		nodePath := GetPathForField("$", parent)

		for arrIdx := range *bsonElem {
			elemVal := (*bsonElem)[arrIdx]

			convertedValue, err := getConvertedValueForNode(ctx, elemVal, schemaNodes, nodePath, translateTo)
			if err != nil {
				return err
			}

			(*bsonElem)[arrIdx] = convertedValue
		}

	// default case handles all primitive types i.e. all leaf nodes of schema tree or all bson doc
	// elements which are not of type bson.D or bson.A.
	default:
		// Transformations related logic starts here

		if len(schemaNode.Props.Transformers) == 0 {
			return nil
		}

		for _, transformer := range schemaNode.Props.Transformers {
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
				err = fmt.Errorf("unknown translateTo enum value %s", translateTo)
			}

			if err != nil {
				return err
			}

			*(bsonDocRef.(*interface{})) = modifiedBSONNodeVal
		}
	}

	return nil
}

func getConvertedValueForNode(
	ctx context.Context,
	nodeVal interface{},
	schemaNodes map[string]*TreeNode,
	parent string,
	translateTo BSONDocTranslateToEnum,
) (interface{}, error) {
	var modifiedVal interface{}
	var err error

	// this switch case provides type support for bson.D and bson.A type of elements.
	// without this, *interface{} type of bsonDoc would be passed in the recursive call,
	// which will then go to the default case and will not be able to handle any nested type.
	switch typedValue := nodeVal.(type) {
	case bson.D:
		err = buildBSONDoc(ctx, &typedValue, schemaNodes, parent, translateTo)
		modifiedVal = typedValue

	case bson.A:
		err = buildBSONDoc(ctx, &typedValue, schemaNodes, parent, translateTo)
		modifiedVal = typedValue

	case interface{}:
		err = buildBSONDoc(ctx, &typedValue, schemaNodes, parent, translateTo)
		modifiedVal = typedValue

	default:
		errorParams := map[string]interface{}{
			"doc": typedValue,
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

// addMissingNodes appends missing nodes in bson doc which have default value.
func addMissingNodes(
	ctx context.Context,
	bsonElem *bson.D,
	immediateChildren []string,
	uniqVisitedSchemaNodes []string,
	schemaNodes map[string]*TreeNode,
	translateTo BSONDocTranslateToEnum,
) error {
	missingSchemaPaths, _ := lo.Difference(immediateChildren, uniqVisitedSchemaNodes)
	for _, missingSchemaPath := range missingSchemaPaths {
		missingSchemaNode, err := getSchemaNodeForPath(ctx, missingSchemaPath, schemaNodes, translateTo)
		if err != nil {
			return err
		} else if missingSchemaNode == nil {
			continue
		}

		// skip the node if it is not required and has no default value
		if !missingSchemaNode.Props.Options.Required && missingSchemaNode.Props.Options.Default == nil {
			continue
		}

		isIDField := missingSchemaNode.BSONKey == "_id"

		// throw error if schema node is not _id field (special field) and is required but has no default value.
		if !isIDField && missingSchemaNode.Props.Options.Default == nil {
			return fmt.Errorf("required field at path %s is missing in bson doc", missingSchemaPath)
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
			bsonNodeToAppend = bson.E{Key: missingSchemaNode.BSONKey, Value: missingSchemaNode.Props.Options.Default}
		}

		// by default, add missing nodes at the end of bson doc.
		(*bsonElem) = append(*bsonElem, bsonNodeToAppend)
	}

	return nil
}

func getSchemaNodeForPath(ctx context.Context, path string, schemaNodes map[string]*TreeNode, translateTo BSONDocTranslateToEnum) (*TreeNode, error) {
	schemaNode, ok := schemaNodes[path]
	if !ok {
		// TODO: remove this check once all schemas are in sync.
		// skip throwing error for nodes which are not present in actual entity schema but present in mongo doc.
		if translateTo == BSONDocTranslateToEnumEntityModel {
			return nil, nil
		}

		logger.Error(ctx, fmt.Sprintf(
			"schema doesn't contains any node at path %s found in bsonDoc", path))

		return nil, fmt.Errorf("unknown path %s found in bson doc", path)
	}

	return schemaNode, nil
}
