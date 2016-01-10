package utils

// ContainsString is a helper function to check whether
// an array of strings contains a specific string
func ContainsString(elements []string, element string) bool {
	for _, v := range elements {
		if element == v {
			return true
		}
	}
	return false
}
