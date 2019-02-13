package utils

import (
	"regexp"
	"time"
)
var (
	patternLegalNodeName = regexp.MustCompile(`^[0-9a-zA-Z_.\[\]:-]+$`)
)

func ValidateNodeName(n string) bool {
	return patternLegalNodeName.MatchString(n)
}

func CheckTime(t int64, duration int) bool {
	nowTime := time.Now().UTC().UnixNano() / 1000000 		//毫秒
	if (nowTime -t ) > int64( duration * 2 * 1000) {
		return false
	}
	return true
}