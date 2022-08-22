package bencode

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshal(t *testing.T) {
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
