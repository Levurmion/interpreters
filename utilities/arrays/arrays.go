package arrays

func Map[A any, O any](array []A, mapper func(A) O) []O {
	result := make([]O, len(array))
	for idx, val := range array {
		result[idx] = mapper(val)
	}
	return result
}

func Filter[A any](array []A, predicate func(A) bool) []A {
	result := []A{}
	for _, val := range array {
		if (predicate(val)) {
			result = append(result, val)
		}
	}
	return result
}

func FindFirstIdx[A any](array []A, predicate func(A) bool) int {
	for idx, val := range array {
		if (predicate(val)) {
			return idx
		}
	}
	return -1
}