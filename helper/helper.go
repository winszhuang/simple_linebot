package helper

func Include(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}
