package mgod

import (
	"github.com/Lyearn/mgod/schema/schemaopt"
)

type entityMongoModelOptions struct {
	connOpts   connectionOptions
	schemaOpts *schemaopt.SchemaOptions
}

type connectionOptions struct {
	db   string
	coll string
}

// NewEntityMongoModelOptions creates a new entityMongoModelOptions instance.
// Its instance is used to provide necessary configuration options to the NewEntityMongoModel function.
//
// dbName is the name of the database in which the entity is stored.
// collection is the name of the mongo collection in which the entity is stored.
// schemaOpts is the schema level options for the entity.
func NewEntityMongoModelOptions(dbName string, collection string, schemaOpts *schemaopt.SchemaOptions) *entityMongoModelOptions {
	return &entityMongoModelOptions{
		connOpts: connectionOptions{
			db:   dbName,
			coll: collection,
		},
		schemaOpts: schemaOpts,
	}
}
