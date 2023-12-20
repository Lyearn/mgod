package schemaopt

// SchemaOptions is Mongo Schema level options (modifies actual MongoDB doc) that needs to be provided when creating a new EntityMongoModel.
type SchemaOptions struct {
	// Timestamps reports whether to add createdAt and updatedAt meta fields for the entity.
	Timestamps bool
	// VersionKey reports whether to add a version key for the entity. Defaults to true.
	VersionKey *bool
	// IsUnionType reports whether the entity is a union type.
	IsUnionType bool
	// DiscriminatorKey is the key used to identify the underlying type in case of a union type entity. Defaults to __t.
	DiscriminatorKey *string // bson key
}
