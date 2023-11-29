package mgod

import (
	"context"
	"fmt"

	"github.com/Lyearn/mgod/errors"
	"github.com/Lyearn/mgod/metafield"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EntityMongoModel[T any] interface {
	GetDocToInsert(ctx context.Context, model T) (bson.D, error)
	InsertOne(ctx context.Context, model interface{}, opts ...*options.InsertOneOptions) (T, error)
	InsertMany(ctx context.Context, docs interface{}, opts ...*options.InsertManyOptions) ([]T, error)
	UpdateMany(ctx context.Context, filter, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	BulkWrite(ctx context.Context, bulkWrites []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]T, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (*T, error)
	FindOneAndUpdate(ctx context.Context, filter, update interface{}, opts ...*options.FindOneAndUpdateOptions) (T, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
	Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error)
	Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) ([]bson.D, error)
}

type entityMongoModel[T any] struct {
	modelType T
	opts      EntityMongoOptions
	coll      *mongo.Collection

	schema *EntityModelSchema

	isUnionType      bool
	discriminatorKey string
}

func NewEntityMongoModel[T any](modelType T, opts EntityMongoOptions) (EntityMongoModel[T], error) {
	dbConnection := opts.dbConnection
	if dbConnection == nil {
		return nil, errors.ErrNoDatabaseConnection
	}

	coll := dbConnection.Collection(opts.schemaOptions.Collection)

	modelName := GetSchemaNameForModel(modelType)
	schemaCacheKey := GetSchemaCacheKey(coll.Name(), modelName)

	var schema *EntityModelSchema
	var err error

	// build schema if not cached.
	if schema, err = entityModelSchemaCacheInstance.GetSchema(schemaCacheKey); err != nil {
		schema, err = BuildSchemaForModel(modelType, opts.schemaOptions)
		if err != nil {
			return nil, err
		}

		entityModelSchemaCacheInstance.SetSchema(schemaCacheKey, schema)
	}

	isUnionTypeModel := opts.schemaOptions.IsUnionType

	discriminatorKey := "__t"
	if isUnionTypeModel && opts.schemaOptions.DiscriminatorKey != nil {
		discriminatorKey = *opts.schemaOptions.DiscriminatorKey
	}

	return &entityMongoModel[T]{
		modelType:        modelType,
		opts:             opts,
		coll:             coll,
		schema:           schema,
		isUnionType:      isUnionTypeModel,
		discriminatorKey: discriminatorKey,
	}, nil
}

func (m entityMongoModel[T]) getEntityModel() T {
	return m.modelType
}

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

	if err = metafield.AddMetaFields(&bsonDoc, m.opts.schemaOptions); err != nil {
		return nil, err
	}

	err = BuildBSONDoc(ctx, &bsonDoc, m.schema, BSONDocTranslateToEnumMongo)
	if err != nil {
		return nil, err
	}

	if m.isUnionType {
		discriminatorVal := getFieldValueFromBSONRootDoc(&bsonDoc, m.discriminatorKey)
		if discriminatorVal == nil {
			discriminatorVal = GetSchemaNameForModel(m.modelType)
			bsonDoc = append(bsonDoc, primitive.E{
				Key:   m.discriminatorKey,
				Value: discriminatorVal,
			})
		}

		cacheKey := GetSchemaCacheKey(m.coll.Name(), discriminatorVal.(string))
		if _, err := entityModelSchemaCacheInstance.GetSchema(cacheKey); err != nil {
			entityModelSchemaCacheInstance.SetSchema(cacheKey, m.schema)
		}
	}

	return bsonDoc, nil
}

