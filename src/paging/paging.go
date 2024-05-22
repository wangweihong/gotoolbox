package paging

import "math"

func Index(length int, page, size int) (sIndex, eIndex int) {
	if length < 0 {
		length = 0
	}

	if page < 0 {
		page = 0
	}

	if size < 0 {
		size = 0
	}

	if page == 0 && size == 0 {
		sIndex = 0
		eIndex = length
		return
	}

	// page从0开始
	sIndex = page * size
	eIndex = int(math.Min(float64(sIndex+size), float64(length)))
	if sIndex > eIndex {
		sIndex = 0
		eIndex = 0
		return
	}

	return
}
