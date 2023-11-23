package mgod

import (
	"fmt"
	"sync"
)

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

var entityModelSchemaCacheInstance = newEntityModelSchemaCache()
