package bencode

import (
	"errors"
	"io"
	"reflect"
	"strings"
)

func Marshal(w io.Writer, i interface{}) (int, error) {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return marshalValue(w, v)
}

func marshalValue(w io.Writer, v reflect.Value) (int, error) {
	curLen := 0
	switch v.Kind() {
	case reflect.Int:
		wLen, err := EncodeInt(w, int(v.Int()))
		if err != nil {
			return 0, err
		}

		curLen += wLen
	case reflect.String:
		wLen, err := EncodeStr(w, v.String())
		if err != nil {
			return 0, err
		}

		curLen += wLen
	case reflect.Slice:
		wLen, err := marshalList(w, v)
		if err != nil {
			return 0, err
		}

		curLen += wLen
	case reflect.Struct:
		wLen, err := marshalDict(w, v)
		if err != nil {
			return 0, err
		}

		curLen += wLen
	default:
		return 0, errors.New("unsupported type")
	}

	return curLen, nil
}

func marshalList(w io.Writer, v reflect.Value) (int, error) {
	curLen := 1
	w.Write([]byte{'l'})

	for i := 0; i < v.Len(); i++ {
		curV := v.Index(i)
		wLen, err := marshalValue(w, curV)
		if err != nil {
			return 0, err
		}
		curLen += wLen
	}

	w.Write([]byte{'e'})
	return curLen + 1, nil
}

func marshalDict(w io.Writer, v reflect.Value) (int, error) {
	curLen := 1
	w.Write([]byte{'d'})

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)

		key := fieldType.Tag.Get("bencode")
		if key == "" {
			key = strings.ToLower(fieldType.Name)
		}

		wLen, err := EncodeStr(w, key)
		if err != nil {
			return 0, err
		}
		curLen += wLen

		wLen, err = marshalValue(w, field)
		if err != nil {
			return 0, err
		}
		curLen += wLen
	}

	w.Write([]byte{'e'})
	return curLen + 1, nil
}
