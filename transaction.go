package mgod

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// TransactionFunc is the function that is executed in a MongoDB transaction.
//
// SessionContext(sc) combines the context.Context and mongo.Session interfaces.
type TransactionFunc func(sc mongo.SessionContext) (interface{}, error)

// WithTransaction executes the given transaction function with a new session.
func WithTransaction(ctx context.Context, transactionFunc TransactionFunc) (interface{}, error) {
	session, err := mClient.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	// Reason behind using read preference:
	// https://www.mongodb.com/community/forums/t/why-can-t-read-preference-be-secondary-in-a-transaction/204432
	payload, transactionErr := session.WithTransaction(ctx, transactionFunc, &options.TransactionOptions{
		ReadPreference: readpref.Primary(),
	})

	return payload, transactionErr
}
