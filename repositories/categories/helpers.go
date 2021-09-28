package categories

import (
	"regexp"
	"strings"
	"sync"
)

var (
	rg   *regexp.Regexp
	once sync.Once
)

func getRegexp() *regexp.Regexp {
	once.Do(func() {
		rg = regexp.MustCompile(`[^A-Za-zА-Яа-яёЁ0-9_\s]`)
	})
	return rg
}

func PrepareSearchByName(query string) (args string) {
	query = getRegexp().ReplaceAllString(query, "")
	return NamesToLtreeArgs(strings.Fields(query))
}

func NamesToLtreeArgs(names []string) (args string) {
	for i, c := range names {
		if len(c) == 0 {
			continue
		}
		if i < len(names)-1 {
			args = args + c + "*@|"
			continue
		}
		args = args + c + "*@"
	}
	return args
}
