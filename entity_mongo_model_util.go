package mgod

import "strings"

// GetSchemaCacheKey returns the cache key for the schema of a model.
// The cache key is in the format of <collection_name>_<model_name_or_value>.
// Unless we are storing the schema of a union type based on the discriminator value, the second parameter
// will be the model name only.
func GetSchemaCacheKey(collName string, modelNameOrVal string) string {
	keyElems := []string{collName, modelNameOrVal}
	key := strings.Join(keyElems, "_")

	return key
}
