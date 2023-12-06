package schema_test

import (
	"reflect"
	"testing"

	"github.com/Lyearn/mgod/schema"
	"github.com/Lyearn/mgod/schema/fieldopt"
	"github.com/Lyearn/mgod/schema/schemaopt"
	"github.com/Lyearn/mgod/schema/transformer"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
)

type EntityModelSchemaSuite struct {
	suite.Suite
	*require.Assertions
}

func TestEntityModelSchemaSuite(t *testing.T) {
	s := new(EntityModelSchemaSuite)
	suite.Run(t, s)
}

func (s *EntityModelSchemaSuite) SetupTest() {
	s.Assertions = require.New(s.T())
}

func (s *EntityModelSchemaSuite) TestBuildSchemaForModel() {
	type TestCase struct {
		Model                 any
		Schema                schema.EntityModelSchema
		SchemaOpts            schemaopt.SchemaOptions
		ValidSchemaNodesCount int
	}

	type UserProject struct {
		ProjectID   string `mgoType:"id"`
		CompletedAt string `bson:"completedAt" mgoType:"date"`
	}

	type Metadata struct {
		JoinedOn  string        `mgoType:"date"`
		TeamIDs   []string      `bson:"teamIds" mgoType:"id"`
		Projects  []UserProject `bson:"projects" mgoID:"false" mgoDefault:"[]"`
		SkipField string        `bson:"-"`
	}

	type NestedInlineProps struct {
		Name        string
		InlineFloat float64 `bson:"inlineFloat"`
	}

	type InlineProps struct {
		NestedInlineProps `bson:",inline"`

		ID           string `bson:"_id" mgoType:"id"`
		InlineBool   bool   `bson:"inlineBool"`
		Metadata     string `bson:"meta"`
		InlineString string
	}

	type NestedModelWithAllTypes struct {
		ID        string      `bson:"_id" mgoType:"id"`
		Name      *string     `bson:",omitempty"`
		Age       int         `mgoDefault:"18"`
		Metadata  *Metadata   `bson:"meta"`
		Props     InlineProps `bson:",inline"`
		Height    float64
		SkipField bool `bson:"-"`
	}

	rootNode := schema.GetDefaultSchemaTreeRootNode()

	// TC 1: NestedModelWithAllTypes
	nestedModelWithAllTypesRootNode := rootNode
	nestedModelWithAllTypesRootNode.Children = []schema.TreeNode{
		{
			Path:    "$root._id",
			BSONKey: "_id",
			Key:     "ID",
			Props: schema.SchemaFieldProps{
				Type:         reflect.String,
				Transformers: []transformer.Transformer{transformer.IDTransformer},
				Options: fieldopt.SchemaFieldOptions{
					Required: true,
				},
			},
		},
		{
			Path:    "$root.name",
			BSONKey: "name",
			Key:     "Name",
			Props: schema.SchemaFieldProps{
				Type:         reflect.String,
				IsPointer:    true,
				Transformers: []transformer.Transformer{},
				Options: fieldopt.SchemaFieldOptions{
					Required: false,
				},
			},
		},
		{
			Path:    "$root.age",
			BSONKey: "age",
			Key:     "Age",
			Props: schema.SchemaFieldProps{
				Type:         reflect.Int,
				Transformers: []transformer.Transformer{},
				Options: fieldopt.SchemaFieldOptions{
					Required: true,
					Default:  18,
				},
			},
		},
		{
			Path:    "$root.meta",
			BSONKey: "meta",
			Key:     "Metadata",
			Props: schema.SchemaFieldProps{
				Type:         reflect.Struct,
				IsPointer:    true,
				Transformers: []transformer.Transformer{},
				Options: fieldopt.SchemaFieldOptions{
					XID:      true,
					Required: true,
				},
			},
			Children: []schema.TreeNode{
				{
					Path:    "$root.meta.joinedon",
					BSONKey: "joinedon",
					Key:     "JoinedOn",
					Props: schema.SchemaFieldProps{
						Type:         reflect.String,
						Transformers: []transformer.Transformer{transformer.DateTransformer},
						Options: fieldopt.SchemaFieldOptions{
							Required: true,
						},
					},
				},
				{
					Path:    "$root.meta.teamIds",
					BSONKey: "teamIds",
					Key:     "TeamIDs",
					Props: schema.SchemaFieldProps{
						Type:         reflect.Slice,
						Transformers: []transformer.Transformer{},
						Options: fieldopt.SchemaFieldOptions{
							Required: true,
						},
					},
					Children: []schema.TreeNode{
						{
							Path:    "$root.meta.teamIds.$",
							BSONKey: "$",
							Key:     "$", // to identify slice element
							Props: schema.SchemaFieldProps{
								Type:         reflect.String,
								Transformers: []transformer.Transformer{transformer.IDTransformer},
							},
						},
					},
				},
				{
					Path:    "$root.meta.projects",
					BSONKey: "projects",
					Key:     "Projects",
					Props: schema.SchemaFieldProps{
						Type:         reflect.Slice,
						Transformers: []transformer.Transformer{},
						Options: fieldopt.SchemaFieldOptions{
							Required: true,
							Default:  bson.A{},
						},
					},
					Children: []schema.TreeNode{
						{
							Path:    "$root.meta.projects.$",
							BSONKey: "$",
							Key:     "$",
							Props: schema.SchemaFieldProps{
								Type:         reflect.Struct,
								Transformers: []transformer.Transformer{},
							},
							Children: []schema.TreeNode{
								{
									Path:    "$root.meta.projects.$.projectid",
									BSONKey: "projectid",
									Key:     "ProjectID",
									Props: schema.SchemaFieldProps{
										Type:         reflect.String,
										Transformers: []transformer.Transformer{transformer.IDTransformer},
										Options: fieldopt.SchemaFieldOptions{
											Required: true,
										},
									},
								},
								{
									Path:    "$root.meta.projects.$.completedAt",
									BSONKey: "completedAt",
									Key:     "CompletedAt",
									Props: schema.SchemaFieldProps{
										Type:         reflect.String,
										Transformers: []transformer.Transformer{transformer.DateTransformer},
										Options: fieldopt.SchemaFieldOptions{
											Required: true,
										},
									},
								},
							},
						},
					},
				},
				{
					Path:    "$root.meta._id",
					BSONKey: "_id",
					Key:     "XID",
					Props: schema.SchemaFieldProps{
						Type:         reflect.String,
						Transformers: []transformer.Transformer{transformer.IDTransformer},
						Options: fieldopt.SchemaFieldOptions{
							Required: true,
						},
					},
				},
			},
		},
		{
			Path:    "$root.inlineFloat",
			BSONKey: "inlineFloat",
			Key:     "InlineFloat",
			Props: schema.SchemaFieldProps{
				Type:         reflect.Float64,
				Transformers: []transformer.Transformer{},
				Options: fieldopt.SchemaFieldOptions{
					Required: true,
				},
			},
		},
		{
			Path:    "$root.inlineBool",
			BSONKey: "inlineBool",
			Key:     "InlineBool",
			Props: schema.SchemaFieldProps{
				Type:         reflect.Bool,
				Transformers: []transformer.Transformer{},
				Options: fieldopt.SchemaFieldOptions{
					Required: true,
				},
			},
		},
		{
			Path:    "$root.inlinestring",
			BSONKey: "inlinestring",
			Key:     "InlineString",
			Props: schema.SchemaFieldProps{
				Type:         reflect.String,
				Transformers: []transformer.Transformer{},
				Options: fieldopt.SchemaFieldOptions{
					Required: true,
				},
			},
		},
		{
			Path:    "$root.height",
			BSONKey: "height",
			Key:     "Height",
			Props: schema.SchemaFieldProps{
				Type:         reflect.Float64,
				Transformers: []transformer.Transformer{},
				Options: fieldopt.SchemaFieldOptions{
					Required: true,
				},
			},
		},
	}
	nestedModelWithAllTypesSchema := schema.EntityModelSchema{Root: nestedModelWithAllTypesRootNode}

	versionKeyEnabled := false
	nestedModelWithAllTypesTestCase := TestCase{
		Model:  NestedModelWithAllTypes{},
		Schema: nestedModelWithAllTypesSchema,
		SchemaOpts: schemaopt.SchemaOptions{
			VersionKey: &versionKeyEnabled,
		},
		ValidSchemaNodesCount: 17,
	}

	// TC 2: NestedModelWithAllTypes with all schema options enabled
	nestedModelWithSchemaOptsRootNode := nestedModelWithAllTypesRootNode

	nestedModelWithSchemaOptsRootNode.Children = append(nestedModelWithSchemaOptsRootNode.Children, schema.TreeNode{
		Path:    "$root.createdAt",
		BSONKey: "createdAt",
		Key:     "createdAt",
		Props: schema.SchemaFieldProps{
			Type:         reflect.String,
			Options:      fieldopt.SchemaFieldOptions{Required: false},
			Transformers: []transformer.Transformer{transformer.DateTransformer},
		},
	})

	nestedModelWithSchemaOptsRootNode.Children = append(nestedModelWithSchemaOptsRootNode.Children, schema.TreeNode{
		Path:    "$root.updatedAt",
		BSONKey: "updatedAt",
		Key:     "updatedAt",
		Props: schema.SchemaFieldProps{
			Type:         reflect.String,
			Options:      fieldopt.SchemaFieldOptions{Required: false},
			Transformers: []transformer.Transformer{transformer.DateTransformer},
		},
	})

	nestedModelWithSchemaOptsRootNode.Children = append(nestedModelWithSchemaOptsRootNode.Children, schema.TreeNode{
		Path:    "$root.__v",
		BSONKey: "__v",
		Key:     "__v",
		Props: schema.SchemaFieldProps{
			Type:         reflect.Int,
			Options:      fieldopt.SchemaFieldOptions{Required: false},
			Transformers: []transformer.Transformer{},
		},
	})

	nestedModelWithAllTypesAndSchemaOptsSchema := schema.EntityModelSchema{Root: nestedModelWithSchemaOptsRootNode}

	nestedModelWithAllTypesAndSchemaOptsTestCase := TestCase{
		Model:  NestedModelWithAllTypes{},
		Schema: nestedModelWithAllTypesAndSchemaOptsSchema,
		SchemaOpts: schemaopt.SchemaOptions{
			Timestamps: true,
			// version key enabled by default
		},
		ValidSchemaNodesCount: 20,
	}

	testCases := []TestCase{
		nestedModelWithAllTypesTestCase,
		nestedModelWithAllTypesAndSchemaOptsTestCase,
	}

	for _, tc := range testCases {
		actualSchema, err := schema.BuildSchemaForModel(tc.Model, tc.SchemaOpts)

		s.Nil(err)
		if !reflect.DeepEqual(tc.Schema.Root, actualSchema.Root) {
			s.Fail("Schema mismatch", "Expected: %v, Got: %v", tc.Schema.Root, actualSchema.Root)
		}
		s.Equal(tc.ValidSchemaNodesCount, len(actualSchema.Nodes))
	}
}
