package torrent

import (
	"github.com/tiandiyijian/go-torrent-cli/client"
	"github.com/tiandiyijian/go-torrent-cli/message"
)

type pieceWork struct {
	index  int
	hash   [20]byte
	length int
}

type pieceResult struct {
	index int
	buf   []byte
}

type pieceProgress struct {
	index      int
	client     *client.Client
	buf        []byte
	downloaded int
	requested  int
	backlog    int
}

func (p *pieceProgress) readMessage() error {
	msg, err := p.client.Read()
	if err != nil {
		return err
	}
	if msg == nil {
		return nil
	}

	switch msg.ID {
	case message.MsgUnchoke:
		p.client.Chocked = false
	case message.MsgChoke:
		p.client.Chocked = true
	case message.MsgHave:
		index, err := message.ParseHave(msg)
		if err != nil {
			return err
		}
		p.client.BitField.SetPiece(index)
	case message.MsgPiece:
		n, err := message.ParsePiece(p.index, p.buf, msg)
		if err != nil {
			return err
		}
		p.downloaded += n
		p.backlog--
	}

	return nil
}
