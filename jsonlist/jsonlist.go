//go:build !solution

package jsonlist

import (
	"bufio"
	"encoding/json"
	"io"
	"reflect"
)

// Marshal преобразует слайс значений в последовательность JSON объектов и записывает её в io.Writer.

func Marshal(w io.Writer, slice interface{}) error {
	sliceType := reflect.TypeOf(slice)
	if sliceType.Kind() != reflect.Slice {
		return &json.UnsupportedTypeError{Type: sliceType}
	}
	buffer := bufio.NewWriter(w)
	encoder := json.NewEncoder(buffer)
	len := reflect.ValueOf(slice).Len()
	for i := 0; i < len; i++ {
		if err := encoder.Encode(reflect.ValueOf(slice).Index(i).Interface()); err != nil {
			return err
		}
	}
	if err := buffer.Flush(); err != nil {
		return err
	}

	return nil
}

func Unmarshal(r io.Reader, slice interface{}) error {
	sliceType := reflect.TypeOf(slice)
	if sliceType.Kind() != reflect.Ptr || sliceType.Elem().Kind() != reflect.Slice {
		return &json.UnsupportedTypeError{Type: sliceType}
	}
	decoder := json.NewDecoder(r)
	sliceValue := reflect.ValueOf(slice).Elem()
	sliceValue.Set(reflect.MakeSlice(sliceType.Elem(), 0, 0))
	for {
		newValue := reflect.New(sliceType.Elem().Elem())
		if err := decoder.Decode(newValue.Interface()); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		sliceValue.Set(reflect.Append(sliceValue, newValue.Elem()))
	}
	return nil
}
