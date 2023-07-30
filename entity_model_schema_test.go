package mongomodel_test

import (
	"reflect"
	"testing"

	"github.com/Lyearn/backend-universe/packages/store/mongomodel"
	"github.com/Lyearn/backend-universe/packages/store/mongomodel/schemaopt"
	"github.com/Lyearn/backend-universe/packages/store/mongomodel/transformer"
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
		Model  any
		Schema mongomodel.EntityModelSchema
	}

	type ActiveSession struct {
		SessionID   string `mgoType:"id"`
		LastLoginAt string `bson:"lastLoginAt" mgoType:"date"`
	}

	type Metadata struct {
		OnboardAt      string          `mgoType:"date"`
		TagIDs         []string        `bson:"tagIds" mgoType:"id"`
		ActiveSessions []ActiveSession `bson:"activeSessions" mgoID:"false" mgoDefault:"[]"`
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
		ID       string      `bson:"_id" mgoType:"id"`
		Name     *string     `bson:",omitempty"`
		Age      int         `mgoDefault:"18"`
		Metadata *Metadata   `bson:"meta"`
		Props    InlineProps `bson:",inline"`
		Height   float64
	}

	rootNode := mongomodel.GetDefaultSchemaTreeRootNode()
	rootNode.Children = []mongomodel.TreeNode{
		{
			Path:    "_id",
			BSONKey: "_id",
			Key:     "ID",
			Props: mongomodel.SchemaFieldProps{
				Type:         reflect.String,
				Transformers: []transformer.Transformer{transformer.IDTransformerInstance},
				Options: schemaopt.SchemaFieldOptions{
					Required: true,
				},
			},
		},
		{
			Path:    "name",
			BSONKey: "name",
			Key:     "Name",
			Props: mongomodel.SchemaFieldProps{
				Type:         reflect.String,
				IsPointer:    true,
				Transformers: []transformer.Transformer{},
				Options: schemaopt.SchemaFieldOptions{
					Required: false,
				},
			},
		},
		{
			Path:    "age",
			BSONKey: "age",
			Key:     "Age",
			Props: mongomodel.SchemaFieldProps{
				Type:         reflect.Int,
				Transformers: []transformer.Transformer{},
				Options: schemaopt.SchemaFieldOptions{
					Required: true,
					Default:  18,
				},
			},
		},
		{
			Path:    "meta",
			BSONKey: "meta",
			Key:     "Metadata",
			Props: mongomodel.SchemaFieldProps{
				Type:         reflect.Struct,
				IsPointer:    true,
				Transformers: []transformer.Transformer{},
				Options: schemaopt.SchemaFieldOptions{
					XID:      true,
					Required: true,
				},
			},
			Children: []mongomodel.TreeNode{
				{
					Path:    "meta._id",
					BSONKey: "_id",
					Key:     "XID",
					Props: mongomodel.SchemaFieldProps{
						Type:         reflect.String,
						Transformers: []transformer.Transformer{transformer.IDTransformerInstance},
						Options: schemaopt.SchemaFieldOptions{
							Required: true,
						},
					},
				},
				{
					Path:    "meta.onboardat",
					BSONKey: "onboardat",
					Key:     "OnboardAt",
					Props: mongomodel.SchemaFieldProps{
						Type:         reflect.String,
						Transformers: []transformer.Transformer{transformer.DateTransformerInstance},
						Options: schemaopt.SchemaFieldOptions{
							Required: true,
						},
					},
				},
				{
					Path:    "meta.tagIds",
					BSONKey: "tagIds",
					Key:     "TagIDs",
					Props: mongomodel.SchemaFieldProps{
						Type:         reflect.Slice,
						Transformers: []transformer.Transformer{},
						Options: schemaopt.SchemaFieldOptions{
							Required: true,
						},
					},
					Children: []mongomodel.TreeNode{
						{
							Path:    "meta.tagIds.$",
							BSONKey: "$",
							Key:     "$", // to identify slice element
							Props: mongomodel.SchemaFieldProps{
								Type:         reflect.String,
								Transformers: []transformer.Transformer{transformer.IDTransformerInstance},
							},
						},
					},
				},
				{
					Path:    "meta.activeSessions",
					BSONKey: "activeSessions",
					Key:     "ActiveSessions",
					Props: mongomodel.SchemaFieldProps{
						Type:         reflect.Slice,
						Transformers: []transformer.Transformer{},
						Options: schemaopt.SchemaFieldOptions{
							Required: true,
							Default:  bson.A{},
						},
					},
					Children: []mongomodel.TreeNode{
						{
							Path:    "meta.activeSessions.$",
							BSONKey: "$",
							Key:     "$",
							Props: mongomodel.SchemaFieldProps{
								Type:         reflect.Struct,
								Transformers: []transformer.Transformer{},
							},
							Children: []mongomodel.TreeNode{
								{
									Path:    "meta.activeSessions.$.sessionid",
									BSONKey: "sessionid",
									Key:     "SessionID",
									Props: mongomodel.SchemaFieldProps{
										Type:         reflect.String,
										Transformers: []transformer.Transformer{transformer.IDTransformerInstance},
										Options: schemaopt.SchemaFieldOptions{
											Required: true,
										},
									},
								},
								{
									Path:    "meta.activeSessions.$.lastLoginAt",
									BSONKey: "lastLoginAt",
									Key:     "LastLoginAt",
									Props: mongomodel.SchemaFieldProps{
										Type:         reflect.String,
										Transformers: []transformer.Transformer{transformer.DateTransformerInstance},
										Options: schemaopt.SchemaFieldOptions{
											Required: true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Path:    "inlineFloat",
			BSONKey: "inlineFloat",
			Key:     "InlineFloat",
			Props: mongomodel.SchemaFieldProps{
				Type:         reflect.Float64,
				Transformers: []transformer.Transformer{},
				Options: schemaopt.SchemaFieldOptions{
					Required: true,
				},
			},
		},
		{
			Path:    "inlineBool",
			BSONKey: "inlineBool",
			Key:     "InlineBool",
			Props: mongomodel.SchemaFieldProps{
				Type:         reflect.Bool,
				Transformers: []transformer.Transformer{},
				Options: schemaopt.SchemaFieldOptions{
					Required: true,
				},
			},
		},
		{
			Path:    "inlinestring",
			BSONKey: "inlinestring",
			Key:     "InlineString",
			Props: mongomodel.SchemaFieldProps{
				Type:         reflect.String,
				Transformers: []transformer.Transformer{},
				Options: schemaopt.SchemaFieldOptions{
					Required: true,
				},
			},
		},
		{
			Path:    "height",
			BSONKey: "height",
			Key:     "Height",
			Props: mongomodel.SchemaFieldProps{
				Type:         reflect.Float64,
				Transformers: []transformer.Transformer{},
				Options: schemaopt.SchemaFieldOptions{
					Required: true,
				},
			},
		},
	}
	nestedModelWithAllTypesSchema := mongomodel.EntityModelSchema{Root: rootNode}

	nestedModelWithAllTypesTestCase := TestCase{
		Model:  NestedModelWithAllTypes{},
		Schema: nestedModelWithAllTypesSchema,
	}

	testCases := []TestCase{
		nestedModelWithAllTypesTestCase,
	}

	for _, tc := range testCases {
		actualSchema, err := mongomodel.BuildSchemaForModel(tc.Model)
		s.Nil(err)
		if !reflect.DeepEqual(tc.Schema, *actualSchema) {
			s.Fail("Schema mismatch", "Expected: %v, Got: %v", tc.Schema, *actualSchema)
		}
	}
}
