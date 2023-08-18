package mongomodel_test

import (
	"context"
	"testing"
	"time"

	"github.com/Lyearn/backend-universe/packages/common/util/dateutil"
	"github.com/Lyearn/backend-universe/packages/common/util/typeutil"
	"github.com/Lyearn/backend-universe/packages/store/acl/model"
	"github.com/Lyearn/backend-universe/packages/store/mongomodel"
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
	onboardAtStr, _ := dateutil.New(onboardAt.Time()).GetISOString()

	tagIDs := []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID()}
	tagIDsStr := lo.Map(tagIDs, func(tagID primitive.ObjectID, _ int) string {
		return tagID.Hex()
	})

	activeSessions := []primitive.DateTime{primitive.NewDateTimeFromTime(time.Now()), primitive.NewDateTimeFromTime(time.Now())}
	activeSessionsStr := lo.Map(activeSessions, func(activeSession primitive.DateTime, _ int) string {
		requestTimestamp, _ := dateutil.New(activeSession.Time()).GetISOString()
		return requestTimestamp
	})

	sessionIds := []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID()}
	sessionIdsStr := lo.Map(sessionIds, func(sessionId primitive.ObjectID, _ int) string {
		return sessionId.Hex()
	})

	createdAt := primitive.NewDateTimeFromTime(time.Now())
	updatedAt := primitive.NewDateTimeFromTime(time.Now())

	type ActiveSession struct {
		SessionID   string `bson:"sessionId" mgoType:"id"`
		LastLoginAt string `bson:"lastLoginAt" mgoType:"date"`
	}

	type Metadata struct {
		OnboardAt      string          `bson:"onboardAt" mgoType:"date"`
		TagIDs         []string        `bson:"tagIds" mgoType:"id"`
		ActiveSessions []ActiveSession `bson:"activeSessions" mgoID:"false" mgoDefault:"[]"`
	}

	type NestedModelWithAllTypes struct {
		ID       string    `bson:"_id" mgoType:"id"`
		Name     *string   `bson:",omitempty"`
		Age      *int      `bson:",omitempty" mgoDefault:"18"`
		Metadata *Metadata `bson:"meta"`
	}

	actualSchema, _ := mongomodel.BuildSchemaForModel(
		NestedModelWithAllTypes{},
		model.SchemaOptions{VersionKey: typeutil.GetPointerForConstType(false)},
	)

	//nolint:stylecheck,revive // ignore linter for test case
	nestedDocWithAllTypes_toMongo := &TestCase{
		TranslateTo: mongomodel.BSONDocTranslateToEnumMongo,

		EntityModelSchema: *actualSchema,

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
				Key:   "age",
				Value: 18,
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

	//nolint:stylecheck,revive // ignore linter for test case
	nestedDocWithAllTypes_toEntityModel := &TestCase{
		TranslateTo: mongomodel.BSONDocTranslateToEnumEntityModel,

		EntityModelSchema: *actualSchema,

		InputDoc: bson.D{
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

		ExpectedDoc: bson.D{
			{
				Key:   "_id",
				Value: id.Hex(),
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
	}

	//nolint:stylecheck,revive // ignore linter for test case
	nestedDocCheckForDefaultValues_toMongo := &TestCase{
		TranslateTo: mongomodel.BSONDocTranslateToEnumMongo,

		EntityModelSchema: *actualSchema,

		InputDoc: bson.D{
			{
				Key:   "_id",
				Value: id.Hex(),
			},
			{
				Key: "meta",
				Value: bson.D{
					{
						Key:   "_id",
						Value: metaID.Hex(),
					},
					{
						Key:   "tagIds",
						Value: bson.A{tagIDsStr[0], tagIDsStr[1]},
					},
					{
						Key:   "onboardAt",
						Value: onboardAtStr,
					},
				},
			},
			{
				Key:   "name",
				Value: "user",
			},
		},

		ExpectedDoc: bson.D{
			{
				Key:   "_id",
				Value: id,
			},
			{
				Key: "meta",
				Value: bson.D{
					{
						Key:   "_id",
						Value: metaID,
					},
					{
						Key:   "tagIds",
						Value: bson.A{tagIDs[0], tagIDs[1]},
					},
					{
						Key:   "onboardAt",
						Value: onboardAt,
					},
					{
						Key:   "activeSessions",
						Value: bson.A{},
					},
				},
			},
			{
				Key:   "name",
				Value: "user",
			},
			{
				Key:   "age",
				Value: 18,
			},
		},
	}

	//nolint:stylecheck,revive // ignore linter for test case
	nestedDocCheckForDefaultValues_toEntityModel := &TestCase{
		TranslateTo: mongomodel.BSONDocTranslateToEnumEntityModel,

		EntityModelSchema: *actualSchema,

		InputDoc: bson.D{
			{
				Key:   "_id",
				Value: id,
			},
			{
				Key: "meta",
				Value: bson.D{
					{
						Key:   "_id",
						Value: metaID,
					},
					{
						Key:   "tagIds",
						Value: bson.A{tagIDs[0], tagIDs[1]},
					},
					{
						Key:   "onboardAt",
						Value: onboardAt,
					},
				},
			},
		},

		ExpectedDoc: bson.D{
			{
				Key:   "_id",
				Value: id.Hex(),
			},
			{
				Key: "meta",
				Value: bson.D{
					{
						Key:   "_id",
						Value: metaID.Hex(),
					},
					{
						Key:   "tagIds",
						Value: bson.A{tagIDsStr[0], tagIDsStr[1]},
					},
					{
						Key:   "onboardAt",
						Value: onboardAtStr,
					},
					{
						Key:   "activeSessions",
						Value: bson.A{},
					},
				},
			},
			{
				Key:   "age",
				Value: 18,
			},
		},
	}

	schemaForNestedDocWithSchemaOpts, _ := mongomodel.BuildSchemaForModel(
		NestedModelWithAllTypes{},
		model.SchemaOptions{
			Timestamps: true,
			// VersionKey enabled by default
		},
	)

	//nolint:stylecheck,revive // ignore linter for test case
	nestedDocWithAllTypesAndSchemaOpts_toMongo := &TestCase{
		TranslateTo: mongomodel.BSONDocTranslateToEnumMongo,

		EntityModelSchema: *schemaForNestedDocWithSchemaOpts,

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
				Key:   "age",
				Value: 18,
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

	//nolint:stylecheck,revive // ignore linter for test case
	nestedDocWithAllTypesAndSchemaOpts_toEntityModel := &TestCase{
		TranslateTo: mongomodel.BSONDocTranslateToEnumEntityModel,

		EntityModelSchema: *schemaForNestedDocWithSchemaOpts,

		InputDoc: bson.D{
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
			{
				Key:   "createdAt",
				Value: createdAt,
			},
			{
				Key:   "updatedAt",
				Value: updatedAt,
			},
			{
				Key:   "__v",
				Value: 0,
			},
		},

		ExpectedDoc: bson.D{
			{
				Key:   "_id",
				Value: id.Hex(),
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
			{
				Key:   "createdAt",
				Value: createdAt,
			},
			{
				Key:   "updatedAt",
				Value: updatedAt,
			},
			{
				Key:   "__v",
				Value: 0,
			},
		},
	}

	testCases := []*TestCase{
		nestedDocWithAllTypes_toMongo,
		nestedDocWithAllTypes_toEntityModel,
		nestedDocCheckForDefaultValues_toMongo,
		nestedDocCheckForDefaultValues_toEntityModel,
		nestedDocWithAllTypesAndSchemaOpts_toMongo,
		nestedDocWithAllTypesAndSchemaOpts_toEntityModel,
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
	onboardAtStr, _ := dateutil.New(onboardAt.Time()).GetISOString()

	type Metadata struct {
		OnboardAt string `bson:"onboardAt" mgoType:"date"`
	}

	type NestedModel struct {
		ID       string    `bson:"_id" mgoType:"id"`
		Name     *string   `bson:",omitempty"`
		Age      int       `mgoDefault:"18"`
		Metadata *Metadata `bson:"meta"`
	}

	actualSchema, _ := mongomodel.BuildSchemaForModel(NestedModel{}, model.SchemaOptions{})

	docWithoutIDCase := &TestCase{
		TranslateTo: mongomodel.BSONDocTranslateToEnumMongo,

		EntityModelSchema: *actualSchema,

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
	s.True(doc[2].Key == "_id")
	s.True(doc[2].Value.(primitive.ObjectID).Hex() != "")
	s.True(doc[1].Key == "meta")
	s.True(doc[1].Value.(primitive.D)[1].Key == "_id")
	s.True(doc[1].Value.(primitive.D)[1].Value.(primitive.ObjectID).Hex() != "")
}
