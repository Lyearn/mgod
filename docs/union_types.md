---
title: Union Types
---
Sometimes its possible that the API needs to be flexible and support a range of types. An example for this might be a tagging functionality on resources such as user, movies, etc. The CRUD layer for tags entity needs to support operations on multiple types of tags like `NumberTag`, `DateTag`, etc. through same functions.

## Usage

In Go, to create union types, we need to create a base type and then use it as a struct embedding in children types. For instance, here we are creating two children tag types - `NumberTag` and `DateTag` with `BaseTag` as the base type.
```go
type TagTypeEnum string

const (
	TagTypeEnumNumber TagTypeEnum = "number"
	TagTypeEnumDate   TagTypeEnum = "date"
)

type BaseTag struct {
	ID   string `bson:"_id" mgoType:"id"`
	Name string
	Type TagTypeEnum
}

func (BaseTag) IsTag() {}

type NumberTag struct {
	BaseTag `bson:",inline"`
	Number  int
}

type DateTag struct {
	BaseTag `bson:",inline"`
	Date    string `mgoType:"date"`
}
```

Though we know the type of the doc while insertion, it might not be possible to know the doc type beforehand while querying a collection that stores multiple types of docs like in the case of union types. So, we need a global type that can receive the doc for any of the union types.
```go
type GlobalTag struct {
	BaseTag `bson:",inline"`
	Number  *int    `bson:",omitempty"`
	Date    *string `bson:",omitempty" mgoType:"date"`
}
```
> [!IMPORTANT]
> Only common fields are kept as compulsory whereas other fields are marked optional.

Configure schema options for the union type entities.
```go
discriminator := "type"
schemaOpts := schemaopt.SchemaOptions{
	Collection:       "unionTest",
	Timestamps:       true,
	IsUnionType:      true,
	DiscriminatorKey: &discriminator,
}
```
We have used the `type` field as the discriminator in this case, as it a unique key to identify underlying union types.

> [!NOTE]
> In case if we do not provide `type` as the discriminator key, `__t` will be added to the MongoDB doc with its value as the name of its underlying struct type.

Create ODM for entities using `mgod`.
```go
globalTagModel, _ := mgod.NewEntityMongoModel(GlobalTag{}, schemaOpts)
numberTagModel, _ := mgod.NewEntityMongoModel(NumberTag{}, schemaOpts)
dateTagModel, _ := mgod.NewEntityMongoModel(DateTag{}, schemaOpts)
```

Now, to insert documents, we have two options -
1. Use ODM specific to the entity we are inserting in case we have liberty to create separate functions to handle different entities.
```go
numberTag := NumberTag{
		BaseTag: BaseTag{
			ID:   primitive.NewObjectID().Hex(),
			Name: "numberTag",
			Type: TagTypeEnumNumber,
		},
		Number: 1,
	}

insertedNumberTag, _ := numberTagModel.InsertOne(context.TODO(), numberTag)
```

**Output:**
```json
{
	"_id" : ObjectId("65718f9c55e90b39cf538b42"),
	"name" : "numberTag",
	"type" : "number",
	"number" : 1,
	"createdAt" : ISODate("2023-12-07T09:25:48.253Z"),
	"updatedAt" : ISODate("2023-12-07T09:25:48.253Z"),
	"__v" : 0
}
```

2. Use global ODM to insert the doc that is created using the entity ODM. This is helpful in case where we want a common function handle the inserting of entities.
```go
date, _ := dateformatter.New(time.Now().UTC()).GetISOString()
dateTag := DateTag{
		BaseTag: BaseTag{
			ID:   primitive.NewObjectID().Hex(),
			Name: "dateTag",
			Type: TagTypeEnumDate,
		},
		Date: date,
	}

dateTagDoc, _ := dateTagModel.GetDocToInsert(context.TODO(), numberTag)

insertedDateTag, _ := globalTagModel.InsertOne(context.TODO(), dateTagDoc)
```

**Output:**
```json
{
	"_id" : ObjectId("65718f9c55e90b39cf538b43"),
	"name" : "dateTag",
	"type" : "date",
	"date" : ISODate("2023-12-07T09:25:48.252Z"),
	"createdAt" : ISODate("2023-12-07T09:25:48.253Z"),
	"updatedAt" : ISODate("2023-12-07T09:25:48.253Z"),
	"__v" : 0
}
```

Use the global ODM to find docs by querying on model properties.
```go
numberTag, _ := globalTagModel.FindOne(context.TODO(), bson.M{"name": "numberTag"})
```

**Output:**
```go
GlobalTag{
	BaseTag{
		ID: "65718f9c55e90b39cf538b42",
		Name: "numberTag",
		Type: "number",
	},
	Number: 0x1400030c380, // pointer to value of 1
	Date: <nil>,
}
```
In the above step, before returning the results, all docs received from the MongoDB are validated and processed against their respective typed models based on the discriminator key (here the `type` field). So, in the above step, the number tag document is processed against the schema for NumberTag type before getting converted to the GlobalTag type.