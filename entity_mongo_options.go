package mgod

import (
	"github.com/Lyearn/mgod/schema/schemaopt"
	"go.mongodb.org/mongo-driver/mongo"
)

type EntityMongoOptions struct {
	schemaOptions schemaopt.SchemaOptions
	dbConnection  *mongo.Database
}

func NewEntityMongoOptions(db *mongo.Database) *EntityMongoOptions {
	return &EntityMongoOptions{
		dbConnection: db,
	}
}

func (o *EntityMongoOptions) SetSchemaOptions(schemaOptions schemaopt.SchemaOptions) *EntityMongoOptions {
	o.schemaOptions = schemaOptions
	return o
}
