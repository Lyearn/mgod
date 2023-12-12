---
title: Field Transformers
sidebar_position: 3
---

Field transformers are an adapter between MongoDB field and Go struct field. They help in transforming field types in both directions i.e. from entity model to mongo doc and vice versa while building intermediate BSON document.

:::note
A field transformer is defined by the tag `mgoType`.
:::

`mgod` supports the following field transformers -

## ID

- Tag Value: `id`

It is a transformer that converts a field of type `string` in Go struct to `primitive.ObjectID` for MongoDB document and vice versa.

### Example

Type with id transformer.

```go
type User struct {
	ID   string `bson:"_id" mgoType:"id"`
	Name string
}

// id = "65697705d4cbed00e8aba717"
id := primitive.NewObjectID().Hex()
userDoc := User{
	ID: id,
	Name: "Gopher",
}

user, _ := userModel.InsertOne(context.TODO(), userDoc)
```

**Output:**

```js
{
	"_id": ObjectId("65697705d4cbed00e8aba717"),
	"name": "Gopher"
}
```

`_id_` field will be of type ObjectId instead of String in MongoDB.

Invalid user doc -

```go
userDoc := User{
	ID: "randomId"
	Name: "Gopher",
}
```

Inserting this doc will throw error as `ID` field cannot be converted to primitive.ObjectID.

---

Type without id transformer.

```go
type User struct {
	ID   string
	Name string
}

userDoc := User{
	ID: "randomId",
	Name: "Gopher",
}
```

This is a valid doc now because there is no transformer applied on `ID` field. Also, note that `ID` field will be converted to `id` instead of `_id` because BSON tag is not present.

## Date

- Tag Value: `date`

It is a transformer that converts a field of type `string` in ISO 8601 format to `primitive.DateTime` for MongoDB document and vice versa.

### Example

Type with date transformer.

```go
type User struct {
	Name     string
	JoinedOn string `bson:"joinedOn" mgoType:"date"`
}

// joinedOn = "2023-12-01T11:32:19.290Z"
joinedOn, _ := dateformatter.New(time.Now()).GetISOString()
userDoc := User{
	Name: "Gopher",
	JoinedOn: joinedOn,
}

user, _ := userModel.InsertOne(context.TODO(), userDoc)
```

**Output:**

```js
{
	"_id": ObjectId("65697705d4cbed00e8aba717"),
	"name": "Gopher",
	"joinedOn": ISODate("2023-12-01T11:32:19.290Z")
}
```

`joinedOn` field will be of type Date instead of String in MongoDB.

Invalid user doc -

```go
userDoc := User{
	Name: "Gopher",
	JoinedOn: "2023-12-01",
}
```

Inserting this doc will throw error as `JoinedOn` field is not in expected ISO 8601 format.

---

Type without date transformer.

```go
type User struct {
	Name     string
	JoinedOn string `bson:"joinedOn"`
}

userDoc := User{
	Name: "Gopher",
	JoinedOn: "2023-12-01",
}
```

This is a valid doc now because there is no transformer applied on `JoinedOn` field.
