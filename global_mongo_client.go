package mgod

/*
 * This file is exact replica of mongoclient/mongodb_client.go file.
 * Reason to have this file is to avoid cyclic dependency between packages.
 */

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/Lyearn/backend-universe/packages/common/configs"
	"github.com/Lyearn/backend-universe/packages/common/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type ConnectionConfigs struct {
	URI          string
	IsTLSEnabled bool
	RootCA       *string
	TLSCert      *string
	TLSPrivKey   *string
}

type ConnectionBuilder struct {
	mu sync.Mutex

	configs ConnectionConfigs
	client  *mongo.Client
}

func NewConnection(configs ConnectionConfigs) *ConnectionBuilder {
	return &ConnectionBuilder{
		configs: configs,
	}
}

func DefaultConnectionConfigs() ConnectionConfigs {
	tlsConfigs := configs.BaseConfigInstance.TLS()
	mongodbConfigs := configs.BaseConfigInstance.Database().MongoDB()

	uri, err := mongodbConfigs.URL()
	if err != nil {
		panic(err)
	}

	rootCAPath := tlsConfigs.RootCAPAth()
	tslCertPath := tlsConfigs.CertificatePath()
	tslPrivKeyPath := tlsConfigs.PrivateKeyPath()

	return ConnectionConfigs{
		URI:          uri,
		IsTLSEnabled: mongodbConfigs.TLSEnabled(),
		RootCA:       &rootCAPath,
		TLSCert:      &tslCertPath,
		TLSPrivKey:   &tslPrivKeyPath,
	}
}

func (b *ConnectionBuilder) getTLSMongoURI() string {
	if !b.configs.IsTLSEnabled {
		logger.Info(context.TODO(), "tls is not enabled, returning base mongo uri")
		return b.configs.URI
	}

	if b.configs.RootCA == nil || b.configs.TLSCert == nil || b.configs.TLSPrivKey == nil {
		logger.Error(context.TODO(), "required tls configs are not provided, returning base mongo uri")
		return b.configs.URI
	}

	tlsArg := fmt.Sprintf("%v=%v", "tls", "true")
	tlsRootCaFileArg := fmt.Sprintf("%v=%v", "tlsCAFile", *b.configs.RootCA)
	tlsCertificateFileArg := fmt.Sprintf("%v=%v", "tlsCertificateFile", *b.configs.TLSCert)
	tlsPrivateKeyFileArg := fmt.Sprintf("%v=%v", "tlsPrivateKeyFile", *b.configs.TLSPrivKey)

	mongoURIArgs := []string{b.configs.URI, tlsArg, tlsPrivateKeyFileArg, tlsCertificateFileArg, tlsRootCaFileArg}
	uri := strings.Join(mongoURIArgs, "&")

	return uri
}

func (b *ConnectionBuilder) Connect() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.client != nil {
		// already connected
		return nil
	}

	uri := b.configs.URI
	if b.configs.IsTLSEnabled {
		uri = b.getTLSMongoURI()
	}

	_, err := connstring.ParseAndValidate(uri)
	if err != nil {
		logger.Error(context.TODO(), "Error validating mongodb connection uri: %s", uri)
		return err
	}

	connectionOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.Background(), connectionOptions)
	if err != nil {
		logger.Error(context.TODO(), "Error connecting to mongodb: %s", err)
		return err
	}

	b.client = client

	// Ping the MongoDB server to check if connection is established
	err = b.Ping()
	if err != nil {
		return err
	}

	logger.Info(context.TODO(), "Connected to MongoDB")

	return nil
}

func (b *ConnectionBuilder) Ping() error {
	if b.client == nil {
		return fmt.Errorf("mongodb client is not initialized")
	}

	return b.client.Ping(context.Background(), nil)
}

func (b *ConnectionBuilder) GetClient() (*mongo.Client, error) {
	if b.client != nil {
		return b.client, nil
	}

	err := b.Connect()
	if err != nil {
		return nil, err
	}

	return b.client, nil
}

// globalClient is used to access tenantInfo collection for any new tenant connection query only.
var globalClient *mongo.Client

func GetGlobalClient() (*mongo.Client, error) {
	if globalClient != nil {
		return globalClient, nil
	}

	var err error
	globalClient, err = NewConnection(DefaultConnectionConfigs()).GetClient()

	return globalClient, err
}

func GetGlobalConnection() (*mongo.Database, error) {
	client, err := GetGlobalClient()
	if err != nil {
		return nil, err
	}

	defaultDB := "starcruiser"
	if env, ok := os.LookupEnv("NODE_ENV"); ok && env == "production" {
		defaultDB = "battlecruiser"
	}

	return client.Database(defaultDB), nil
}
