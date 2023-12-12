---
title: Field Options
sidebar_position: 2
---

Field Options are custom schema options available at field level (for fields of struct type). These options either modifies the schema or adds validations to the field on which it is applied.

`mgod` supports the following field options -

## \_\_id

- BSON Tag: `mgoID`
- Accepts Type: `bool`
- Default Value: `true` for custom type fields

It defines if `_id` field needs to be added in a object.

:::note
This option is only applicable for fields holding structs.
:::

### Example

```go
type UserProject struct {
	Name string
}

type Metadata struct {
	TeamIDs  []string      `bson:"teamIds"`
	Projects []UserProject `mgoID:"false"`
}

type User struct {
	Name     string
	Metadata *Metadata `bson:"meta"`
}

userDoc := User{
	Name: "Gopher",
	Metadata: Metadata{
		TeamIDs: []string{"team1", "team2"},
		Projects: []UserProject{
			UserProject{Name: "Project 1"},
			UserProject{Name: "Project 2"},
		},
	},
}

user, _ := userModel.InsertOne(context.TODO(), userDoc)
```

**Output:**

```js
{
	"_id": ObjectId("65697705d4cbed00e8aba717"),
	"name": "Gopher",
	"meta": {
		"_id": ObjectId("65697705d4cbed00e8aba718"),
		"teamIds": ["team1", "team2"],
		"projects": [
			{
				"name": "Project 1"
			},
			{
				"name": "Project 2"
			}
		]
	}
}
```

See how `_id` field is added for `meta` object because `mgoID` is true by default for struct type fields. Also, note how `_id` field is skipped for `projects` object as it was set to false explicitly in type declaration.

## required

- Accepts Type: `bool`
- Default Value: `true`

It defines if a field is required or not. The option can be invalidated using `omitempty` property of `bson` tag.

### Example

```go
type User struct {
	Name   string
	Age    int32
	Height *float `bson:",omitempty"`
}
```

In the above type, height field is set to not required.

```go
userDoc := User{
	Name: "Gopher",
}
```

The above doc will throw error because `Age` field is required.

```go
userDoc := User{
	Name: "Gopher",
	Age: 18,
}
```

This doc will work fine.

## default

- BSON Tag: `mgoDefault`
- Accepts Type: `string`
- Default Value: `nil`

It provides the default value for a field. The value of this option is used when the field is not present in the input document.

:::note
This option is applicable only for fields that are not of custom type (custom structs).
:::

### Example

```go
type UserProject struct {
	Name string
}

type Metadata struct {
	TeamIDs  []string      `bson:"teamIds"`
	Projects *[]UserProject `mgoID:"false" mgoDefault:"[]"`
}

type User struct {
	Name     string
	Age      *int32 `mgoDefault:"18"`
	Metadata *Metadata `bson:"meta"`
}

userDoc := User{
	Name: "Gopher",
	Age: 30,
	Metadata: Metadata{
		TeamIDs: []string{"team1", "team2"},
	},
}

user, _ := userModel.InsertOne(context.TODO(), userDoc)
```

**Output:**

```js
{
	"_id": ObjectId("65697705d4cbed00e8aba717"),
	"name": "Gopher",
	"age": 30,
	"meta": {
		"_id": ObjectId("65697705d4cbed00e8aba718"),
		"teamIds": ["team1", "team2"],
		"projects": []
	}
}
```

See how the value of `age` field was used because it was provided in the input doc and how the default value of `projects` field is used because it was missing from the input doc.
