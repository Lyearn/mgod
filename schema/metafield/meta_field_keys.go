package metafield

// MetaFieldKey is the unique field name of a meta field.
type MetaFieldKey string

const (
	MetaFieldKeyCreatedAt  MetaFieldKey = "createdAt"
	MetaFieldKeyUpdatedAt  MetaFieldKey = "updatedAt"
	MetaFieldKeyDocVersion MetaFieldKey = "__v"
)
