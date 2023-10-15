package lib

import "time"

const (
	layoutISO = "2006-01-02"
)

func ConvertToDate(timestamp time.Time) string {
	return timestamp.Format(layoutISO)
}

func ConvertToEpoch(timestamp time.Time) int64 {
	var epochTimestamp int64
	if !timestamp.IsZero() {
		epochTimestamp = timestamp.Unix()
	}

	return epochTimestamp
}
