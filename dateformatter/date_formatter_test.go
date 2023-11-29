package dateformatter_test

import (
	"testing"
	"time"

	"github.com/Lyearn/mgod/dateformatter"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DateFormatterSuite struct {
	suite.Suite
	*require.Assertions
}

func TestDateFormatterSuite(t *testing.T) {
	s := new(DateFormatterSuite)
	suite.Run(t, s)
}

func (s *DateFormatterSuite) SetupTest() {
	s.Assertions = require.New(s.T())
}

func (s *DateFormatterSuite) TestDateFormatterISOString() {
	type TestCase struct {
		InputTime         time.Time
		ExpectedISOString string
	}

	loc, _ := time.LoadLocation("Asia/Kolkata")

	testCases := []TestCase{
		{
			InputTime:         time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			ExpectedISOString: "2023-01-01T00:00:00.000Z",
		},
		{
			InputTime:         time.Date(2023, 1, 1, 0, 0, 0, 1, time.UTC),
			ExpectedISOString: "2023-01-01T00:00:00.000Z",
		},
		{
			InputTime:         time.Date(2023, 1, 1, 0, 0, 0, 999999, time.UTC),
			ExpectedISOString: "2023-01-01T00:00:00.000Z",
		},
		{
			InputTime:         time.Date(2023, 1, 1, 0, 0, 0, 1000000, time.UTC),
			ExpectedISOString: "2023-01-01T00:00:00.001Z",
		},
		{
			InputTime:         time.Date(2023, 1, 1, 0, 0, 0, 1000001, time.UTC),
			ExpectedISOString: "2023-01-01T00:00:00.001Z",
		},
		{
			InputTime:         time.Date(2023, 1, 1, 0, 0, 0, 999999999, time.UTC),
			ExpectedISOString: "2023-01-01T00:00:00.999Z",
		},
		{
			InputTime:         time.Date(2023, 1, 1, 0, 0, 0, 0, loc),
			ExpectedISOString: "2022-12-31T18:30:00.000Z",
		},
		{
			InputTime:         time.Date(2023, 1, 1, 0, 0, 0, 1000000, loc),
			ExpectedISOString: "2022-12-31T18:30:00.001Z",
		},
	}

	for _, testCase := range testCases {
		df := dateformatter.New(testCase.InputTime)
		isoString, err := df.GetISOString()

		s.Nil(err)
		s.Equal(testCase.ExpectedISOString, isoString)
	}
}
