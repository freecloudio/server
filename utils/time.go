package utils

import "time"

func GetCurrentTime() time.Time {
	return time.Now().UTC()
}

func GetTimeIn(duration time.Duration) time.Time {
	return time.Now().Add(duration).UTC()
}
