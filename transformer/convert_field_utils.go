package transformer

import (
	"time"

	"github.com/Lyearn/backend-universe/packages/common/util/dateutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

func convertDateTimeToString(dateTimes ...primitive.DateTime) ([]string, error) {
	dates := []string{}

	for _, dateTime := range dateTimes {
		date, err := dateutil.New(dateTime.Time()).GetISOString()
		if err != nil {
			return nil, err
		}

		dates = append(dates, date)
	}

	return dates, nil
}

func convertObjectIDToString(objectIDs ...primitive.ObjectID) ([]string, error) {
	ids := []string{}

	for _, objectID := range objectIDs {
		id := objectID.Hex()

		ids = append(ids, id)
	}

	return ids, nil
}
