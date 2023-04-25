package slices

func MapSlice[IN any, OUT any](input []IN, mapper func(IN) OUT) []OUT {
	res := make([]OUT, 0)

	for _, value := range input {
		res = append(res, mapper(value))
	}

	return res
}

func MapSliceWithError[IN any, OUT any](input []IN, mapper func(IN) (OUT, error)) ([]OUT, error) {
	res := make([]OUT, 0)

	for _, value := range input {
		out, err := mapper(value)
		if err != nil {
			return nil, err
		}
		res = append(res, out)
	}

	return res, nil
}

func MapBy[IN any](input []IN, mapper func(IN) string) map[string]IN {
	res := make(map[string]IN)
	for _, value := range input {
		res[mapper(value)] = value
	}
	return res
}

func SplitSlice[IN any](input []IN, maxChunkSize int) [][]IN {
	chunks := calculateChunks(len(input), maxChunkSize)

	var res [][]IN
	for i := 0; i < chunks; i++ {
		res = append(res, make([]IN, 0))
	}

	for i, value := range input {
		res[i%chunks] = append(res[i%chunks], value)
	}

	return res
}

func calculateChunks(count int, maxChunkSize int) int {
	chunks := count / maxChunkSize

	if count%maxChunkSize == 0 {
		return chunks
	}
	return chunks + 1
}
