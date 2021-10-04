package categories

import (
	"regexp"
	"strings"
	"sync"
)

var (
	re   *regexp.Regexp
	once sync.Once
)

func init() {
	getRegexp()
}

func getRegexp() *regexp.Regexp {
	once.Do(func() {
		re = regexp.MustCompile(`[^A-Za-zА-Яа-яёЁ0-9_\s]`)
	})
	return re
}

func PrepareSearchByName(query string) (args []string) {
	query = getRegexp().ReplaceAllString(query, " ")
	return strings.Fields(strings.ToLower(query))
}
