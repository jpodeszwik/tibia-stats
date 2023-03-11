package slices

func MapSlice[IN any, OUT any](input []IN, mapper func(IN) OUT) []OUT {
	res := make([]OUT, 0)

	for _, value := range input {
		res = append(res, mapper(value))
	}

	return res
}
