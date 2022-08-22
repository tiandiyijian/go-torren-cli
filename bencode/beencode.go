package bencode

import (
	"errors"
)

type BType uint8
type BVal interface{}

const (
	BINT BType = iota
	BSTR
	BLIST
	BDICT
)

var (
	ErrNum = errors.New("expect num")
	ErrCol = errors.New("expect colon")
	ErrEpI = errors.New("expect char i")
	ErrEpE = errors.New("expect char e")
	ErrTyp = errors.New("wrong type")
	ErrIvd = errors.New("invalid bencode")
)

type BObj struct {
	type_ BType
	val_  BVal
}

func NewBObj(typ BType, val BVal) *BObj {
	return &BObj{
		type_: typ,
		val_:  val,
	}
}

func (o *BObj) Int() (int, error) {
	if o.type_ != BINT {
		return 0, ErrTyp
	}

	return o.val_.(int), nil
}

func (o *BObj) Str() (string, error) {
	if o.type_ != BSTR {
		return "", ErrTyp
	}

	return o.val_.(string), nil
}

func (o *BObj) List() ([]*BObj, error) {
	if o.type_ != BLIST {
		return nil, ErrTyp
	}

	return o.val_.([]*BObj), nil
}

func (o *BObj) Dict() (map[string]*BObj, error) {
	if o.type_ != BDICT {
		return nil, ErrTyp
	}

	return o.val_.(map[string]*BObj), nil
}
