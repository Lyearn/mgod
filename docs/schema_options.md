---
title: Schema Options
---

Schema Options is Mongo Schema level options (which modifies actual MongoDB doc) that needs to be provided when creating a new EntityMongoModel.

`mgod` supports the following schema options -

## Collection
* Accepts Type: `string`
* Is Optional: `No`

It is the name of the mongo collection in which the entity is stored. For example, `users` collection of MongoDB for `User` model in Golang.

## Timestamps
* Accepts Type: `bool`
* Default Value: `false`
* Is Optional: `Yes`

This reports whether to track `createdAt` and `updatedAt` meta fields for the entity. See `Meta Fields` sections for examples.

## VersionKey
* Accepts Type: `bool`
* Default Value: `true`
* Is Optional: `Yes`

This reports whether to add a version key (`__v`) for the entity. See `Meta Fields` sections for examples.

## IsUnionType
* Accepts Type: `bool`
* Default Value: `false`
* Is Optional: `Yes`

This reports whether the entity is a union type. See `Union Types` section for more details on unions.

## DiscriminatorKey
* Accepts Type: `string`
* Default Value: `__t`
* Is Optional: `Yes`

It is the key used to identify the underlying type in case of a union type entity.