package timefuncs

import "time"

func SecondToHHMMSS(timeInt int64) string {
	return time.Unix(timeInt, 0).UTC().Format("15:04:05")
}
