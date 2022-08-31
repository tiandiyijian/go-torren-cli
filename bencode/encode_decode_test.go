package bencode

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInt(t *testing.T) {
	assert := assert.New(t)
	buf := new(bytes.Buffer)
	buf.WriteString("i233e")

	bObj, err := Decode(buf)
	assert.Equal(bObj.type_, BINT)
	assert.Equal(bObj.val_, 233)

	bLen, err := bObj.Encode(buf)
	assert.Nil(err)
	assert.Equal(bLen, 5)
	assert.Equal(buf.String(), "i233e")

	buf.Reset()
	buf.WriteString("i-233e")
	bObj, err = Decode(buf)
	assert.Equal(bObj.type_, BINT)
	assert.Equal(bObj.val_, -233)

	bLen, err = bObj.Encode(buf)
	assert.Nil(err)
	assert.Equal(bLen, 6)
	assert.Equal(buf.String(), "i-233e")
}

func TestStr(t *testing.T) {
	assert := assert.New(t)
	buf := new(bytes.Buffer)
	buf.WriteString("3:233")

	bObj, err := Decode(buf)
	assert.Nil(err)
	assert.Equal(bObj.type_, BSTR)
	assert.Equal(bObj.val_, "233")

	bLen, err := bObj.Encode(buf)
	assert.Nil(err)
	assert.Equal(bLen, 5)
	assert.Equal(buf.String(), "3:233")
}

func TestList(t *testing.T) {
	assert := assert.New(t)
	buf := new(bytes.Buffer)
	buf.WriteString("l3:233i123ee")

	bObj, err := Decode(buf)
	assert.Equal(bObj.type_, BLIST)

	l := bObj.val_.([]*BObj)
	assert.Equal(l[0].type_, BSTR)
	assert.Equal(l[0].val_, "233")
	assert.Equal(l[1].type_, BINT)
	assert.Equal(l[1].val_, 123)

	bLen, err := bObj.Encode(buf)
	assert.Nil(err)
	assert.Equal(bLen, len("l3:233i123ee"))
	assert.Equal(buf.String(), "l3:233i123ee")
}

func TestObj(t *testing.T) {
	assert := assert.New(t)

	code := "d4:key13:2334:key2i123ee"
	buf := new(bytes.Buffer)
	buf.WriteString(code)

	bObj, err := Decode(buf)
	assert.Equal(bObj.type_, BDICT)

	dict := bObj.val_.(map[string]*BObj)
	assert.Equal(dict["key1"].type_, BSTR)
	assert.Equal(dict["key1"].val_, "233")
	assert.Equal(dict["key2"].type_, BINT)
	assert.Equal(dict["key2"].val_, 123)

	bLen, err := bObj.Encode(buf)
	assert.Nil(err)
	assert.Equal(bLen, len(code))
	newCode := buf.String()
	assert.Equal(newCode, code)
}
