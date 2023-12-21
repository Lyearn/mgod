package mgod_test

import (
	"context"
	"errors"
	"fmt"
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
}

type transactionTestUser struct {
	Name    string
	EmailID string `bson:"emailId"`
}

func TestTransactionSuite(t *testing.T) {
	s := new(TransactionSuite)
	suite.Run(t, s)
}

func (s *TransactionSuite) SetupSuite() {
	s.setupConnection()
}

func (s *TransactionSuite) SetupTest() {
	s.Assertions = require.New(s.T())
}

func (s *TransactionSuite) setupConnection() {
	cfg := &mgod.ConnectionConfig{Timeout: 5 * time.Second}
	// Can use the `mlaunch` tool to start a local replica set using command `mlaunch --repl`.
	uri := "mongodb://localhost:27017/?replicaSet=replset&authSource=admin"

	err := mgod.ConfigureDefaultClient(cfg, options.Client().ApplyURI(uri))
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *TransactionSuite) getModelForDB(dbName string) mgod.EntityMongoModel[transactionTestUser] {
	schemaOpts := schemaopt.SchemaOptions{Timestamps: true}
	opts := mgod.NewEntityMongoModelOptions(dbName, "users", &schemaOpts)
	userModel, err := mgod.NewEntityMongoModel(transactionTestUser{}, *opts)
	if err != nil {
		s.T().Fatal(err)
	}

	return userModel
}

func (s *TransactionSuite) TestWithTransaction() {
	userModel := s.getModelForDB("mgod1")
	userDoc := transactionTestUser{Name: "Gopher", EmailID: "gopher@mgod.com"}

	p, err := mgod.WithTransaction(context.Background(), func(sc mongo.SessionContext) (interface{}, error) {
		_, err := userModel.InsertOne(sc, userDoc)
		if err != nil {
			return nil, err
		}

		userCount, err := userModel.CountDocuments(sc, bson.M{})
		if err != nil {
			return nil, err
		}

		_, err = userModel.DeleteOne(sc, bson.M{})
		if err != nil {
			return nil, err
		}

		return userCount, nil
	})

	userCount, ok := p.(int64)

	s.NoError(err)
	s.True(ok)
	s.Equal(userCount, int64(1))
}

func (s *TransactionSuite) TestWithTransactionAbort() {
	userModel := s.getModelForDB("mgod1")
	userDoc := transactionTestUser{Name: "Gopher", EmailID: "gopher@mgod.com"}

	abortErr := errors.New("dummy error to abort transaction")

	_, err := mgod.WithTransaction(context.Background(), func(sc mongo.SessionContext) (interface{}, error) {
		_, err := userModel.InsertOne(sc, userDoc)
		if err != nil {
			return nil, err
		}

		return nil, abortErr
	})

	s.EqualError(err, abortErr.Error())

	userCount, err := userModel.CountDocuments(context.Background(), bson.M{})

	s.NoError(err)
	s.Equal(userCount, int64(0))
}

func (s *TransactionSuite) TestWithTransactionForMultiTenancy() {
	userModelTenant1 := s.getModelForDB("mgod1")
	userModelTenant2 := s.getModelForDB("mgod2")

	userDoc := transactionTestUser{Name: "Gopher", EmailID: "gopher@mgod.com"}

	p, err := mgod.WithTransaction(context.Background(), func(sc mongo.SessionContext) (interface{}, error) {
		_, err := userModelTenant1.InsertOne(sc, userDoc)
		if err != nil {
			return nil, err
		}
		_, err = userModelTenant2.InsertOne(sc, userDoc)
		if err != nil {
			return nil, err
		}

		userCount1, err := userModelTenant1.CountDocuments(sc, bson.M{})
		if err != nil {
			return nil, err
		}
		userCount2, err := userModelTenant2.CountDocuments(sc, bson.M{})
		if err != nil {
			return nil, err
		}

		_, err = userModelTenant1.DeleteOne(sc, bson.M{})
		if err != nil {
			return nil, err
		}
		_, err = userModelTenant2.DeleteOne(sc, bson.M{})
		if err != nil {
			return nil, err
		}

		return fmt.Sprintf("%d%d", userCount1, userCount2), nil
	})

	userCountStr, ok := p.(string)

	s.NoError(err)
	s.True(ok)
	s.Equal(userCountStr, "11")
}
