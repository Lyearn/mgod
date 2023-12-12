---
title: Meta Fields
sidebar_position: 4
---

Meta fields are those fields that tracks extra information about the document which can be helpful to determine the state of a document.

`mgod` supports the following meta fields -

We will assume the following `User` type for the rest of this section -

```go
type User struct {
	Name     string
	EmailID  string `bson:"emailId"`
}
```

## createdAt

- Format - `ISO String`

It is the meta field that stores the timestamp of the document creation. This field is automatically added (if not present in the input) to the schema if the `SchemaOptions.Timestamps` is set to true.

### Example

```go
schemaOpts := schemaopt.SchemaOptions{
	Collection: "users",
	Timestamps: true,
}

userDoc := User{
	Name: "Gopher",
	EmailID: "gopher@mgod.com",
}
user, _ := userModel.InsertOne(context.TODO(), userDoc)
```

**Output:**

```js
{
	"_id": ObjectId("65697705d4cbed00e8aba717"),
	"name": "Gopher",
	"emailId": "gopher@mgod.com",
	"createdAt": ISODate("2023-12-01T11:32:19.290Z"),
	"updatedAt": ISODate("2023-12-01T11:32:19.290Z"),
	"__v": 0
}
```

## updatedAt

- Format - `ISO String`

It is the meta field that stores the timestamp of the document updation. This field is automatically upserted (if not present in the input) to the schema if the `SchemaOptions.Timestamps` is set to true.

### Example

```go
schemaOpts := schemaopt.SchemaOptions{
	Collection: "users",
	Timestamps: true,
}

result, _ := userModel.UpdateMany(context.TODO(), bson.M{"emailId": "gopher@mgod.com"}, bson.M{"$inc": {"__v": 1}})
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
	"createdAt": ISODate("2023-12-01T11:32:19.290Z"),
	"updatedAt": ISODate("2023-12-02 10:40:00.670Z"),
	"__v": 1
}
```

## \_\_v

- Format - `Number`

It is the field that stores the version of the document. This field is automatically added (if not present in the input) to the schema if the `SchemaOptions.VersionKey` is set to true. This field starts with a default value of 0.

### Example

```go
schemaOpts := schemaopt.SchemaOptions{
	Collection: "users",
	VersionKey: true
}

userDoc := User{
	Name: "Gopher",
	EmailID: "gopher@mgod.com",
}
user, _ := userModel.InsertOne(context.TODO(), userDoc)
```

**Output:**

```js
{
	"_id": ObjectId("65697705d4cbed00e8aba717"),
	"name": "Gopher",
	"emailId": "gopher@mgod.com",
	"__v": 0
}
```

If `VersionKey` is set to `false`.

```go
schemaOpts := schemaopt.SchemaOptions{
	Collection: "users",
	VersionKey: false
}

userDoc := User{
	Name: "Gopher",
	EmailID: "gopher@mgod.com",
}
user, _ := userModel.InsertOne(context.TODO(), userDoc)
```

**Output:**

```js
{
	"_id": ObjectId("65697705d4cbed00e8aba717"),
	"name": "Gopher",
	"emailId": "gopher@mgod.com"
}
```
