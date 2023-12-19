package mgod_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Lyearn/mgod"
	"github.com/Lyearn/mgod/schema/schemaopt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TransactionSuite struct {
	suite.Suite
	*require.Assertions

	userModel mgod.EntityMongoModel[TransactionTestUser]
}

type TransactionTestUser struct {
	Name    string
	EmailID string `bson:"emailId"`
}

func TestTransactionSuite(t *testing.T) {
	s := new(TransactionSuite)
	suite.Run(t, s)
}

func (s *TransactionSuite) SetupTest() {
	s.Assertions = require.New(s.T())
	s.SetupConnectionAndModel()
}

func (s *TransactionSuite) SetupConnectionAndModel() {
	cfg := &mgod.ConnectionConfig{Timeout: 5 * time.Second}
	dbName := "mgod_test"
	opts := options.Client().ApplyURI("mongodb://localhost:27017/?replicaSet=mgod_rs&authSource=admin")

	err := mgod.ConfigureDefaultConnection(cfg, dbName, opts)
	if err != nil {
		s.T().Fatal(err)
	}

	schemaOpts := schemaopt.SchemaOptions{
		Collection: "users",
		Timestamps: true,
	}
	userModel, err := mgod.NewEntityMongoModel(TransactionTestUser{}, schemaOpts)
	if err != nil {
		s.T().Fatal(err)
	}

	s.userModel = userModel
}

func (s *TransactionSuite) TestWithTransaction() {
	userDoc := TransactionTestUser{Name: "Gopher", EmailID: "gopher@mgod.com"}

	p, err := mgod.WithTransaction(context.Background(), func(sc mongo.SessionContext) (interface{}, error) {
		user, err := s.userModel.InsertOne(sc, userDoc)
		return user, err
	})

	user, ok := p.(TransactionTestUser)

	s.True(ok)
	s.NoError(err)
	s.Equal(user.Name, userDoc.Name)
	s.Equal(user.EmailID, userDoc.EmailID)

	userCount, err := s.userModel.CountDocuments(context.Background(), bson.M{})

	s.NoError(err)
	s.Equal(userCount, int64(1))
}

func (s *TransactionSuite) TestWithTransactionAbort() {
	userDoc := TransactionTestUser{Name: "Gopher", EmailID: "gopher@mgod.com"}

	abortErr := errors.New("dummy error to abort transaction")

	_, err := mgod.WithTransaction(context.Background(), func(sc mongo.SessionContext) (interface{}, error) {
		_, err := s.userModel.InsertOne(sc, userDoc)
		if err != nil {
			return nil, err
		}

		return nil, abortErr
	})

	s.EqualError(err, abortErr.Error())

	userCount, err := s.userModel.CountDocuments(context.Background(), bson.M{})

	s.NoError(err)
	s.Equal(userCount, int64(0))
}
