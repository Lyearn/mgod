package mongomodel

import (
	"context"
	"net/http"

	"github.com/Lyearn/backend-universe/packages/observability/errorhandler"
	"github.com/Lyearn/backend-universe/packages/store/mongomodel/metafield"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EntityMongoModel[T any] interface {
	InsertOne(ctx context.Context, model T, opts ...*options.InsertOneOptions) (T, error)
	InsertMany(ctx context.Context, docs []T, opts ...*options.InsertManyOptions) ([]T, error)
	UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	BulkWrite(ctx context.Context, bulkWrites []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]T, error)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (*T, error)
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) (T, error)
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
}

func NewEntityMongoModel[T any](modelType T, opts EntityMongoOptions) (EntityMongoModel[T], error) {
	dbConnection := opts.dbConnection

	// Use global client if no connection is provided
	if dbConnection == nil {
		if conn, err := GetGlobalConnection(); err == nil {
			dbConnection = conn
		}
	}

	var model T
	schema, err := BuildSchemaForModel(model, opts.schemaOptions)
	if err != nil {
		return nil, err
	}

	return &entityMongoModel[T]{
		modelType: modelType,
		opts:      opts,
		coll:      dbConnection.Collection(opts.schemaOptions.Collection),
		schema:    schema,
	}, nil
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

	return bsonDoc, nil
}

func (m entityMongoModel[T]) getEntityModelFromMongoDoc(ctx context.Context, bsonDoc bson.D) (T, error) {
	var model T

	if bsonDoc == nil {
		// empty bson doc
		return model, nil
	}

	err := BuildBSONDoc(ctx, &bsonDoc, m.schema, BSONDocTranslateToEnumEntityModel)
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
		errParams := map[string]interface{}{
			"updateQuery": update,
		}

		return nil, errorhandler.NewLyearnError(errorhandler.LyearnErrorProps{
			Message:    "Update query is not of required type bson.D",
			Where:      "entityMongoModel." + funcName,
			HttpStatus: http.StatusInternalServerError,
			Params:     &errParams,
			ReportBug:  true,
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

func (m entityMongoModel[T]) InsertOne(ctx context.Context, doc T,
	opts ...*options.InsertOneOptions,
) (T, error) {
	var model T

	bsonDoc, err := m.getMongoDocFromEntityModel(ctx, doc)
	if err != nil {
		return model, err
	}

	// metafield.AddMetaFields(&bsonDoc, m.opts.schemaOptions)

	// TODO: add an extra strict check to ensure that the doc to be inserted contains _id field

	result, err := m.coll.InsertOne(ctx, bsonDoc, opts...)
	if err != nil {
		return model, err
	}

	if _, ok := result.InsertedID.(primitive.ObjectID); !ok {
		errorParams := map[string]interface{}{
			"docId":    result.InsertedID,
			"inputDoc": doc,
		}

		return model, errorhandler.NewLyearnError(errorhandler.LyearnErrorProps{
			Message:    "Invalid doc id returned while inserting document in mongo",
			Where:      "entityMongoModel.Create",
			HttpStatus: http.StatusInternalServerError,
			Params:     &errorParams,
			ReportBug:  true,
		})
	}

	// TODO: add an extra strict check to ensure that the docID returned by mongo is the same as the one passed in bsonDoc

	model, err = m.getEntityModelFromMongoDoc(ctx, bsonDoc)

	return model, err
}

func (m entityMongoModel[T]) InsertMany(ctx context.Context, docs []T,
	opts ...*options.InsertManyOptions,
) ([]T, error) {
	bsonDocs := make([]interface{}, len(docs))

	for i, doc := range docs {
		bsonDoc, err := m.getMongoDocFromEntityModel(ctx, doc)
		if err != nil {
			return nil, err
		}

		bsonDocs[i] = bsonDoc
	}

	result, err := m.coll.InsertMany(ctx, bsonDocs, opts...)
	if err != nil {
		return nil, err
	}

	if result == nil || len(result.InsertedIDs) != len(docs) {
		errParams := map[string]interface{}{
			"mongoResult": result,
			"inputDocs":   docs,
		}

		return nil, errorhandler.NewLyearnError(errorhandler.LyearnErrorProps{
			Message:    "Failed to create some documents in mongo",
			Where:      "entityMongoModel.InsertMany",
			HttpStatus: http.StatusInternalServerError,
			Params:     &errParams,
			ReportBug:  true,
		})
	}

	models := []T{}

	for i, insertedID := range result.InsertedIDs {
		if _, ok := insertedID.(primitive.ObjectID); !ok {
			errorParams := map[string]interface{}{
				"insertedDocId": insertedID,
				"inputDoc":      docs[i],
			}

			return nil, errorhandler.NewLyearnError(errorhandler.LyearnErrorProps{
				Message:    "Invalid doc id returned while inserting document in mongo",
				Where:      "entityMongoModel.InsertMany",
				HttpStatus: http.StatusInternalServerError,
				Params:     &errorParams,
				ReportBug:  true,
			})
		}

		model, transformErr := m.getEntityModelFromMongoDoc(ctx, bsonDocs[i].(primitive.D))
		if transformErr != nil {
			return nil, transformErr
		} else {
			models = append(models, model)
		}
	}

	return models, err
}

func (m entityMongoModel[T]) UpdateMany(ctx context.Context, filter interface{}, update interface{},
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

	var model T
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

func (m entityMongoModel[T]) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.FindOneAndUpdateOptions,
) (T, error) {
	var model T
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
