//go:build !solution

package reversemap

import (
	"reflect"
)

func ReverseMap(forward interface{}) interface{} {
	// получаем значения forward и его тип
	forwardValue := reflect.ValueOf(forward)
	forwardType := forwardValue.Type()
	if forwardType.Kind() != reflect.Map {
		panic("forward must be map")
	}
	// создаем обратную мапу
	reversedType := reflect.MapOf(forwardType.Elem(), forwardType.Key())
	reverseMap := reflect.MakeMap(reversedType)
	keys := forwardValue.MapKeys()
	for _, key := range keys {
		value := forwardValue.MapIndex(key)
		reverseMap.SetMapIndex(value, key)
	}
	return reverseMap.Interface()
}
