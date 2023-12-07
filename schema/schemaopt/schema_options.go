package schemaopt

// SchemaOptions is Mongo Schema level options (modifies actual MongoDB doc) that needs to be provided when creating a new EntityMongoModel.
//
// These options are used to identify the collection name, whether to add timestamps meta fields, etc.
type SchemaOptions struct {
	// Collection is the name of the mongo collection in which the entity is stored.
	Collection string
	// Timestamps reports whether to add createdAt and updatedAt meta fields for the entity.
	Timestamps bool
	// VersionKey reports whether to add a version key for the entity. Defaults to true.
	VersionKey *bool
	// IsUnionType reports whether the entity is a union type.
	IsUnionType bool
	// DiscriminatorKey is the key used to identify the underlying type in case of a union type entity. Defaults to __t.
	DiscriminatorKey *string // bson key
}
