package bencode

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

type sub struct {
	A int    `bencode:"ee"`
	B string `bencode:"ff"`
}

type example struct {
	A int    `bencode:"aa"`
	B string `bencode:"bb"`
	C []int  `bencode:"cc"`
	D sub    `bencode:"dd"`
}

func TestMarshal(t *testing.T) {
	assert := assert.New(t)

	e := example{
		A: 233,
		B: "i am a str",
		C: []int{666, 777, 888},
		D: sub{A: 555, B: "i am a str too"},
	}
	code := "d2:aai233e2:bb10:i am a str2:ccli666ei777ei888ee2:ddd2:eei555e2:ff14:i am a str tooee"

	buf := new(bytes.Buffer)
	wLen, err := Marshal(buf, e)
	assert.Nil(err)
	assert.Equal(wLen, len(code))
	assert.Equal(buf.String(), code)
}
