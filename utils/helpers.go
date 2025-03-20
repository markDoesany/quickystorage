package utils

import "time"

func FormatTimestamp(t time.Time) string {
	return t.Format("January 2, 2006 @ 3:04pm")
}
