package utils

func Map[T interface{}, S interface{}](arr []T, f func(T) S) []S {
	var newArr []S
	for _, t := range arr {
		newArr = append(newArr, f(t))
	}
	return newArr
}

func Filter[T interface{}](arr []T, f func(T) bool) []T {
	var newArr []T
	for _, t := range arr {
		if f(t) {
			newArr = append(newArr, t)
		}
	}
	return newArr
}
