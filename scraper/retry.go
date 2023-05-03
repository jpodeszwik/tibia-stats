package scraper

import "fmt"

func retry[T any](f func() (T, error), times int) (T, error) {
	var errs []error
	for i := 0; i < times; i++ {
		val, err := f()
		if err == nil {
			return val, err
		}
		errs = append(errs, err)
	}

	var zero T
	return zero, fmt.Errorf("failed after %v tries, %v", times, errs)
}
