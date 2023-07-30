package mongomodel_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/Lyearn/backend-universe/packages/common/dateformatter"
	"github.com/Lyearn/backend-universe/packages/store/mongomodel"
	"github.com/Lyearn/backend-universe/packages/store/mongomodel/schemaopt"
	"github.com/Lyearn/backend-universe/packages/store/mongomodel/transformer"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BuildBSONDocSuite struct {
	suite.Suite
	*require.Assertions
}

func TestBuildBSONDocSuite(t *testing.T) {
	s := new(BuildBSONDocSuite)
	suite.Run(t, s)
}

func (s *BuildBSONDocSuite) SetupTest() {
	s.Assertions = require.New(s.T())
}

func (s *BuildBSONDocSuite) TestBuildBSONDoc() {
	type TestCase struct {
		TranslateTo       mongomodel.BSONDocTranslateToEnum
		EntityModelSchema mongomodel.EntityModelSchema
		InputDoc          bson.D
		ExpectedDoc       bson.D
	}

	id := primitive.NewObjectID()
	metaID := primitive.NewObjectID()

	onboardAt := primitive.NewDateTimeFromTime(time.Now())
	onboardAtStr, _ := dateformatter.New(onboardAt.Time()).GetISOString()

	tagIDs := []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID()}
	tagIDsStr := lo.Map(tagIDs, func(tagID primitive.ObjectID, _ int) string {
		return tagID.Hex()
	})

	activeSessions := []primitive.DateTime{primitive.NewDateTimeFromTime(time.Now()), primitive.NewDateTimeFromTime(time.Now())}
	activeSessionsStr := lo.Map(activeSessions, func(activeSession primitive.DateTime, _ int) string {
		requestTimestamp, _ := dateformatter.New(activeSession.Time()).GetISOString()
		return requestTimestamp
	})

	sessionIds := []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID()}
	sessionIdsStr := lo.Map(sessionIds, func(sessionId primitive.ObjectID, _ int) string {
		return sessionId.Hex()
	})

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
				Type:         reflect.TypeOf((*string)(nil)).Kind(),
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
					Path:    "meta.onboardAt",
					BSONKey: "onboardAt",
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
									Path:    "meta.activeSessions.$.sessionId",
									BSONKey: "sessionId",
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
	}

	nestedDocWithAllTypes := &TestCase{
		TranslateTo: mongomodel.BSONDocTranslateToEnumMongo,

		EntityModelSchema: mongomodel.EntityModelSchema{Root: rootNode},

		InputDoc: bson.D{
			{
				Key:   "_id",
				Value: id.Hex(),
			},
			{
				Key:   "name",
				Value: "user",
			},
			{
				Key: "meta",
				Value: bson.D{
					{
						Key:   "_id",
						Value: metaID.Hex(),
					},
					{
						Key:   "onboardAt",
						Value: onboardAtStr,
					},
					{
						Key:   "tagIds",
						Value: bson.A{tagIDsStr[0], tagIDsStr[1]},
					},
					{
						Key: "activeSessions",
						Value: bson.A{
							bson.D{
								{
									Key:   "sessionId",
									Value: sessionIdsStr[0],
								},
								{
									Key:   "lastLoginAt",
									Value: activeSessionsStr[0],
								},
							},
							bson.D{
								{
									Key:   "sessionId",
									Value: sessionIdsStr[1],
								},
								{
									Key:   "lastLoginAt",
									Value: activeSessionsStr[1],
								},
							},
						},
					},
				},
			},
		},

		ExpectedDoc: bson.D{
			{
				Key:   "_id",
				Value: id,
			},
			{
				Key:   "name",
				Value: "user",
			},
			{
				Key:   "age",
				Value: 18,
			},
			{
				Key: "meta",
				Value: bson.D{
					{
						Key:   "_id",
						Value: metaID,
					},
					{
						Key:   "onboardAt",
						Value: onboardAt,
					},
					{
						Key:   "tagIds",
						Value: bson.A{tagIDs[0], tagIDs[1]},
					},
					{
						Key: "activeSessions",
						Value: bson.A{
							bson.D{
								{
									Key:   "sessionId",
									Value: sessionIds[0],
								},
								{
									Key:   "lastLoginAt",
									Value: activeSessions[0],
								},
							},
							bson.D{
								{
									Key:   "sessionId",
									Value: sessionIds[1],
								},
								{
									Key:   "lastLoginAt",
									Value: activeSessions[1],
								},
							},
						},
					},
				},
			},
		},
	}

	nestedDocCheckForDefaultValues := &TestCase{
		TranslateTo: mongomodel.BSONDocTranslateToEnumMongo,

		EntityModelSchema: mongomodel.EntityModelSchema{Root: rootNode},

		InputDoc: bson.D{
			{
				Key:   "_id",
				Value: id.Hex(),
			},
			{
				Key:   "name",
				Value: "user",
			},
			{
				Key: "meta",
				Value: bson.D{
					{
						Key:   "_id",
						Value: metaID.Hex(),
					},
					{
						Key:   "onboardAt",
						Value: onboardAtStr,
					},
					{
						Key:   "tagIds",
						Value: bson.A{tagIDsStr[0], tagIDsStr[1]},
					},
				},
			},
		},

		ExpectedDoc: bson.D{
			{
				Key:   "_id",
				Value: id,
			},
			{
				Key:   "name",
				Value: "user",
			},
			{
				Key:   "age",
				Value: 18,
			},
			{
				Key: "meta",
				Value: bson.D{
					{
						Key:   "_id",
						Value: metaID,
					},
					{
						Key:   "onboardAt",
						Value: onboardAt,
					},
					{
						Key:   "tagIds",
						Value: bson.A{tagIDs[0], tagIDs[1]},
					},
					{
						Key:   "activeSessions",
						Value: bson.A{},
					},
				},
			},
		},
	}

	testCases := []*TestCase{
		nestedDocWithAllTypes,
		nestedDocCheckForDefaultValues,
	}

	for _, testCase := range testCases {
		doc := testCase.InputDoc
		err := mongomodel.BuildBSONDoc(context.TODO(), &doc, &testCase.EntityModelSchema, testCase.TranslateTo)

		s.Nil(err)
		s.Equal(testCase.ExpectedDoc, doc)
	}
}

