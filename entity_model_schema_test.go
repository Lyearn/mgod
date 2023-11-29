package mgod_test

import (
	"reflect"
	"testing"

	"github.com/Lyearn/mgod"
	"github.com/Lyearn/mgod/fieldopt"
	"github.com/Lyearn/mgod/schemaopt"
	"github.com/Lyearn/mgod/transformer"
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
		Schema                mgod.EntityModelSchema
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

	rootNode := mgod.GetDefaultSchemaTreeRootNode()

	// TC 1: NestedModelWithAllTypes
	nestedModelWithAllTypesRootNode := rootNode
	nestedModelWithAllTypesRootNode.Children = []mgod.TreeNode{
		{
			Path:    "$root._id",
			BSONKey: "_id",
			Key:     "ID",
			Props: mgod.SchemaFieldProps{
				Type:         reflect.String,
				Transformers: []transformer.Transformer{transformer.IDTransformerInstance},
				Options: fieldopt.SchemaFieldOptions{
					Required: true,
				},
			},
		},
		{
			Path:    "$root.name",
			BSONKey: "name",
			Key:     "Name",
			Props: mgod.SchemaFieldProps{
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
			Props: mgod.SchemaFieldProps{
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
			Props: mgod.SchemaFieldProps{
				Type:         reflect.Struct,
				IsPointer:    true,
				Transformers: []transformer.Transformer{},
				Options: fieldopt.SchemaFieldOptions{
					XID:      true,
					Required: true,
				},
			},
			Children: []mgod.TreeNode{
				{
					Path:    "$root.meta.joinedon",
					BSONKey: "joinedon",
					Key:     "JoinedOn",
					Props: mgod.SchemaFieldProps{
						Type:         reflect.String,
						Transformers: []transformer.Transformer{transformer.DateTransformerInstance},
						Options: fieldopt.SchemaFieldOptions{
							Required: true,
						},
					},
				},
				{
					Path:    "$root.meta.teamIds",
					BSONKey: "teamIds",
					Key:     "TeamIDs",
					Props: mgod.SchemaFieldProps{
						Type:         reflect.Slice,
						Transformers: []transformer.Transformer{},
						Options: fieldopt.SchemaFieldOptions{
							Required: true,
						},
					},
					Children: []mgod.TreeNode{
						{
							Path:    "$root.meta.teamIds.$",
							BSONKey: "$",
							Key:     "$", // to identify slice element
							Props: mgod.SchemaFieldProps{
								Type:         reflect.String,
								Transformers: []transformer.Transformer{transformer.IDTransformerInstance},
							},
						},
					},
				},
				{
					Path:    "$root.meta.projects",
					BSONKey: "projects",
					Key:     "Projects",
					Props: mgod.SchemaFieldProps{
						Type:         reflect.Slice,
						Transformers: []transformer.Transformer{},
						Options: fieldopt.SchemaFieldOptions{
							Required: true,
							Default:  bson.A{},
						},
					},
					Children: []mgod.TreeNode{
						{
							Path:    "$root.meta.projects.$",
							BSONKey: "$",
							Key:     "$",
							Props: mgod.SchemaFieldProps{
								Type:         reflect.Struct,
								Transformers: []transformer.Transformer{},
							},
							Children: []mgod.TreeNode{
								{
									Path:    "$root.meta.projects.$.projectid",
									BSONKey: "projectid",
									Key:     "ProjectID",
									Props: mgod.SchemaFieldProps{
										Type:         reflect.String,
										Transformers: []transformer.Transformer{transformer.IDTransformerInstance},
										Options: fieldopt.SchemaFieldOptions{
											Required: true,
										},
									},
								},
								{
									Path:    "$root.meta.projects.$.completedAt",
									BSONKey: "completedAt",
									Key:     "CompletedAt",
									Props: mgod.SchemaFieldProps{
										Type:         reflect.String,
										Transformers: []transformer.Transformer{transformer.DateTransformerInstance},
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
					Props: mgod.SchemaFieldProps{
						Type:         reflect.String,
						Transformers: []transformer.Transformer{transformer.IDTransformerInstance},
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
			Props: mgod.SchemaFieldProps{
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
			Props: mgod.SchemaFieldProps{
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
			Props: mgod.SchemaFieldProps{
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
			Props: mgod.SchemaFieldProps{
				Type:         reflect.Float64,
				Transformers: []transformer.Transformer{},
				Options: fieldopt.SchemaFieldOptions{
					Required: true,
				},
			},
		},
	}
	nestedModelWithAllTypesSchema := mgod.EntityModelSchema{Root: nestedModelWithAllTypesRootNode}

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

	nestedModelWithSchemaOptsRootNode.Children = append(nestedModelWithSchemaOptsRootNode.Children, mgod.TreeNode{
		Path:    "$root.createdAt",
		BSONKey: "createdAt",
		Key:     "createdAt",
		Props: mgod.SchemaFieldProps{
			Type:         reflect.String,
			Options:      fieldopt.SchemaFieldOptions{Required: false},
			Transformers: []transformer.Transformer{transformer.DateTransformerInstance},
		},
	})

	nestedModelWithSchemaOptsRootNode.Children = append(nestedModelWithSchemaOptsRootNode.Children, mgod.TreeNode{
		Path:    "$root.updatedAt",
		BSONKey: "updatedAt",
		Key:     "updatedAt",
		Props: mgod.SchemaFieldProps{
			Type:         reflect.String,
			Options:      fieldopt.SchemaFieldOptions{Required: false},
			Transformers: []transformer.Transformer{transformer.DateTransformerInstance},
		},
	})

	nestedModelWithSchemaOptsRootNode.Children = append(nestedModelWithSchemaOptsRootNode.Children, mgod.TreeNode{
		Path:    "$root.__v",
		BSONKey: "__v",
		Key:     "__v",
		Props: mgod.SchemaFieldProps{
			Type:         reflect.Int,
			Options:      fieldopt.SchemaFieldOptions{Required: false},
			Transformers: []transformer.Transformer{},
		},
	})

	nestedModelWithAllTypesAndSchemaOptsSchema := mgod.EntityModelSchema{Root: nestedModelWithSchemaOptsRootNode}

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
		actualSchema, err := mgod.BuildSchemaForModel(tc.Model, tc.SchemaOpts)

		s.Nil(err)
		if !reflect.DeepEqual(tc.Schema.Root, actualSchema.Root) {
			s.Fail("Schema mismatch", "Expected: %v, Got: %v", tc.Schema.Root, actualSchema.Root)
		}
		s.Equal(tc.ValidSchemaNodesCount, len(actualSchema.Nodes))
	}
}
