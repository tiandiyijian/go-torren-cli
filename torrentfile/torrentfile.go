package torrentfile

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/tiandiyijian/go-torrent-cli/bencode"
	"github.com/tiandiyijian/go-torrent-cli/torrent"
)

type Info struct {
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
}

type TorrentFile struct {
	Announce string `bencode:"announce"`
	Info     Info   `bencode:"info"`
}

func (i *Info) Hash() ([20]byte, error) {
	var buf bytes.Buffer
	_, err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func (i *Info) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20
	buf := []byte(i.Pieces)
	if len(buf)%hashLen != 0 {
		return nil, fmt.Errorf("malformed pieces hash length: %d", len(buf))
	}

	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}

	return hashes, nil
}

func (tf *TorrentFile) ToTorrent() (torrent.Torrent, error) {
	infoHash, err := tf.Info.Hash()
	if err != nil {
		return torrent.Torrent{}, err
	}

	piecesHash, err := tf.Info.splitPieceHashes()
	if err != nil {
		return torrent.Torrent{}, err
	}

	return torrent.Torrent{
		Announce:    tf.Announce,
		Name:        tf.Info.Name,
		Length:      tf.Info.Length,
		InfoHash:    infoHash,
		PieceLength: tf.Info.PieceLength,
		PieceHashes: piecesHash,
	}, nil
}
