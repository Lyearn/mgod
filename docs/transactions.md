---
title: Transactions
---

`mgod` provides a wrapper function `WithTransaction` that supports MongoDB transactions, allowing users to perform a series of read and write operations as a single atomic unit.

## Usage

Configure default connection with `mgod`.

```go
cfg := &mgod.ConnectionConfig{Timeout: 5 * time.Second}
opts := options.Client().ApplyURI("mongodb://localhost:27017/?replicaSet=mgod_rs&authSource=admin")

err := mgod.ConfigureDefaultClient(cfg, opts)
```

:::info
To use Transactions, it is compulsory to run MongoDB daemon as a replica set.
Refer Community Forum Discussion - [Why replica set is mandatory for transactions in MongoDB?](https://www.mongodb.com/community/forums/t/why-replica-set-is-mandatory-for-transactions-in-mongodb/9533)
:::

Create models to be used inside a MongoDB transaction.

```go
type User struct {
	Name    string
	EmailID string `bson:"emailId"`
}

dbName := "mgoddb"
collection := "users"
schemaOpts := schemaopt.SchemaOptions{
	Timestamps: true,
}

userModel, _ := mgod.NewEntityMongoModelOptions(dbName, collection, &schemaOpts)
```

Use `WithTransaction` function to perform multiple CRUD operations as an atomic unit.

```go
userDoc1 := User{Name: "Gopher1", EmailID: "gopher1@mgod.com"}
userDoc2 := User{Name: "Gopher2", EmailID: "gopher2@mgod.com"}

_, err := mgod.WithTransaction(context.Background(), func(sc mongo.SessionContext) (interface{}, error) {
	_, err1 := s.userModel.InsertOne(sc, userDoc1)
	_, err2 := s.userModel.InsertOne(sc, userDoc2)

	if err1 != nil || err2 != nil {
		return nil, errors.New("abort transaction")
	}

	return nil, nil
})
```

:::warning
Make sure to pass the session's context (`sc` here) only in EntityMongoModel's operation functions.
:::
