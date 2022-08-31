package bencode

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalDict(t *testing.T) {
	assert := assert.New(t)

	code := "d2:aai233e2:bb10:i am a str2:ccli666ei777ei888ee2:ddd2:eei555e2:ff14:i am a str tooee"
	buf := bytes.NewBufferString(code)

	var e example
	exp := example{
		A: 233,
		B: "i am a str",
		C: []int{666, 777, 888},
		D: sub{A: 555, B: "i am a str too"},
	}
	err := Unmarshal(buf, &e)

	assert.Nil(err)
	assert.Equal(e, exp)
}

func TestUnmarshalList(t *testing.T) {
	assert := assert.New(t)

	code := "l2:aa3:bbbe"
	buf := bytes.NewBufferString(code)

	var e []string
	exp := []string{"aa", "bbb"}
	err := Unmarshal(buf, &e)

	assert.Nil(err)
	assert.Equal(e, exp)
}

func TestUnmarshalListList(t *testing.T) {
	assert := assert.New(t)

	code := "ll2:aael3:bbb3:cccee"
	buf := bytes.NewBufferString(code)

	var e [][]string
	exp := [][]string{{"aa"}, {"bbb", "ccc"}}
	err := Unmarshal(buf, &e)

	assert.Nil(err)
	assert.Equal(e, exp)
}
