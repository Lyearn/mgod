package bsondoc_test

import (
	"context"
	"testing"
	"time"

	"github.com/Lyearn/mgod/bsondoc"
	"github.com/Lyearn/mgod/dateformatter"
	"github.com/Lyearn/mgod/schema"
	"github.com/Lyearn/mgod/schema/schemaopt"
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
		TranslateTo       bsondoc.TranslateToEnum
		EntityModelSchema schema.EntityModelSchema
		InputDoc          bson.D
		ExpectedDoc       bson.D
	}

	id := primitive.NewObjectID()
	metaID := primitive.NewObjectID()

	joinedOn := primitive.NewDateTimeFromTime(time.Now())
	joinedOnStr, _ := dateformatter.New(joinedOn.Time()).GetISOString()

	teamIDs := []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID()}
	teamIDsStr := lo.Map(teamIDs, func(teamID primitive.ObjectID, _ int) string {
		return teamID.Hex()
	})

	projectsCompletedAt := []primitive.DateTime{primitive.NewDateTimeFromTime(time.Now()), primitive.NewDateTimeFromTime(time.Now())}
	projectsCompletedAtStr := lo.Map(projectsCompletedAt, func(projectCompletedAt primitive.DateTime, _ int) string {
		requestTimestamp, _ := dateformatter.New(projectCompletedAt.Time()).GetISOString()
		return requestTimestamp
	})

	projectIDs := []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID()}
	projectIDsStr := lo.Map(projectIDs, func(projectID primitive.ObjectID, _ int) string {
		return projectID.Hex()
	})

	createdAt := primitive.NewDateTimeFromTime(time.Now())
	createdAtTimestamp, _ := dateformatter.New(createdAt.Time()).GetISOString()
	updatedAt := primitive.NewDateTimeFromTime(time.Now())
	updatedAtTimestamp, _ := dateformatter.New(updatedAt.Time()).GetISOString()

	type UserProject struct {
		ProjectID   string `bson:"projectId" mgoType:"id"`
		CompletedAt string `bson:"completedAt" mgoType:"date"`
	}

	type Metadata struct {
		JoinedOn string        `bson:"joinedOn" mgoType:"date"`
		TeamIDs  []string      `bson:"teamIds" mgoType:"id"`
		Projects []UserProject `bson:"projects" mgoID:"false" mgoDefault:"[]"`
	}

	type NestedModelWithAllTypes struct {
		ID       string    `bson:"_id" mgoType:"id"`
		Name     *string   `bson:",omitempty"`
		Age      *int      `bson:",omitempty" mgoDefault:"18"`
		Metadata *Metadata `bson:"meta"`
	}

	versionKeyEnabled := false
	actualSchema, _ := schema.BuildSchemaForModel(
		NestedModelWithAllTypes{},
		schemaopt.SchemaOptions{VersionKey: &versionKeyEnabled},
	)

	//nolint:stylecheck,revive // ignore linter for test case
	nestedDocWithAllTypes_toMongo := &TestCase{
		TranslateTo: bsondoc.TranslateToEnumMongo,

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
						Key:   "joinedOn",
						Value: joinedOnStr,
					},
					{
						Key:   "teamIds",
						Value: bson.A{teamIDsStr[0], teamIDsStr[1]},
					},
					{
						Key: "projects",
						Value: bson.A{
							bson.D{
								{
									Key:   "projectId",
									Value: projectIDsStr[0],
								},
								{
									Key:   "completedAt",
									Value: projectsCompletedAtStr[0],
								},
							},
							bson.D{
								{
									Key:   "projectId",
									Value: projectIDsStr[1],
								},
								{
									Key:   "completedAt",
									Value: projectsCompletedAtStr[1],
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
						Key:   "joinedOn",
						Value: joinedOn,
					},
					{
						Key:   "teamIds",
						Value: bson.A{teamIDs[0], teamIDs[1]},
					},
					{
						Key: "projects",
						Value: bson.A{
							bson.D{
								{
									Key:   "projectId",
									Value: projectIDs[0],
								},
								{
									Key:   "completedAt",
									Value: projectsCompletedAt[0],
								},
							},
							bson.D{
								{
									Key:   "projectId",
									Value: projectIDs[1],
								},
								{
									Key:   "completedAt",
									Value: projectsCompletedAt[1],
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
		TranslateTo: bsondoc.TranslateToEnumEntityModel,

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
						Key:   "joinedOn",
						Value: joinedOn,
					},
					{
						Key:   "teamIds",
						Value: bson.A{teamIDs[0], teamIDs[1]},
					},
					{
						Key: "projects",
						Value: bson.A{
							bson.D{
								{
									Key:   "projectId",
									Value: projectIDs[0],
								},
								{
									Key:   "completedAt",
									Value: projectsCompletedAt[0],
								},
							},
							bson.D{
								{
									Key:   "projectId",
									Value: projectIDs[1],
								},
								{
									Key:   "completedAt",
									Value: projectsCompletedAt[1],
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
						Key:   "joinedOn",
						Value: joinedOnStr,
					},
					{
						Key:   "teamIds",
						Value: bson.A{teamIDsStr[0], teamIDsStr[1]},
					},
					{
						Key: "projects",
						Value: bson.A{
							bson.D{
								{
									Key:   "projectId",
									Value: projectIDsStr[0],
								},
								{
									Key:   "completedAt",
									Value: projectsCompletedAtStr[0],
								},
							},
							bson.D{
								{
									Key:   "projectId",
									Value: projectIDsStr[1],
								},
								{
									Key:   "completedAt",
									Value: projectsCompletedAtStr[1],
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
		TranslateTo: bsondoc.TranslateToEnumMongo,

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
						Key:   "teamIds",
						Value: bson.A{teamIDsStr[0], teamIDsStr[1]},
					},
					{
						Key:   "joinedOn",
						Value: joinedOnStr,
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
						Key:   "teamIds",
						Value: bson.A{teamIDs[0], teamIDs[1]},
					},
					{
						Key:   "joinedOn",
						Value: joinedOn,
					},
					{
						Key:   "projects",
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
		TranslateTo: bsondoc.TranslateToEnumEntityModel,

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
						Key:   "teamIds",
						Value: bson.A{teamIDs[0], teamIDs[1]},
					},
					{
						Key:   "joinedOn",
						Value: joinedOn,
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
						Key:   "teamIds",
						Value: bson.A{teamIDsStr[0], teamIDsStr[1]},
					},
					{
						Key:   "joinedOn",
						Value: joinedOnStr,
					},
					{
						Key:   "projects",
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

	schemaForNestedDocWithSchemaOpts, _ := schema.BuildSchemaForModel(
		NestedModelWithAllTypes{},
		schemaopt.SchemaOptions{
			Timestamps: true,
			// VersionKey enabled by default
		},
	)

	//nolint:stylecheck,revive // ignore linter for test case
	nestedDocWithAllTypesAndSchemaOpts_toMongo := &TestCase{
		TranslateTo: bsondoc.TranslateToEnumMongo,

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
						Key:   "joinedOn",
						Value: joinedOnStr,
					},
					{
						Key:   "teamIds",
						Value: bson.A{teamIDsStr[0], teamIDsStr[1]},
					},
					{
						Key: "projects",
						Value: bson.A{
							bson.D{
								{
									Key:   "projectId",
									Value: projectIDsStr[0],
								},
								{
									Key:   "completedAt",
									Value: projectsCompletedAtStr[0],
								},
							},
							bson.D{
								{
									Key:   "projectId",
									Value: projectIDsStr[1],
								},
								{
									Key:   "completedAt",
									Value: projectsCompletedAtStr[1],
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
						Key:   "joinedOn",
						Value: joinedOn,
					},
					{
						Key:   "teamIds",
						Value: bson.A{teamIDs[0], teamIDs[1]},
					},
					{
						Key: "projects",
						Value: bson.A{
							bson.D{
								{
									Key:   "projectId",
									Value: projectIDs[0],
								},
								{
									Key:   "completedAt",
									Value: projectsCompletedAt[0],
								},
							},
							bson.D{
								{
									Key:   "projectId",
									Value: projectIDs[1],
								},
								{
									Key:   "completedAt",
									Value: projectsCompletedAt[1],
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
		TranslateTo: bsondoc.TranslateToEnumEntityModel,

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
						Key:   "joinedOn",
						Value: joinedOn,
					},
					{
						Key:   "teamIds",
						Value: bson.A{teamIDs[0], teamIDs[1]},
					},
					{
						Key: "projects",
						Value: bson.A{
							bson.D{
								{
									Key:   "projectId",
									Value: projectIDs[0],
								},
								{
									Key:   "completedAt",
									Value: projectsCompletedAt[0],
								},
							},
							bson.D{
								{
									Key:   "projectId",
									Value: projectIDs[1],
								},
								{
									Key:   "completedAt",
									Value: projectsCompletedAt[1],
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
						Key:   "joinedOn",
						Value: joinedOnStr,
					},
					{
						Key:   "teamIds",
						Value: bson.A{teamIDsStr[0], teamIDsStr[1]},
					},
					{
						Key: "projects",
						Value: bson.A{
							bson.D{
								{
									Key:   "projectId",
									Value: projectIDsStr[0],
								},
								{
									Key:   "completedAt",
									Value: projectsCompletedAtStr[0],
								},
							},
							bson.D{
								{
									Key:   "projectId",
									Value: projectIDsStr[1],
								},
								{
									Key:   "completedAt",
									Value: projectsCompletedAtStr[1],
								},
							},
						},
					},
				},
			},
			{
				Key:   "createdAt",
				Value: createdAtTimestamp,
			},
			{
				Key:   "updatedAt",
				Value: updatedAtTimestamp,
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
		err := bsondoc.Build(context.TODO(), &doc, &testCase.EntityModelSchema, testCase.TranslateTo)

		s.Nil(err)
		s.Equal(testCase.ExpectedDoc, doc)
	}
}

func (s *BuildBSONDocSuite) TestBuildBSONDocWithoutID() {
	type TestCase struct {
		TranslateTo       bsondoc.TranslateToEnum
		EntityModelSchema schema.EntityModelSchema
		InputDoc          bson.D
	}

	joinedOn := primitive.NewDateTimeFromTime(time.Now())
	joinedOnStr, _ := dateformatter.New(joinedOn.Time()).GetISOString()

	type Metadata struct {
		JoinedOn string `bson:"joinedOn" mgoType:"date"`
	}

	type NestedModel struct {
		ID       string    `bson:"_id" mgoType:"id"`
		Name     *string   `bson:",omitempty"`
		Age      int       `mgoDefault:"18"`
		Metadata *Metadata `bson:"meta"`
	}

	actualSchema, _ := schema.BuildSchemaForModel(NestedModel{}, schemaopt.SchemaOptions{})

	docWithoutIDCase := &TestCase{
		TranslateTo: bsondoc.TranslateToEnumMongo,

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
						Key:   "joinedOn",
						Value: joinedOnStr,
					},
				},
			},
		},
	}

	doc := docWithoutIDCase.InputDoc
	err := bsondoc.Build(context.TODO(), &doc, &docWithoutIDCase.EntityModelSchema, docWithoutIDCase.TranslateTo)

	s.Nil(err)
	s.True(doc[2].Key == "_id")
	s.True(doc[2].Value.(primitive.ObjectID).Hex() != "")
	s.True(doc[1].Key == "meta")
	s.True(doc[1].Value.(primitive.D)[1].Key == "_id")
	s.True(doc[1].Value.(primitive.D)[1].Value.(primitive.ObjectID).Hex() != "")
}
