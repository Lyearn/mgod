package mgod

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mClient *mongo.Client
var dbConn *mongo.Database
var defaultTimeout = 10 * time.Second

// ConnectionConfig is the configuration options available for a MongoDB connection.
type ConnectionConfig struct {
	// Timeout is the timeout for various operations performed on the MongoDB server like Connect, Ping etc.
	Timeout time.Duration
}

// SetDefaultConnection sets the default connection to be used by the package.
func SetDefaultConnection(client *mongo.Client, dbName string) {
	mClient = client
	dbConn = mClient.Database(dbName)
}

// ConfigureDefaultConnection opens a new connection using the provided config options and sets it as a default connection to be used by the package.
func ConfigureDefaultConnection(cfg *ConnectionConfig, dbName string, opts ...*options.ClientOptions) (err error) {
	if cfg == nil {
		cfg = defaultConnectionConfig()
	}

	mClient, err = newClient(cfg, opts...)
	if err != nil {
		return err
	}

	ctx, cancel := newCtx(cfg.Timeout)
	defer cancel()

	// Ping the MongoDB server to check if connection is established.
	err = mClient.Ping(ctx, nil)
	if err != nil {
		return err
	}

	dbConn = mClient.Database(dbName)

	return nil
}

// newClient creates a new MongoDB client.
func newClient(cfg *ConnectionConfig, opts ...*options.ClientOptions) (*mongo.Client, error) {
	ctx, cancel := newCtx(cfg.Timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// newCtx returns a context with timeout (from the configured connection config).
func newCtx(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// defaultConnectionConfig returns the default connection config.
func defaultConnectionConfig() *ConnectionConfig {
	return &ConnectionConfig{
		Timeout: defaultTimeout,
	}
}
