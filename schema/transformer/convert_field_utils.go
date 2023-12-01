package transformer

import (
	"time"

	"github.com/Lyearn/mgod/dateformatter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ConvertStringToDateTime converts the provided string dates to primitive.DateTime.
func convertStringToDateTime(dates ...string) ([]primitive.DateTime, error) {
	dateTimes := []primitive.DateTime{}

	for _, date := range dates {
		goTime, err := time.Parse(time.RFC3339Nano, date)
		if err != nil {
			return nil, err
		}

		dateTime := primitive.NewDateTimeFromTime(goTime)

		dateTimes = append(dateTimes, dateTime)
	}

	return dateTimes, nil
}

// ConvertStringToObjectID converts the provided string ids to primitive.ObjectID.
func convertStringToObjectID(ids ...string) ([]primitive.ObjectID, error) {
	objectIDS := []primitive.ObjectID{}

	for _, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}

		objectIDS = append(objectIDS, objectID)
	}

	return objectIDS, nil
}

// ConvertDateTimeToString converts the provided primitive.DateTime to string.
func convertDateTimeToString(dateTimes ...primitive.DateTime) ([]string, error) {
	dates := []string{}

	for _, dateTime := range dateTimes {
		date, err := dateformatter.New(dateTime.Time()).GetISOString()
		if err != nil {
			return nil, err
		}

		dates = append(dates, date)
	}

	return dates, nil
}

// ConvertObjectIDToString converts the provided primitive.ObjectID to string.
//
//nolint:unparam // error field added to keep the signature of convert field functions consistent
func convertObjectIDToString(objectIDs ...primitive.ObjectID) ([]string, error) {
	ids := []string{}

	for _, objectID := range objectIDs {
		id := objectID.Hex()

		ids = append(ids, id)
	}

	return ids, nil
}
