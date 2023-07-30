package mongomodel

import (
	"github.com/Lyearn/backend-universe/packages/store/acl/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type EntityMongoOptions struct {
	schemaOptions model.SchemaOptions
	dbConnection  *mongo.Database
}

func NewEntityMongoOptions() *EntityMongoOptions {
	return &EntityMongoOptions{}
}

func (o *EntityMongoOptions) SetSchemaOptions(schemaOptions model.SchemaOptions) *EntityMongoOptions {
	o.schemaOptions = schemaOptions
	return o
}

func (o *EntityMongoOptions) SetConnection(db *mongo.Database) *EntityMongoOptions {
	o.dbConnection = db
	return o
}
