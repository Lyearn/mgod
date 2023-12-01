package schema

import (
	"fmt"
	"sync"
)

// EntityModelSchemaCache is the cache implementation than can hold [EntityModelSchema].
// It has the following to main use cases -
// 1. Avoid re-computing the schema for the same entity model.
// 2. Fetch the relevant schema based on the discriminator key in case of union type models to validate the bson doc fetched from
// MongoDB against a concrete type.
type EntityModelSchemaCache interface {
	GetSchema(schemaName string) (*EntityModelSchema, error)
	SetSchema(schemaName string, schema *EntityModelSchema)
}

type entityModelSchemaCache struct {
	cache map[string]*EntityModelSchema
	mux   sync.RWMutex
}

func newEntityModelSchemaCache() EntityModelSchemaCache {
	return &entityModelSchemaCache{
		cache: map[string]*EntityModelSchema{},
	}
}

func (c *entityModelSchemaCache) GetSchema(schemaName string) (*EntityModelSchema, error) {
	if schema, ok := c.cache[schemaName]; ok {
		return schema, nil
	}

	return nil, fmt.Errorf("%s schema not found in cache", schemaName)
}

func (c *entityModelSchemaCache) SetSchema(schemaName string, schema *EntityModelSchema) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.cache[schemaName] = schema
}

// EntityModelSchemaCacheInstance is the singleton instance of [EntityModelSchemaCache].
var EntityModelSchemaCacheInstance = newEntityModelSchemaCache()
