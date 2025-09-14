package utils

import "time"

type ShortDate time.Time

func (s ShortDate) MarshalJSON() ([]byte, error) {
	// format time into Django date "YYYY_MM_DD"
	formattedDate := time.Time(s).Format("2006_01_02")
	// return as JSON string
	return []byte(`"` + formattedDate + `"`), nil
}
