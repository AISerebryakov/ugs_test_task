package companyrepos

func categoriesToLtreeArgs(categories []string) (args string) {
	for i, c := range categories {
		if len(c) == 0 {
			continue
		}
		if i < len(categories)-1 {
			args = args + c + "*@|"
			continue
		}
		args = args + c + "*@"
	}
	return args
}

func categoryNameToLtreeArg(name string) (arg string) {
	return name + "*@"
}
