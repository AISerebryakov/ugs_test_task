package repository

import "strconv"

func pgSqlArguments(n int) []string {
	args := make([]string, n)
	for i := 1; i <= n; i++ {
		args[i-1] = "$" + strconv.Itoa(i)
	}
	return args
}
