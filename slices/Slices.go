package slices

func MapSlice[IN any, OUT any](input []IN, mapper func(IN) OUT) []OUT {
	res := make([]OUT, 0)

	for _, value := range input {
		res = append(res, mapper(value))
	}

	return res
}

func SplitSlice[IN any](input []IN, chunks int) [][]IN {
	res := make([][]IN, 0)
	for i := 0; i < chunks; i++ {
		res = append(res, make([]IN, 0))
	}

	for i, value := range input {
		res[i%chunks] = append(res[i%chunks], value)
	}

	return res
}
