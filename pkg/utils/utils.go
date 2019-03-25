package utils

func IndexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k //index
		}
	}
	return -1 //not found.
}
