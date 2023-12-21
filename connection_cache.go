package mgod

import (
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

// dbConnCache is a cache of MongoDB database connections.
var dbConnCache *connectionCache

func init() {
	dbConnCache = newConnectionCache()
}

// connectionCache is a thread safe construct to cache MongoDB database connections.
type connectionCache struct {
	cache map[string]*mongo.Database
	mux   sync.RWMutex
}

func newConnectionCache() *connectionCache {
	return &connectionCache{
		cache: map[string]*mongo.Database{},
	}
}

func (c *connectionCache) Get(dbName string) *mongo.Database {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.cache[dbName]
}

func (c *connectionCache) Set(dbName string, db *mongo.Database) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.cache[dbName] = db
}

// getDBConn returns a MongoDB database connection from the cache.
// If the connection is not present in the cache, it creates a new connection and adds it to the cache (Write-through policy).
func getDBConn(dbName string) *mongo.Database {
	dbConn := dbConnCache.Get(dbName)

	// Initialize the cache entry if it is not present.
	if dbConn == nil {
		dbConn = mClient.Database(dbName)
		dbConnCache.Set(dbName, dbConn)
	}

	return dbConn
}
