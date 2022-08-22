package bencode

import (
	"bufio"
	"io"
	"strconv"
)

func (o *BObj) Encode(w io.Writer) (int, error) {
	bw, ok := w.(*bufio.Writer)
	if !ok {
		bw = bufio.NewWriter(w)
	}

	curLen := 0
	switch o.type_ {
	case BINT:
		num, _ := o.Int()
		wLen, err := EncodeInt(bw, num)
		if err != nil {
			return 0, err
		}
		curLen += wLen
	case BSTR:
		s, _ := o.Str()
		wLen, err := EncodeStr(bw, s)
		if err != nil {
			return 0, err
		}
		curLen += wLen
	case BLIST:
		list, _ := o.List()

		curLen += 1
		bw.WriteByte('l')

		for _, bobj := range list {
			wLen, err := bobj.Encode(bw)
			if err != nil {
				return 0, err
			}
			curLen += wLen
		}

		curLen += 1
		bw.WriteByte('e')
	case BDICT:
		dict, _ := o.Dict()

		curLen += 1
		bw.WriteByte('d')

		for key, bobj := range dict {
			wLen, err := EncodeStr(bw, key)
			if err != nil {
				return 0, err
			}
			curLen += wLen

			wLen, err = bobj.Encode(bw)
			if err != nil {
				return 0, err
			}
			curLen += wLen
		}

		curLen += 1
		bw.WriteByte('e')
	}

	if err := bw.Flush(); err != nil {
		return 0, err
	}

	return curLen, nil
}

func EncodeStr(w io.Writer, s string) (int, error) {
	bw := bufio.NewWriter(w)
	wLen := 0
	size := len(s)
	Ssize := strconv.Itoa(size)

	wLen += len(Ssize)
	bw.WriteString(Ssize)

	wLen += 1
	bw.WriteByte(':')

	wLen += len(s)
	bw.WriteString(s)

	if err := bw.Flush(); err != nil {
		return 0, err
	}

	return wLen, nil
}

func EncodeInt(w io.Writer, num int) (int, error) {
	bw := bufio.NewWriter(w)

	wLen := 1
	bw.WriteByte('i')

	s := strconv.Itoa(num)
	wLen += 1 + len(s)
	bw.WriteString(s)
	bw.WriteByte('e')

	if err := bw.Flush(); err != nil {
		return 0, err
	}

	return wLen, nil
}
