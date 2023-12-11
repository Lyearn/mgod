---
title: Schema Options
sidebar_position: 1
---

Schema Options is Mongo Schema level options (which modifies actual MongoDB doc) that needs to be provided when creating a new EntityMongoModel.

`mgod` supports the following schema options -

## Collection
* Accepts Type: `string`
* Is Optional: `No`

It is the name of the mongo collection in which the entity is stored. For example, `users` collection of MongoDB for `User` model in Golang.

**Usage**
```go
schemaOpts := schemaopt.SchemaOptions{
	Collection: "users", // MongoDB collection name
}
```

## Timestamps
* Accepts Type: `bool`
* Default Value: `false`
* Is Optional: `Yes`

It is used to track `createdAt` and `updatedAt` meta fields for the entity. See [Meta Fields](meta_fields.md) for examples.

**Usage**
```go
schemaOpts := schemaopt.SchemaOptions{
	Collection: "users",
	Timestamps: true,
}
```

## VersionKey
* Accepts Type: `bool`
* Default Value: `true`
* Is Optional: `Yes`

This reports whether to add a version key (`__v`) for the entity. See [Meta Fields](meta_fields.md) for examples.

**Usage**
```go
schemaOpts := schemaopt.SchemaOptions{
	Collection: "users",
	VersionKey: true,
}
```

## IsUnionType
* Accepts Type: `bool`
* Default Value: `false`
* Is Optional: `Yes`

It defines whether the entity is a union type. See [Union Types](union_types.md) for more details on unions.

**Usage**
```go
schemaOpts := schemaopt.SchemaOptions{
	Collection: "resources",
	IsUnionType: true,
}
```
If `IsUnionType` is set to true, then `__t` will be used as the `DiscriminatorKey` by default.

## DiscriminatorKey
* Accepts Type: `string`
* Default Value: `__t`
* Is Optional: `Yes`

It is the key used to identify the underlying type in case of a union type entity.

**Usage**

:::note
`IsUnionType` needs to be set to `true` to use the `DiscriminatorKey` field.
:::

```go
schemaOpts := schemaopt.SchemaOptions{
	Collection: "resources",
	IsUnionType: true,
	DiscriminatorKey: "type",
}
```

:::info
The provided `DiscriminatorKey` should be present in the Go struct as a compulsory field.
:::

Default `DiscriminatorKey` will be overwritten by the provided `type` field.