func (m entityMongoModel[T]) getEntityModelFromMongoDoc(ctx context.Context, bsonDoc bson.D) (T, error) {
	model := m.getEntityModel()

	if bsonDoc == nil {
		// empty bson doc
		return model, nil
	}

	entityModelSchema := m.schema

	if m.isUnionType {
		discriminatorVal := getFieldValueFromBSONRootDoc(&bsonDoc, m.discriminatorKey)
		if discriminatorVal != nil {
			cacheKey := GetSchemaCacheKey(m.coll.Name(), discriminatorVal.(string))
			if unionElemSchema, err := entityModelSchemaCacheInstance.GetSchema(cacheKey); err == nil {
				entityModelSchema = unionElemSchema
			}
		}
	}

	err := BuildBSONDoc(ctx, &bsonDoc, entityModelSchema, BSONDocTranslateToEnumEntityModel)
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

func (m entityMongoModel[T]) handleTimestampsForUpdateQuery(update interface{}, funcName string) (interface{}, error) {
	updateQuery, ok := update.(bson.D)
	if !ok {
		return nil, errors.NewBadRequestError(errors.BadRequestError{
			Underlying: "update query",
			Got:        fmt.Sprintf("%T", update),
			Expected:   "bson.D",
		})
	}

	if m.opts.schemaOptions.Timestamps {
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

// Converts bulkWrite entity models to mongo models.
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

func (m entityMongoModel[T]) GetDocToInsert(ctx context.Context, doc T) (bson.D, error) {
	bsonDoc, err := m.getMongoDocFromEntityModel(ctx, doc)
	if err != nil {
		return nil, err
	}

	return bsonDoc, nil
}

func (m entityMongoModel[T]) InsertOne(ctx context.Context, doc interface{},
	opts ...*options.InsertOneOptions,
) (T, error) {
	model := m.getEntityModel()

	var bsonDoc primitive.D
	var err error

	switch typedDoc := doc.(type) {
	case bson.D:
		bsonDoc = typedDoc
	case T:
		bsonDoc, err = m.getMongoDocFromEntityModel(ctx, typedDoc)
		if err != nil {
			return model, err
		}
	}

	// TODO: add an extra strict check to ensure that the doc to be inserted contains _id field

	_, err = m.coll.InsertOne(ctx, bsonDoc, opts...)
	if err != nil {
		return model, err
	}

	model, err = m.getEntityModelFromMongoDoc(ctx, bsonDoc)

	return model, err
}

func (m entityMongoModel[T]) InsertMany(ctx context.Context, docs interface{},
	opts ...*options.InsertManyOptions,
) ([]T, error) {
	bsonDocs := []interface{}{}

	switch typedDocs := docs.(type) {
	case []T:
		for _, doc := range typedDocs {
			bsonDoc, err := m.getMongoDocFromEntityModel(ctx, doc)
			if err != nil {
				return nil, err
			}

			bsonDocs = append(bsonDocs, bsonDoc)
		}
	case []bson.D:
		bsonDocs = lo.Map(typedDocs, func(typedDoc bson.D, _ int) interface{} {
			return typedDoc
		})
	default:
		var dummyTypedVar T
		return nil, errors.NewBadRequestError(errors.BadRequestError{
			Underlying: "insertMany docs",
			Got:        fmt.Sprintf("%T", typedDocs),
			Expected:   fmt.Sprintf("array of %T or bson.D", dummyTypedVar),
		})
	}

	_, err := m.coll.InsertMany(ctx, bsonDocs, opts...)
	if err != nil {
		return nil, err
	}

	models := []T{}

	// TODO: transform and return only those docs which are inserted successfully (use result from InsertMany)
	for _, bsonDoc := range bsonDocs {
		model, transformErr := m.getEntityModelFromMongoDoc(ctx, bsonDoc.(primitive.D))
		if transformErr != nil {
			return nil, transformErr
		} else {
			models = append(models, model)
		}
	}

	return models, err
}

func (m entityMongoModel[T]) UpdateMany(ctx context.Context, filter, update interface{},
	opts ...*options.UpdateOptions,
) (*mongo.UpdateResult, error) {
	updateQuery, err := m.handleTimestampsForUpdateQuery(update, "UpdateMany")
	if err != nil {
		return nil, err
	}

	result, err := m.coll.UpdateMany(ctx, filter, updateQuery, opts...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m entityMongoModel[T]) BulkWrite(ctx context.Context, bulkWrites []mongo.WriteModel,
	opts ...*options.BulkWriteOptions,
) (*mongo.BulkWriteResult, error) {
	err := m.transformToBulkWriteBSONDocs(ctx, bulkWrites)
	if err != nil {
		return nil, err
	}
	result, err := m.coll.BulkWrite(ctx, bulkWrites, opts...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m entityMongoModel[T]) Find(ctx context.Context, filter interface{},
	opts ...*options.FindOptions,
) ([]T, error) {
	cursor, err := m.coll.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	var docs []bson.D
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	var models []T
	for _, doc := range docs {
		model, err := m.getEntityModelFromMongoDoc(ctx, doc)
		if err != nil {
			return nil, err
		}

		models = append(models, model)
	}

	return models, nil
}

func (m entityMongoModel[T]) FindOne(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions,
) (*T, error) {
	cursor := m.coll.FindOne(ctx, filter, opts...)

	var doc bson.D

	model := m.getEntityModel()
	var err error

	if err = cursor.Decode(&doc); err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			//nolint:nilnil // this is the expected behavior
			return nil, nil
		}
		return nil, err
	}

	model, err = m.getEntityModelFromMongoDoc(ctx, doc)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (m entityMongoModel[T]) FindOneAndUpdate(ctx context.Context, filter, update interface{},
	opts ...*options.FindOneAndUpdateOptions,
) (T, error) {
	model := m.getEntityModel()
	var err error

	updateQuery, err := m.handleTimestampsForUpdateQuery(update, "FindOneAndUpdate")
	if err != nil {
		return model, err
	}

	cursor := m.coll.FindOneAndUpdate(ctx, filter, updateQuery, opts...)

	var doc bson.D

	if err = cursor.Decode(&doc); err != nil {
		return model, err
	}

	model, err = m.getEntityModelFromMongoDoc(ctx, doc)
	if err != nil {
		return model, err
	}

	return model, nil
}

func (m entityMongoModel[T]) DeleteOne(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions,
) (*mongo.DeleteResult, error) {
	return m.coll.DeleteOne(ctx, filter, opts...)
}

func (m entityMongoModel[T]) DeleteMany(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions,
) (*mongo.DeleteResult, error) {
	return m.coll.DeleteMany(ctx, filter, opts...)
}

func (m entityMongoModel[T]) CountDocuments(ctx context.Context, filter interface{},
	opts ...*options.CountOptions,
) (int64, error) {
	return m.coll.CountDocuments(ctx, filter, opts...)
}

func (m entityMongoModel[T]) Distinct(ctx context.Context, fieldName string, filter interface{},
	opts ...*options.DistinctOptions,
) ([]interface{}, error) {
	return m.coll.Distinct(ctx, fieldName, filter, opts...)
}

func (m entityMongoModel[T]) Aggregate(ctx context.Context, pipeline interface{},
	opts ...*options.AggregateOptions,
) ([]bson.D, error) {
	cursor, err := m.coll.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return nil, err
	}

	var docs []bson.D
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	return docs, nil
}
