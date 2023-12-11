package mgod

import (
	"context"
	"fmt"
	"strings"

	"github.com/Lyearn/mgod/bsondoc"
	"github.com/Lyearn/mgod/errors"
	"github.com/Lyearn/mgod/schema"
	"github.com/Lyearn/mgod/schema/metafield"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetSchemaCacheKey returns the cache key for the schema of a model.
// The cache key is in the format of <collection_name>_<model_name_or_value>.
// Unless we are storing the schema of a union type based on the discriminator value, the second parameter
// will be the model name only.
func GetSchemaCacheKey(collName, modelNameOrVal string) string {
	keyElems := []string{collName, modelNameOrVal}
	key := strings.Join(keyElems, "_")

	return key
}

func (m entityMongoModel[T]) getEntityModel() T {
	return m.modelType
}

// getMongoDocFromEntityModel converts the provided entity model to a bson.D doc.
func (m entityMongoModel[T]) getMongoDocFromEntityModel(ctx context.Context, model T) (bson.D, error) {
	marshalledDoc, err := bson.Marshal(model)
	if err != nil {
		return nil, err
	}

	var bsonDoc bson.D
	err = bson.Unmarshal(marshalledDoc, &bsonDoc)
	if err != nil {
		return nil, err
	}

	if bsonDoc == nil {
		// empty bson doc
		return bsonDoc, nil
	}

	if err = metafield.AddMetaFields(&bsonDoc, m.schemaOpts); err != nil {
		return nil, err
	}

	err = bsondoc.Build(ctx, &bsonDoc, m.schema, bsondoc.TranslateToEnumMongo)
	if err != nil {
		return nil, err
	}

	if m.isUnionType {
		discriminatorVal := bsondoc.GetFieldValueFromRootDoc(&bsonDoc, m.discriminatorKey)
		if discriminatorVal == nil {
			discriminatorVal = schema.GetSchemaNameForModel(m.modelType)
			bsonDoc = append(bsonDoc, primitive.E{
				Key:   m.discriminatorKey,
				Value: discriminatorVal,
			})
		}

		cacheKey := GetSchemaCacheKey(m.coll.Name(), discriminatorVal.(string))
		if _, err := schema.EntityModelSchemaCacheInstance.GetSchema(cacheKey); err != nil {
			schema.EntityModelSchemaCacheInstance.SetSchema(cacheKey, m.schema)
		}
	}

	return bsonDoc, nil
}

// getEntityModelFromMongoDoc converts the provided bson.D doc to an entity model.
func (m entityMongoModel[T]) getEntityModelFromMongoDoc(ctx context.Context, bsonDoc bson.D) (T, error) {
	model := m.getEntityModel()

	if bsonDoc == nil {
		// empty bson doc
		return model, nil
	}

	entityModelSchema := m.schema

	if m.isUnionType {
		discriminatorVal := bsondoc.GetFieldValueFromRootDoc(&bsonDoc, m.discriminatorKey)
		if discriminatorVal != nil {
			cacheKey := GetSchemaCacheKey(m.coll.Name(), discriminatorVal.(string))
			if unionElemSchema, err := schema.EntityModelSchemaCacheInstance.GetSchema(cacheKey); err == nil {
				entityModelSchema = unionElemSchema
			}
		}
	}

	err := bsondoc.Build(ctx, &bsonDoc, entityModelSchema, bsondoc.TranslateToEnumEntityModel)
	if err != nil {
		return model, err
	}

	marshalledDoc, err := bson.Marshal(bsonDoc)
	if err != nil {
		return model, err
	}

	err = bson.Unmarshal(marshalledDoc, &model)
	if err != nil {
		return model, err
	}

	return model, nil
}

// handleTimestampsForUpdateQuery adds updatedAt field to the update query if the schema options has timestamps enabled.
func (m entityMongoModel[T]) handleTimestampsForUpdateQuery(update interface{}, funcName string) (interface{}, error) {
	updateQuery, ok := update.(bson.D)
	if !ok {
		return nil, errors.NewBadRequestError(errors.BadRequestError{
			Underlying: "update query",
			Got:        fmt.Sprintf("%T", update),
			Expected:   "bson.D",
		})
	}

	if m.schemaOpts.Timestamps {
		updatedAtCommand := bson.E{
			Key: "$currentDate",
			Value: bson.D{{
				Key:   "updatedAt",
				Value: true,
			}},
		}

		updateQuery = append(updateQuery, updatedAtCommand)
	}

	return updateQuery, nil
}

// transformToBulkWriteBSONDocs converts bulkWrite entity models to mongo models.
func (m entityMongoModel[T]) transformToBulkWriteBSONDocs(ctx context.Context, bulkWrites []mongo.WriteModel) error {
	for _, bulkWrite := range bulkWrites {
		switch bulkWriteType := bulkWrite.(type) {
		case *mongo.InsertOneModel:
			doc := bulkWriteType.Document
			if doc == nil {
				continue
			}

			bsonDoc, err := m.getMongoDocFromEntityModel(ctx, doc.(T))
			if err != nil {
				return err
			}

			bulkWriteType.Document = bsonDoc
		case *mongo.UpdateOneModel:
			updateQuery, err := m.handleTimestampsForUpdateQuery(bulkWriteType.Update, "BulkWrite")
			if err != nil {
				return err
			}

			bulkWriteType.Update = updateQuery
		case *mongo.UpdateManyModel:
			updateQuery, err := m.handleTimestampsForUpdateQuery(bulkWriteType.Update, "BulkWrite")
			if err != nil {
				return err
			}

			bulkWriteType.Update = updateQuery
		}
	}
	return nil
}
