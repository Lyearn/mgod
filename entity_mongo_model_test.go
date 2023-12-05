package mgod_test

import (
	"context"
	"testing"
	"time"

	"github.com/Lyearn/mgod"
	"github.com/Lyearn/mgod/dateformatter"
	"github.com/Lyearn/mgod/schema/schemaopt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

type EntityMongoModelSuite struct {
	suite.Suite
	*require.Assertions

	dbName   string
	collName string

	mt     *mtest.T
	mtOpts *mtest.Options
}

func TestEntityMongoModelSuite(t *testing.T) {
	s := new(EntityMongoModelSuite)
	s.dbName = "foo"
	s.collName = "bar"

	suite.Run(t, s)
}

func (s *EntityMongoModelSuite) SetupTest() {
	s.Assertions = require.New(s.T())

	mtOpts := mtest.NewOptions()
	mtOpts = mtOpts.ClientType(mtest.Mock)
	mtOpts = mtOpts.DatabaseName(s.dbName)
	mtOpts = mtOpts.CollectionName(s.collName)

	mt := mtest.New(s.T(), mtOpts)

	s.mt = mt
	s.mtOpts = mtOpts
}

func (s *EntityMongoModelSuite) ns() string {
	return s.dbName + "." + s.collName
}

type TestEntity struct {
	ID       string `bson:"_id" mgoType:"id"`
	Name     string
	JoinedOn string `mgoType:"date"`
	Age      *int   `bson:",omitempty" mgoDefault:"18"`
}

func (s *EntityMongoModelSuite) TestFind() {
	defer s.mt.Close()

	s.mt.RunOpts("find", s.mtOpts, func(mt *mtest.T) {
		currentTime := time.Now()
		currentTimeStr, _ := dateformatter.New(currentTime).GetISOString()

		firstID := primitive.NewObjectID()
		secondID := primitive.NewObjectID()

		//nolint:govet // this is a mock entity.
		firstEntity := TestEntity{
			ID:       firstID.Hex(),
			Name:     "test1",
			JoinedOn: currentTimeStr,
		}
		//nolint:govet // this is a mock entity.
		secondEntity := TestEntity{
			ID:       secondID.Hex(),
			Name:     "test2",
			JoinedOn: currentTimeStr,
		}

		first := mtest.CreateCursorResponse(1, s.ns(), mtest.FirstBatch, bson.D{
			{Key: "_id", Value: firstID},
			{Key: "name", Value: firstEntity.Name},
			{Key: "joinedon", Value: primitive.NewDateTimeFromTime(currentTime)},
		})
		second := mtest.CreateCursorResponse(1, s.ns(), mtest.NextBatch, bson.D{
			{Key: "_id", Value: secondID},
			{Key: "name", Value: secondEntity.Name},
			{Key: "joinedon", Value: primitive.NewDateTimeFromTime(currentTime)},
		})
		killCursors := mtest.CreateCursorResponse(0, s.ns(), mtest.NextBatch)

		mt.AddMockResponses(first, second, killCursors)

		opts := mgod.NewEntityMongoOptions(mt.DB, schemaopt.SchemaOptions{Collection: s.collName})

		entityMongoModel, err := mgod.NewEntityMongoModel(TestEntity{}, *opts)
		s.Nil(err)

		testEntities, err := entityMongoModel.Find(context.Background(), bson.D{
			{Key: "name", Value: firstEntity.Name},
		})

		s.Nil(err)
		s.Equal(2, len(testEntities))
	})
}

func (s *EntityMongoModelSuite) TestFindOne() {
	defer s.mt.Close()

	s.mt.RunOpts("find one", s.mtOpts, func(mt *mtest.T) {
		currentTime := time.Now()
		currentTimeStr, _ := dateformatter.New(currentTime).GetISOString()

		id := primitive.NewObjectID()

		//nolint:govet // this is a mock entity.
		entity := TestEntity{
			ID:       id.Hex(),
			Name:     "test",
			JoinedOn: currentTimeStr,
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, s.ns(), mtest.FirstBatch, bson.D{
			{Key: "_id", Value: id},
			{Key: "name", Value: entity.Name},
			{Key: "joinedon", Value: primitive.NewDateTimeFromTime(currentTime)},
		}))

		opts := mgod.NewEntityMongoOptions(mt.DB, schemaopt.SchemaOptions{Collection: s.collName})

		entityMongoModel, err := mgod.NewEntityMongoModel(TestEntity{}, *opts)
		s.Nil(err)

		testEntity, err := entityMongoModel.FindOne(context.Background(), bson.D{
			{Key: "id", Value: entity.ID},
		})

		s.Nil(err)
		s.Equal(entity.ID, testEntity.ID)
	})
}

func (s *EntityMongoModelSuite) TestInsertOne() {
	defer s.mt.Close()

	id := primitive.NewObjectID()
	currentTime := time.Now()
	currentTimeStr, _ := dateformatter.New(currentTime).GetISOString()

	s.mt.RunOpts("insert one", s.mtOpts, func(mt *mtest.T) {
		entity := TestEntity{
			ID:       id.Hex(),
			Name:     "test",
			JoinedOn: currentTimeStr,
		}

		mt.AddMockResponses(mtest.CreateCursorResponse(1, s.ns(), mtest.FirstBatch, bson.D{
			{Key: "_id", Value: id},
			{Key: "name", Value: entity.Name},
			{Key: "joinedon", Value: primitive.NewDateTimeFromTime(currentTime)},
			{Key: "age", Value: 18},
		}))

		opts := mgod.NewEntityMongoOptions(mt.DB, schemaopt.SchemaOptions{Collection: s.collName})

		entityMongoModel, err := mgod.NewEntityMongoModel(TestEntity{}, *opts)
		s.Nil(err)

		doc, err := entityMongoModel.InsertOne(context.Background(), entity)

		s.Nil(err)
		s.Equal(entity.ID, doc.ID)
		s.Equal(18, *doc.Age)
	})

	s.mt.RunOpts("insert one with error", s.mtOpts, func(mt *mtest.T) {
		entity := TestEntity{
			ID:       id.Hex(),
			Name:     "test",
			JoinedOn: currentTimeStr,
		}

		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))

		opts := mgod.NewEntityMongoOptions(mt.DB, schemaopt.SchemaOptions{Collection: s.collName})

		entityMongoModel, err := mgod.NewEntityMongoModel(TestEntity{}, *opts)
		s.Nil(err)

		docID, err := entityMongoModel.InsertOne(context.Background(), entity)

		s.Empty(docID)
		s.NotNil(err)
		s.True(mongo.IsDuplicateKeyError(err))
	})
}
