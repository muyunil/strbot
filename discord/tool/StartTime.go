package tool

import (
	"strings"
	"time"
)

func StartTime(t *time.Time) string {
	Time := time.Now()
	str := Time.Sub(*t).String()
	if strings.Contains(str, "m") {
		index := strings.Index(str, "m")
		return str[:index+1]
	} else {
		return "0m"
	}
}
