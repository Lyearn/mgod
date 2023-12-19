---
title: Basic Usage
---

Use existing MongoDB connection, or setup a new one to register a default database connection.

For existing database connection,
```go
import "github.com/Lyearn/mgod"

func init() {
	dbName := "mgod-test"

	// client is the MongoDB client obtained using Go Mongo Driver's Connect method.
	mgod.SetDefaultConnection(client, dbName)
}
```

To setup a new connection,
```go
import (
	"time"

	"github.com/Lyearn/mgod"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	// `cfg` is optional. Can rely on default configurations by providing `nil` value in argument.
	cfg := &mgod.ConnectionConfig{Timeout: 5 * time.Second}
	dbName := "mgod-test"
	opts := options.Client().ApplyURI("mongodb://root:mgod123@localhost:27017")

	err := mgod.ConfigureDefaultConnection(cfg, dbName, opts)
}
```

:::note
The above `err` variable will be a connection error (if occurs) returned by the Go Mongo Driver. So, handle the error accordingly.
:::

Add tags _(wherever applicable)_ in existing struct _(or define a new model)_.

```go
type User struct {
	Name     string
	EmailID  string `bson:"emailId"`
	Age      *int32 `bson:",omitempty"`
	JoinedOn string `bson:"joinedOn" mgoType:"date"`
}
```

Use `mgod` to get the entity ODM.

```go
import (
	"github.com/Lyearn/mgod"
	"github.com/Lyearn/mgod/schema/schemaopt"
)

model := User{}
schemaOpts := schemaopt.SchemaOptions{
	Collection: "users",
	Timestamps: true,
}

userModel, _ := mgod.NewEntityMongoModel(model, schemaOpts)
```

Use the entity ODM to perform CRUD operations with ease.

## Inserting a new document

```go
joinedOn, _ := dateformatter.New(time.Now()).GetISOString()
userDoc := User{
	Name: "Gopher",
	EmailID: "gopher@mgod.com",
	JoinedOn: joinedOn,
}
user, _ := userModel.InsertOne(context.TODO(), userDoc)
```

**Output:**

```js
{
	"_id": ObjectId("65697705d4cbed00e8aba717"),
	"name": "Gopher",
	"emailId": "gopher@mgod.com",
	"joinedOn": ISODate("2023-12-01T11:32:19.290Z"),
	"createdAt": ISODate("2023-12-01T11:32:19.290Z"),
	"updatedAt": ISODate("2023-12-01T11:32:19.290Z"),
	"__v": 0
}
```

Notice how `_id`, `createdAt`, `updatedAt` and `__v` fields are added automatically.

## Finding documents using model properties

```go
users, _ := userModel.Find(context.TODO(), bson.M{"name": userDoc.Name})
```

**Output:**

```go
[]User{
	User{
		Name: "Gopher",
		EmailID: "gopher@mgod.com",
		JoinedOn: "2023-12-01T11:32:19.290Z",
	}
}
```

## Updating document properties

```go
result, _ := userModel.UpdateMany(context.TODO(), bson.M{"joinedOn": bson.M{"$gte": "2023-12-01T00:00:00.000Z"}}, bson.M{"$inc": {"__v": 1}})
```

**Output:**

```go
mongo.UpdateResult{
	MatchedCount: 1,
	ModifiedCount: 1,
	UpsertedCount: 0,
}
```

```js
// User Doc
{
	"_id": ObjectId("65697705d4cbed00e8aba717"),
	"name": "Gopher",
	"emailId": "gopher@mgod.com",
	"joinedOn": ISODate("2023-12-01T11:32:19.290Z"),
	"createdAt": ISODate("2023-12-01T11:32:19.290Z"),
	"updatedAt": ISODate("2023-12-02T10:40:00.670Z"),
	"__v": 1
}
```

Notice the updation of the `updatedAt` field.

## Removing documents matching certain or all model properties

```go
result, _ := userModel.DeleteMany(context.TODO(), bson.M{"name": userDoc.Name})
```

**Output:**

```go
mongo.DeleteResult{
	DeletedCount: 1
}
```
