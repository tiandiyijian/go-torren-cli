package bencode

import (
	"bufio"
	"io"
)

func Decode(r io.Reader) (*BObj, error) {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}

	b, err := br.Peek(1)
	if err != nil {
		return nil, err
	}

	switch {
	case b[0] >= '0' && b[0] <= '9':
		s, err := DecodeString(br)
		if err != nil {
			return nil, err
		}

		return NewBObj(BSTR, s), nil
	case b[0] == 'i':
		num, err := DecodeInt(br)
		if err != nil {
			return nil, err
		}

		return NewBObj(BINT, num), nil
	case b[0] == 'l':
		br.ReadByte() // start l
		var list []*BObj
		for {
			b, err := br.Peek(1)
			if err != nil {
				return nil, err
			}

			if b[0] == 'e' { // end e
				br.ReadByte()
				return NewBObj(BLIST, list), nil
			}

			bobj, err := Decode(br)
			if err != nil {
				return nil, err
			}
			list = append(list, bobj)
		}
	case b[0] == 'd':
		br.ReadByte() // start d
		dict := map[string]*BObj{}
		for {
			b, err := br.Peek(1)
			if err != nil {
				return nil, err
			}

			if b[0] == 'e' { // end e
				br.ReadByte()
				return NewBObj(BDICT, dict), nil
			}

			key, err := DecodeString(br)
			if err != nil {
				return nil, err
			}

			bobj, err := Decode(br)
			if err != nil {
				return nil, err
			}
			dict[key] = bobj
		}
	}

	return nil, ErrIvd
}

func DecodeString(br *bufio.Reader) (string, error) {
	var num int
	for peek, _ := br.Peek(1); peek[0] >= '0' && peek[0] <= '9'; peek, _ = br.Peek(1) {
		c, err := br.ReadByte()
		if err != nil {
			return "", err
		}
		num = num*10 + int(c-'0')
	}

	c, err := br.ReadByte()
	if err != nil {
		return "", err
	}
	if c != ':' {
		return "", ErrCol
	}

	buf := make([]byte, num)
	_, err = io.ReadAtLeast(br, buf, num)
	if err != nil {
		return "", err
	}
	//fmt.Println(string(buf))
	return string(buf), nil
}

func DecodeInt(br *bufio.Reader) (int, error) {
	br.ReadByte() // start i

	sign := 1
	b, err := br.Peek(1)
	if err != nil {
		return 0, err
	}
	if b[0] == '-' {
		sign = -1
		br.ReadByte()
	}

	var num int
	for b, _ := br.Peek(1); b[0] >= '0' && b[0] <= '9'; b, _ = br.Peek(1) {
		c, err := br.ReadByte()
		if err != nil {
			return 0, err
		}
		num = num*10 + int(c-'0')
	}

	c, err := br.ReadByte() // end e
	if err != nil {
		return 0, err
	}
	if c != 'e' {
		return 0, ErrEpE
	}

	return sign * num, nil
}
