package helper

import (
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/schema"
	"github.com/oklog/ulid/v2"
)

func DecodeQuery(query string, target interface{}) interface{} {
	values, _ := url.ParseQuery(query)

	encoder := schema.NewDecoder()
	encoder.IgnoreUnknownKeys(true)

	encoder.Decode(target, values)

	return target
}

func GenerateId() string {
	//return uuid.New().String()[0:8]
	return ulid.Make().String()
}

func DateToUnixString(date string) string {
	datetime, _ := time.Parse("2006-01-02", date)
	timestamp := datetime.Unix()

	return strconv.FormatInt(timestamp, 10)
}

func UnixToDateString(timestamp string) string {
	i, _ := strconv.ParseInt(timestamp, 10, 64)
	tm := time.Unix(i, 0)
	return tm.Format("2006-01-02")
}

func UnixToDate(timestamp string) time.Time {
	i, _ := strconv.ParseInt(timestamp, 10, 64)
	return time.Unix(i, 0)
}
