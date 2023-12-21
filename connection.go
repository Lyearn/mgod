package mgod

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mClient *mongo.Client
var defaultTimeout = 10 * time.Second

// ConnectionConfig is the configuration options available for a MongoDB connection.
type ConnectionConfig struct {
	// Timeout is the timeout for various operations performed on the MongoDB server like Connect, Ping etc.
	Timeout time.Duration
}

// SetDefaultClient sets the default MongoDB client to be used by the package.
func SetDefaultClient(client *mongo.Client) {
	mClient = client
}

// ConfigureDefaultClient opens a new connection using the provided config options and sets the default MongoDB client to be used by the package.
func ConfigureDefaultClient(cfg *ConnectionConfig, opts ...*options.ClientOptions) (err error) {
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
