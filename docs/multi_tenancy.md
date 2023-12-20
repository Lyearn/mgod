---
title: Multi Tenancy
---

`mgod` comes with the built-in support for multi-tenancy, enabling the use of a single Go struct with multiple databases. This feature allows creation of multiple `EntityMongoModel` of the same Go struct to be attached to different databases while using the same underlying MongoDB client connection.

## Usage

Create separate `EntityMongoModel` for different tenants using same Go struct and corresponding databases.

```go
type User struct {
	Name    string
	EmailID string `bson:"emailId"`
	Amount  float32
}
collection := "users"

tenant1DB := "tenant1"
tenant2DB := "tenant2"

tenant1Model, _ := mgod.NewEntityMongoModelOptions(tenant1DB, collection, nil)
tenant2Model, _ := mgod.NewEntityMongoModelOptions(tenant2DB, collection, nil)
```

These models can now be used simultaneously inside the same service logic as well as in a transaction operation.

```go
amount := 10000

tenant1Model.UpdateMany(context.TODO(), bson.M{"name": "Gopher Tenant 1"}, bson.M{"$inc": {"amount": -amount}})
tenant2Model.UpdateMany(context.TODO(), bson.M{"name": "Gopher Tenant 2"}, bson.M{"$inc": {"amount": amount}})
```

:::note
The `EntityMongoModel` is always bound to the specified database at the time of its declaration and, as such, cannot be used to perform operations across multiple databases.
:::

```go
amount := 10000

result, _ := tenant1Model.FindOne(context.TODO(), bson.M{"name": "Gopher Tenant 2"})
// result will be <nil> value in this case
```
