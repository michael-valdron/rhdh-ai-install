package util

import (
	"fmt"
	"reflect"
)

func ExistsInUnstructured(v any, ids ...any) bool {
	if len(ids) == 0 {
		return true
	} else if v == nil {
		return false
	}

	value := reflect.ValueOf(v)

	switch value.Kind() {
	case reflect.Map:
		idxValue := reflect.ValueOf(ids[0])
		idxElement := value.MapIndex(idxValue)
		if !idxElement.IsValid() {
			return false
		}
		return ExistsInUnstructured(idxElement.Interface(), ids[1:]...)
	case reflect.Slice:
		idxValue := ids[0].(int)
		if idxValue >= value.Len() {
			return false
		}
		return ExistsInUnstructured(value.Index(idxValue).Interface(), ids[1:]...)
	default:
		return false
	}
}

func GetValueInUnstructured(v any, ids ...any) any {
	if v == nil || len(ids) == 0 {
		return v
	}

	value := reflect.ValueOf(v)

	switch value.Kind() {
	case reflect.Map:
		idxValue := reflect.ValueOf(ids[0])
		return GetValueInUnstructured(value.MapIndex(idxValue).Interface(), ids[1:]...)
	case reflect.Slice:
		idxValue := ids[0].(int)
		return GetValueInUnstructured(value.Index(idxValue).Interface(), ids[1:]...)
	default:
		return nil
	}
}

func SetValueInUnstructured(v any, value any, ids ...any) {
	if v == nil || len(ids) == 0 {
		return
	}

	specValue := reflect.ValueOf(v)

	switch specValue.Kind() {
	case reflect.Map:
		idxValue := reflect.ValueOf(ids[0])
		if len(ids[1:]) == 0 {
			specValue.SetMapIndex(idxValue, reflect.ValueOf(value))
		} else {
			SetValueInUnstructured(specValue.MapIndex(idxValue).Interface(), value, ids[1:]...)
		}
	case reflect.Slice:
		idxValue := ids[0].(int)
		if len(ids[1:]) == 0 {
			specValue.Index(idxValue).Set(reflect.ValueOf(value))
		} else {
			SetValueInUnstructured(specValue.Index(idxValue).Interface(), value, ids[1:]...)
		}
	}
}

func UnstructuredToMap[T any](v any) (map[string]T, error) {
	switch result := v.(type) {
	case map[string]T:
		return result, nil
	default:
		return make(map[string]T), fmt.Errorf("invalid map entity: %v", v)
	}
}

func UnstructuredToSlice[T any](v any) ([]T, error) {
	switch result := v.(type) {
	case []T:
		return result, nil
	default:
		return []T{}, fmt.Errorf("invalid slice entity: %v", v)
	}
}
