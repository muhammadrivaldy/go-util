package util

// Pagination is a function for calculate of limit & offset
func Pagination(limit, offset int) (limits, offsets int) {
	if offset == 0 {
		return 0, 0
	}

	offset--
	offset *= limit

	return limit, offset
}