func (s *BuildBSONDocSuite) TestBuildBSONDocWithoutID() {
	type TestCase struct {
		TranslateTo       mongomodel.BSONDocTranslateToEnum
		EntityModelSchema mongomodel.EntityModelSchema
		InputDoc          bson.D
	}

	onboardAt := primitive.NewDateTimeFromTime(time.Now())
	onboardAtStr, _ := dateformatter.New(onboardAt.Time()).GetISOString()

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
				Type:         reflect.TypeOf((*string)(nil)).Kind(),
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
					Path:    "meta.onboardAt",
					BSONKey: "onboardAt",
					Key:     "OnboardAt",
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
	}

	docWithoutIDCase := &TestCase{
		TranslateTo: mongomodel.BSONDocTranslateToEnumMongo,

		EntityModelSchema: mongomodel.EntityModelSchema{Root: rootNode},

		InputDoc: bson.D{
			{
				Key:   "name",
				Value: "user",
			},
			{
				Key: "meta",
				Value: bson.D{
					{
						Key:   "onboardAt",
						Value: onboardAtStr,
					},
				},
			},
		},
	}

	doc := docWithoutIDCase.InputDoc
	err := mongomodel.BuildBSONDoc(context.TODO(), &doc, &docWithoutIDCase.EntityModelSchema, docWithoutIDCase.TranslateTo)

	s.Nil(err)
	s.True(doc[0].Key == "_id")
	s.True(doc[0].Value.(primitive.ObjectID).Hex() != "")
	s.True(doc[3].Key == "meta")
	s.True(doc[3].Value.(primitive.D)[0].Key == "_id")
	s.True(doc[3].Value.(primitive.D)[0].Value.(primitive.ObjectID).Hex() != "")
}
