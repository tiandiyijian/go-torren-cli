package handshake

import (
	"fmt"
	"io"
)

type Handshake struct {
	Pstr     string // always "BitTorrent protocol"
	InfoHash [20]byte
	PeerID   [20]byte
}

func New(infoHash, peerID [20]byte) *Handshake {
	return &Handshake{
		Pstr:     "BitTorrent protocol",
		InfoHash: infoHash,
		PeerID:   peerID,
	}
}

// Serialize encode form: <19(pstr length)><pstr(BitTorrent protocol)><8 reserved bytes all 0><infohash><peerid>
func (h *Handshake) Serialize() []byte {
	buf := make([]byte, len(h.Pstr)+49)
	buf[0] = byte(len(h.Pstr))

	offset := 1
	offset += copy(buf[offset:], h.Pstr)
	offset += copy(buf[offset:], make([]byte, 8))
	offset += copy(buf[offset:], h.InfoHash[:])
	offset += copy(buf[offset:], h.PeerID[:])

	return buf
}

func Read(r io.Reader) (*Handshake, error) {
	buf := make([]byte, 1) // start length
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}

	pstrLen := int(buf[0])
	if pstrLen == 0 {
		return nil, fmt.Errorf("pstrlen cannot be 0")
	}

	bufContent := make([]byte, pstrLen+48)
	_, err = io.ReadFull(r, bufContent)
	if err != nil {
		return nil, err
	}

	var infoHash, peerID [20]byte
	copy(infoHash[:], bufContent[pstrLen+8:pstrLen+28])
	copy(peerID[:], bufContent[pstrLen+28:])

	h := Handshake{
		Pstr:     string(bufContent[:pstrLen]),
		InfoHash: infoHash,
		PeerID:   peerID,
	}

	return &h, err
}
