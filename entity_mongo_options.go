package mgod

import (
	"github.com/Lyearn/mgod/schema/schemaopt"
	"go.mongodb.org/mongo-driver/mongo"
)

// EntityMongoOptions is the options to be configured/provided when creating a new [EntityMongoModel].
type EntityMongoOptions struct {
	schemaOptions schemaopt.SchemaOptions
	dbConnection  *mongo.Database
}

// NewEntityMongoOptions creates a new instance of [EntityMongoOptions] with the provided database connection.
func NewEntityMongoOptions(db *mongo.Database) *EntityMongoOptions {
	return &EntityMongoOptions{
		dbConnection: db,
	}
}

// SetSchemaOptions sets the updated [schemaopt.SchemaOptions].
func (o *EntityMongoOptions) SetSchemaOptions(schemaOptions schemaopt.SchemaOptions) *EntityMongoOptions {
	o.schemaOptions = schemaOptions
	return o
}
