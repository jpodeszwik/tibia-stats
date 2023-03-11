package slices

import (
	"reflect"
	"strconv"
	"testing"
)

func TestMapSlice_Multiply(t *testing.T) {
	numbers := []int{1, 2, 3, 4}

	response := MapSlice(numbers, func(in int) int {
		return in * 2
	})

	expected := []int{2, 4, 6, 8}
	if !reflect.DeepEqual(expected, response) {
		t.Errorf("Response: %v does not match expected: %v", response, expected)
	}

}

func TestMapSlice_ChangeType(t *testing.T) {
	numbers := []int{1, 2, 3, 4}

	response := MapSlice(numbers, func(in int) string {
		return strconv.Itoa(in)
	})

	expected := []string{"1", "2", "3", "4"}
	if !reflect.DeepEqual(expected, response) {
		t.Errorf("Response: %v does not match expected: %v", response, expected)
	}
}
