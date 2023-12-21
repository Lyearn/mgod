// Package mgod implements ODM logic for MongoDB in Go.
package mgod

import (
	"context"
	"fmt"

	"github.com/Lyearn/mgod/errors"
	"github.com/Lyearn/mgod/schema"
	"github.com/Lyearn/mgod/schema/schemaopt"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EntityMongoModel is a generic interface of available wrapper functions on MongoDB collection.
type EntityMongoModel[T any] interface {
	// GetDocToInsert returns the bson.D doc to be inserted in the collection for the provided struct object.
	// This function is mainly used while creating a doc to be inserted for Union Type models because the underlying type of a union
	// type model is interface{}, so it's not possible to identify the underlying concrete type to validate and insert the doc.
	GetDocToInsert(ctx context.Context, model T) (bson.D, error)

	// InsertOne inserts a single document in the collection.
	// Model is kept as interface{} to support Union Type models i.e. accept both bson.D (generated using GetDocToInsert()) and struct object.
	InsertOne(ctx context.Context, model interface{}, opts ...*options.InsertOneOptions) (T, error)

	// InsertMany inserts multiple documents in the collection.
	// Docs is kept as interface{} to support Union Type models i.e. accept both []bson.D (generated using GetDocToInsert()) and []struct objects.
	InsertMany(ctx context.Context, docs interface{}, opts ...*options.InsertManyOptions) ([]T, error)

	// UpdateMany updates multiple filtered documents in the collection based on the provided update query.
	UpdateMany(ctx context.Context, filter, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)

	// BulkWrite performs multiple write operations on the collection at once.
	// Currently, only InsertOne, UpdateOne, and UpdateMany operations are supported.
	BulkWrite(ctx context.Context, bulkWrites []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error)

	// Find returns all documents in the collection matching the provided filter.
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]T, error)

	// FindOne returns a single document from the collection matching the provided filter.
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (*T, error)

	// FindOneAndUpdate returns a single document from the collection based on the provided filter and updates it.
	FindOneAndUpdate(ctx context.Context, filter, update interface{}, opts ...*options.FindOneAndUpdateOptions) (T, error)

	// DeleteOne deletes a single document in the collection based on the provided filter.
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)

	// DeleteMany deletes multiple documents in the collection based on the provided filter.
	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)

	// CountDocuments returns the number of documents in the collection for the provided filter.
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)

	// Distinct returns the distinct values for the provided field name in the collection for the provided filter.
	Distinct(ctx context.Context, fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error)

	// Aggregate performs an aggregation operation on the collection and returns the results.
	Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) ([]bson.D, error)
}

type entityMongoModel[T any] struct {
	modelType  T
	schemaOpts schemaopt.SchemaOptions
	coll       *mongo.Collection

	schema *schema.EntityModelSchema

	isUnionType      bool
	discriminatorKey string
}

// NewEntityMongoModel returns a new instance of EntityMongoModel for the provided model type and options.
func NewEntityMongoModel[T any](modelType T, opts entityMongoModelOptions) (EntityMongoModel[T], error) {
	dbConn := getDBConn(opts.connOpts.db)
	if dbConn == nil {
		return nil, errors.ErrNoDatabaseConnection
	}

	coll := dbConn.Collection(opts.connOpts.coll)

	modelName := schema.GetSchemaNameForModel(modelType)
	schemaCacheKey := GetSchemaCacheKey(coll.Name(), modelName)

	var entityModelSchema *schema.EntityModelSchema
	var err error

	schemaOpts := schemaopt.SchemaOptions{}
	if opts.schemaOpts != nil {
		schemaOpts = *opts.schemaOpts
	}

	// build schema if not cached.
	if entityModelSchema, err = schema.EntityModelSchemaCacheInstance.GetSchema(schemaCacheKey); err != nil {
		entityModelSchema, err = schema.BuildSchemaForModel(modelType, schemaOpts)
		if err != nil {
			return nil, err
		}

		schema.EntityModelSchemaCacheInstance.SetSchema(schemaCacheKey, entityModelSchema)
	}

	isUnionTypeModel := schemaOpts.IsUnionType

	discriminatorKey := "__t"
	if isUnionTypeModel && schemaOpts.DiscriminatorKey != nil {
		discriminatorKey = *schemaOpts.DiscriminatorKey
	}

	return &entityMongoModel[T]{
		modelType:        modelType,
		schemaOpts:       schemaOpts,
		coll:             coll,
		schema:           entityModelSchema,
		isUnionType:      isUnionTypeModel,
		discriminatorKey: discriminatorKey,
	}, nil
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
