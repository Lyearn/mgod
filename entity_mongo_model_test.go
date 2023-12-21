package mgod_test

import (
	"context"
	"testing"
	"time"

	"github.com/Lyearn/mgod"
	"github.com/Lyearn/mgod/schema/schemaopt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EntityMongoModelSuite struct {
	suite.Suite
	*require.Assertions
}

type testEntity struct {
	ID   string `bson:"_id" mgoType:"id"`
	Name string
	Age  *int `bson:",omitempty" mgoDefault:"18"`
}

func TestEntityMongoModelSuite(t *testing.T) {
	s := new(EntityMongoModelSuite)
	suite.Run(t, s)
}

func (s *EntityMongoModelSuite) SetupSuite() {
	s.setupConnection()
	s.setupData()
}

func (s *EntityMongoModelSuite) SetupTest() {
	s.Assertions = require.New(s.T())
}

func (s *EntityMongoModelSuite) TearDownSuite() {
	entityMongoModel := s.getModel()
	_, err := entityMongoModel.DeleteMany(context.Background(), bson.D{})
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *EntityMongoModelSuite) setupConnection() {
	cfg := &mgod.ConnectionConfig{Timeout: 5 * time.Second}
	uri := "mongodb://localhost:27017/?replicaSet=replset&authSource=admin"

	err := mgod.ConfigureDefaultClient(cfg, options.Client().ApplyURI(uri))
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *EntityMongoModelSuite) setupData() {
	firstID := primitive.NewObjectID()
	secondID := primitive.NewObjectID()

	age1 := 30
	age2 := 40

	entities := []testEntity{
		{
			ID:   firstID.Hex(),
			Name: "Default User 1",
			Age:  &age1,
		},
		{
			ID:   secondID.Hex(),
			Name: "Default User 2",
			Age:  &age2,
		},
	}

	entityMongoModel := s.getModel()
	_, err := entityMongoModel.InsertMany(context.Background(), entities)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *EntityMongoModelSuite) getModel() mgod.EntityMongoModel[testEntity] {
	dbName := "mgoddb"
	collection := "entityMongoModel"
	schemaOpts := schemaopt.SchemaOptions{Timestamps: true}

	opts := mgod.NewEntityMongoModelOptions(dbName, collection, &schemaOpts)
	model, err := mgod.NewEntityMongoModel(testEntity{}, *opts)
	if err != nil {
		s.T().Fatal(err)
	}

	return model
}

func (s *EntityMongoModelSuite) TestFind() {
	entityMongoModel := s.getModel()
	entities, err := entityMongoModel.Find(context.Background(), bson.M{
		"age": bson.M{
			"$gt": 20,
		},
		"name": bson.M{
			"$regex": "Default",
		},
	})

	s.NoError(err)
	s.Equal(2, len(entities))
}

func (s *EntityMongoModelSuite) TestFindOne() {
	entityMongoModel := s.getModel()
	entity, err := entityMongoModel.FindOne(context.Background(), bson.M{
		"age": bson.M{
			"$gt": 30,
		},
		"name": bson.M{
			"$regex": "Default",
		},
	})

	s.NoError(err)
	s.Equal("Default User 2", entity.Name)
}

func (s *EntityMongoModelSuite) TestInsertOne() {
	id := primitive.NewObjectID()
	age := 18

	entity := testEntity{
		ID:   id.Hex(),
		Name: "test",
		Age:  &age,
	}

	entityMongoModel := s.getModel()
	doc, err := entityMongoModel.InsertOne(context.Background(), entity)

	s.Nil(err)
	s.Equal(entity.ID, doc.ID)
	s.Equal(18, *doc.Age)

	// Test duplicate key error
	docID, err := entityMongoModel.InsertOne(context.Background(), entity)
	s.Empty(docID)
	s.NotNil(err)
	s.True(mongo.IsDuplicateKeyError(err))
}
