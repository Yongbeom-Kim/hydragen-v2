package utils

func Map[T interface{}, S interface{}](arr []T, f func(T) S) []S {
	var newArr []S
	for _, t := range arr {
		newArr = append(newArr, f(t))
	}
	return newArr
}
