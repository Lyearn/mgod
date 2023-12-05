package mgod

import (
	"github.com/Lyearn/mgod/schema/schemaopt"
	"go.mongodb.org/mongo-driver/mongo"
)

// EntityMongoOptions is the options to be configured/provided when creating a new [EntityMongoModel].
type EntityMongoOptions struct {
	dbConnection  *mongo.Database
	schemaOptions schemaopt.SchemaOptions
}

// NewEntityMongoOptions creates a new instance of [EntityMongoOptions].
func NewEntityMongoOptions(dbConn *mongo.Database, schemaOpts schemaopt.SchemaOptions) *EntityMongoOptions {
	return &EntityMongoOptions{
		dbConnection:  dbConn,
		schemaOptions: schemaOpts,
	}
}
