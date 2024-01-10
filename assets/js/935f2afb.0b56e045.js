"use strict";(self.webpackChunkmgod=self.webpackChunkmgod||[]).push([[53],{1109:e=>{e.exports=JSON.parse('{"pluginId":"default","version":"current","label":"Next","banner":null,"badge":false,"noIndex":false,"className":"docs-version-current","isLast":true,"docsSidebars":{"docsSidebar":[{"type":"category","label":"Introduction","items":[{"type":"link","label":"About","href":"/mgod/docs/about","docId":"about","unlisted":false}],"collapsed":false,"collapsible":true},{"type":"category","label":"Beginner\'s Guide","items":[{"type":"link","label":"Installation","href":"/mgod/docs/installation","docId":"installation","unlisted":false},{"type":"link","label":"Basic Usage","href":"/mgod/docs/basic_usage","docId":"basic_usage","unlisted":false}],"collapsed":false,"collapsible":true},{"type":"category","label":"Features","items":[{"type":"link","label":"Schema Options","href":"/mgod/docs/schema_options","docId":"schema_options","unlisted":false},{"type":"link","label":"Field Options","href":"/mgod/docs/field_options","docId":"field_options","unlisted":false},{"type":"link","label":"Field Transformers","href":"/mgod/docs/field_transformers","docId":"field_transformers","unlisted":false},{"type":"link","label":"Meta Fields","href":"/mgod/docs/meta_fields","docId":"meta_fields","unlisted":false}],"collapsed":false,"collapsible":true},{"type":"category","label":"Advanced Guide","items":[{"type":"link","label":"Multi Tenancy","href":"/mgod/docs/multi_tenancy","docId":"multi_tenancy","unlisted":false},{"type":"link","label":"Union Types","href":"/mgod/docs/union_types","docId":"union_types","unlisted":false},{"type":"link","label":"Transactions","href":"/mgod/docs/transactions","docId":"transactions","unlisted":false}],"collapsed":false,"collapsible":true}]},"docs":{"about":{"id":"about","title":"About","description":"What is mgod?","sidebar":"docsSidebar"},"basic_usage":{"id":"basic_usage","title":"Basic Usage","description":"Use existing MongoDB connection, or setup a new one to register a default database connection.","sidebar":"docsSidebar"},"field_options":{"id":"field_options","title":"Field Options","description":"Field Options are custom schema options available at field level (for fields of struct type). These options either modifies the schema or adds validations to the field on which it is applied.","sidebar":"docsSidebar"},"field_transformers":{"id":"field_transformers","title":"Field Transformers","description":"Field transformers are an adapter between MongoDB field and Go struct field. They help in transforming field types in both directions i.e. from entity model to mongo doc and vice versa while building intermediate BSON document.","sidebar":"docsSidebar"},"installation":{"id":"installation","title":"Installation","description":"Requirements","sidebar":"docsSidebar"},"meta_fields":{"id":"meta_fields","title":"Meta Fields","description":"Meta fields are those fields that tracks extra information about the document which can be helpful to determine the state of a document.","sidebar":"docsSidebar"},"multi_tenancy":{"id":"multi_tenancy","title":"Multi Tenancy","description":"mgod comes with the built-in support for multi-tenancy, enabling the use of a single Go struct with multiple databases. This feature allows creation of multiple EntityMongoModel of the same Go struct to be attached to different databases while using the same underlying MongoDB client connection.","sidebar":"docsSidebar"},"schema_options":{"id":"schema_options","title":"Schema Options","description":"Schema Options is Mongo Schema level options (which modifies actual MongoDB doc) that needs to be provided when creating a new EntityMongoModel.","sidebar":"docsSidebar"},"transactions":{"id":"transactions","title":"Transactions","description":"mgod provides a wrapper function WithTransaction that supports MongoDB transactions, allowing users to perform a series of read and write operations as a single atomic unit.","sidebar":"docsSidebar"},"union_types":{"id":"union_types","title":"Union Types","description":"Sometimes its possible that the API needs to be flexible and support a range of types. An example for this might be a tagging functionality on resources such as user, movies, etc. The CRUD layer for tags entity needs to support operations on multiple types of tags like NumberTag, DateTag, etc. through same functions.","sidebar":"docsSidebar"}}}')}}]);