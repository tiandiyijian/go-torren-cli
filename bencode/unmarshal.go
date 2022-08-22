package bencode

import (
	"errors"
	"io"
	"reflect"
	"strings"
)

func Unmarshal(r io.Reader, s interface{}) error {
	o, err := Decode(r)
	if err != nil {
		return err
	}

	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr {
		return errors.New("dst must be a pointer")
	}

	switch o.type_ {
	case BLIST:
		list, _ := o.List()
		l := reflect.MakeSlice(v.Elem().Type(), len(list), len(list))
		v.Elem().Set(l)
		err = unmarshalList(list, v)
		if err != nil {
			return err
		}
	case BDICT:
		dict, _ := o.Dict()
		err := unmarshalDict(dict, v)
		if err != nil {
			return err
		}
	default:
		return errors.New("src must be struct or slice")
	}

	return err
}

func unmarshalList(list []*BObj, v reflect.Value) error {
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		return errors.New("dst must be pointer of slice")
	}

	v = v.Elem()
	if len(list) == 0 {
		return nil
	}

	switch list[0].type_ {
	case BSTR:
		for i, bObj := range list {
			s, err := bObj.Str()
			if err != nil {
				return err
			}
			v.Index(i).SetString(s)
		}
	case BINT:
		for i, bObj := range list {
			num, err := bObj.Int()
			if err != nil {
				return err
			}
			v.Index(i).SetInt(int64(num))
		}
	case BLIST:
		if v.Type().Elem().Kind() != reflect.Slice {
			return ErrTyp
		}
		for i, bObj := range list {
			bList, err := bObj.List()
			if err != nil {
				return err
			}

			newPtr := reflect.New(v.Type().Elem())
			newSlice := reflect.MakeSlice(v.Type().Elem(), len(bList), len(bList))
			newPtr.Elem().Set(newSlice)
			err = unmarshalList(bList, newPtr)
			if err != nil {
				return err
			}
			v.Index(i).Set(newPtr.Elem())
		}
	case BDICT:
		if v.Type().Elem().Kind() != reflect.Struct {
			return ErrTyp
		}
		for i, bObj := range list {
			bDict, err := bObj.Dict()
			if err != nil {
				return err
			}

			newPtr := reflect.New(v.Type().Elem())
			err = unmarshalDict(bDict, newPtr)
			if err != nil {
				return err
			}
			v.Index(i).Set(newPtr.Elem())
		}
	}

	return nil
}

func unmarshalDict(bDict map[string]*BObj, v reflect.Value) error {
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("dst must be pointer of struct")
	}

	v = v.Elem()
	for i, n := 0, v.NumField(); i < n; i++ {
		field := v.Field(i)
		if !field.CanSet() {
			continue
		}

		fieldType := v.Type().Field(i)
		key := fieldType.Tag.Get("bencode")
		if key == "" {
			key = strings.ToLower(fieldType.Name)
		}

		bObj, ok := bDict[key]
		if !ok {
			continue
		}

		switch bObj.type_ {
		case BINT:
			if fieldType.Type.Kind() != reflect.Int {
				break
			}

			num, _ := bObj.Int()
			field.SetInt(int64(num))
		case BSTR:
			if fieldType.Type.Kind() != reflect.String {
				break
			}

			s, _ := bObj.Str()
			field.SetString(s)
		case BLIST:
			if fieldType.Type.Kind() != reflect.Slice {
				break
			}

			bList, _ := bObj.List()
			newPtr := reflect.New(fieldType.Type)
			newSlice := reflect.MakeSlice(fieldType.Type, len(bList), len(bList))
			newPtr.Elem().Set(newSlice)
			err := unmarshalList(bList, newPtr)
			if err != nil {
				break
			}

			field.Set(newPtr.Elem())
		case BDICT:
			if fieldType.Type.Kind() != reflect.Struct {
				break
			}

			bDict, _ := bObj.Dict()
			newPtr := reflect.New(fieldType.Type)
			err := unmarshalDict(bDict, newPtr)
			if err != nil {
				break
			}

			field.Set(newPtr.Elem())
		}
	}
	return nil
}
