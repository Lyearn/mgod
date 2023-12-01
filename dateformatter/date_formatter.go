package dateformatter

import (
	"errors"
	"strings"
	"time"
)

// DateFormatter provides utilities to get date time in different formats.
// NOTE: DateFormatter processes the provided date time string in UTC format.
type DateFormatter struct {
	t time.Time
}

func New(t time.Time) *DateFormatter {
	return &DateFormatter{t: t.UTC()}
}

func (d *DateFormatter) GetISOString() (string, error) {
	// ISO 8601 format: YYYY-MM-DDTHH:mm:ss.sssZ

	// using RFC3339Nano format to capture milliseconds
	formattedTime := d.t.Format(time.RFC3339Nano)

	// add zeros to milliseconds if not present
	formattedTime = formatMilliSecondsString(formattedTime)

	// Go time format: 2023-01-01T10:10:10.999+5:30
	isoStringLen := 23

	// Go time format with zero millis value: 2023-01-01T10:10:10Z
	isoStringWithoutMillisLen := 19

	if len(formattedTime) >= isoStringLen {
		// milliseconds present in the time string
		formattedTime = formattedTime[:isoStringLen] + "Z"

		return formattedTime, nil
	} else if len(formattedTime) >= isoStringWithoutMillisLen {
		// milliseconds not present in the time string. hence, adding its zero value.
		// milliseconds are ignored for its zero value i.e. 000
		formattedTime = formattedTime[:isoStringWithoutMillisLen] + ".000Z"

		return formattedTime, nil
	}

	return "", errors.New("invalid time format")
}

// formatMilliSecondsString adds zeros to milliseconds if not present.
func formatMilliSecondsString(dateTimeStr string) string {
	parts := strings.Split(dateTimeStr, ".")
	partsAfterMillisSplit := 2
	if len(parts) != partsAfterMillisSplit {
		return dateTimeStr
	}

	millisStr := strings.Split(parts[1], "Z")
	expectedMillisLen := 3
	zerosToAppend := expectedMillisLen - len(millisStr[0])

	zerosStr := ""
	for i := 0; i < zerosToAppend; i++ {
		zerosStr += "0"
	}

	return parts[0] + "." + millisStr[0] + zerosStr + "Z"
}
