package mgod

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbConn *mongo.Database
var defaultTimeout = 10 * time.Second

// ConnectionConfig is the configuration for a MongoDB connection.
type ConnectionConfig struct {
	// Timeout is the timeout for the connection.
	Timeout time.Duration
}

// SetDefaultConnection sets the default connection to be used by the package.
func SetDefaultConnection(conn *mongo.Database) {
	dbConn = conn
}

// ConfigureDefaultConnection opens a new connection using the provided config options and sets it as a default connection to be used by the package.
func ConfigureDefaultConnection(config *ConnectionConfig, dbName string, opts ...*options.ClientOptions) error {
	if config == nil {
		config = defaultConnectionConfig()
	}

	client, err := newClient(config, opts...)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	// Ping the MongoDB server to check if connection is established.
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	dbConn = client.Database(dbName)

	return nil
}

// defaultConnectionConfig returns the default connection config.
func defaultConnectionConfig() *ConnectionConfig {
	return &ConnectionConfig{
		Timeout: defaultTimeout,
	}
}

// newClient creates a new MongoDB client.
func newClient(config *ConnectionConfig, opts ...*options.ClientOptions) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return client, nil
}
