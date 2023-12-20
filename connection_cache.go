package mgod

import "go.mongodb.org/mongo-driver/mongo"

// dbConnCache is a cache of MongoDB database connections.
var dbConnCache map[string]*mongo.Database

func init() {
	dbConnCache = make(map[string]*mongo.Database)
}

// getDBConn returns a MongoDB database connection from the cache.
// If the connection is not present in the cache, it creates a new connection and adds it to the cache (Write-through policy).
func getDBConn(dbName string) *mongo.Database {
	// Initialize the cache entry if it is not present.
	if dbConnCache[dbName] == nil {
		dbConnCache[dbName] = mClient.Database(dbName)
	}

	return dbConnCache[dbName]
}
