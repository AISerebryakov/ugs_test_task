package config

import "strings"

func parseArrayStr(str string) []string {
	str = strings.ReplaceAll(str, " ", ",")
	str = strings.ReplaceAll(str, ",,", ",")
	return strings.Split(str, ",")
}
