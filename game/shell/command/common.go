package command

func filterArgs(args []string) []string {
	var result []string
	for _, arg := range args {
		if arg != "" {
			result = append(result, arg)
		}
	}

	return result
}
