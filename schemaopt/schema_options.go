package schemaopt

type SchemaOptions struct {
	// Collection is the name of the mongo collection in which the entity is stored.
	Collection string
	// Timestamps reports whether to add createdAt and updatedAt meta fields for the entity.
	Timestamps bool
	// EnableVersionKey reports whether to add a version key for the entity. Defaults to true.
	VersionKey *bool
	// IsUnionType reports whether the entity is a union type.
	IsUnionType bool
	// DiscriminatorKey is the key used to identify the underlying type in case of a union type entity. Defaults to __t.
	DiscriminatorKey *string // bson key
}
