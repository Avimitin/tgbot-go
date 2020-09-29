package timer

import (
	"time"
)

func UnixToString(unixTime int64) string {
	return time.Unix(unixTime, 0).String()
}