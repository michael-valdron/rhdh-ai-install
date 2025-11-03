package util

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func indexInRange[S []E, E any](s S, idx int) bool {
	return idx > -1 && idx < len(s)
}

func GetElementFromSliceByFunc[S []E, E any](s S, accessor func(E) bool) (*E, error) {
	idx := slices.IndexFunc(s, accessor)

	if !indexInRange(s, idx) {
		return nil, fmt.Errorf("element could not be found with provided accessor function")
	}

	return &s[idx], nil
}

func ModifySliceElementByFunc[S []E, E any](s S, accessor func(E) bool, modifier func(*E)) error {
	idx := slices.IndexFunc(s, accessor)

	if indexInRange(s, idx) {
		modifier(&s[idx])
	} else {
		return fmt.Errorf("element could not be found with provided accessor function")
	}

	return nil
}

func RemoveFromSliceByIndex[S []E, E any](s S, idx int) (S, error) {
	if indexInRange(s, idx) {
		return append(s[:idx], s[idx+1:]...), nil
	} else {
		return s, fmt.Errorf("index out of range")
	}
}

func RemoveFromSliceByValue[S []E, E comparable](s S, v E) (S, error) {
	idx := slices.Index(s, v)

	if result, err := RemoveFromSliceByIndex(s, idx); err != nil {
		return s, fmt.Errorf("element '%v' was not found in slice", v)
	} else {
		return result, nil
	}
}

func RemoveFromSliceByFunc[S []E, E any](s S, f func(E) bool) (S, error) {
	idx := slices.IndexFunc(s, f)

	if result, err := RemoveFromSliceByIndex(s, idx); err != nil {
		return s, fmt.Errorf("element was not found in slice using provided function")
	} else {
		return result, nil
	}
}

func Join(s []any, sep string) string {
	result := make([]string, len(s))

	for idx := range s {
		switch typed := s[idx].(type) {
		case string:
			result[idx] = typed
		case int, int8, int16, int32, int64:
			result[idx] = strconv.FormatInt(typed.(int64), 10)
		case uint, uint8, uint16, uint32, uint64:
			result[idx] = strconv.FormatUint(typed.(uint64), 10)
		case float32:
			result[idx] = strconv.FormatFloat(float64(typed), 'f', -1, 32)
		case float64:
			result[idx] = strconv.FormatFloat(typed, 'f', -1, 64)
		case bool:
			result[idx] = strconv.FormatBool(typed)
		case fmt.Stringer:
			result[idx] = typed.String()
		case error:
			result[idx] = typed.Error()
		case []string:
			result[idx] = "[" + strings.Join(typed, ", ") + "]"
		case []int:
			strSlice := make([]string, len(typed))
			for i, val := range typed {
				strSlice[i] = strconv.Itoa(val)
			}
			result[idx] = "{" + strings.Join(strSlice, ", ") + "}"
		default:
			result[idx] = fmt.Sprintf("%v", typed)
		}
	}

	return strings.Join(result, sep)
}
