package util

func Paginate[T any](data []T, pageIndex int, pageSize int) []T {
	if pageIndex <= 0 || pageSize <= 0 {
		return []T{}
	}
	start := (pageIndex - 1) * pageSize
	end := pageIndex * pageSize
	if start >= len(data) {
		return []T{}
	}
	if end > len(data) {
		end = len(data)
	}
	return data[start:end]
}
