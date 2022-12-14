package client

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"github.com/tiandiyijian/go-torrent-cli/bitfield"
	"github.com/tiandiyijian/go-torrent-cli/handshake"
	"github.com/tiandiyijian/go-torrent-cli/message"
	"github.com/tiandiyijian/go-torrent-cli/peer"
)

type Client struct {
	Conn     net.Conn
	Chocked  bool
	BitField bitfield.BitField
	peer     peer.Peer
	infoHash [20]byte
	peerID   [20]byte
}

func completeHandshake(conn net.Conn, infoHash, peerID [20]byte) (*handshake.Handshake, error) {
	conn.SetDeadline(time.Now().Add(time.Second * 3))
	defer conn.SetDeadline(time.Time{}) // Disable the deadline

	req := handshake.New(infoHash, peerID)
	_, err := conn.Write(req.Serialize())
	if err != nil {
		return nil, err
	}

	res, err := handshake.Read(conn)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(res.InfoHash[:], infoHash[:]) {
		return nil, fmt.Errorf("Expected infohash %x but got %x", infoHash, res.InfoHash)
	}

	return res, nil
}

func recvBitField(conn net.Conn) (bitfield.BitField, error) {
	conn.SetDeadline(time.Now().Add(time.Second * 5))
	defer conn.SetDeadline(time.Time{}) // Disable the deadline

	msg, err := message.Read(conn)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		err := fmt.Errorf("expected bitfield but got %s", msg)
		return nil, err
	}

	if msg.ID != message.MsgBitfield {
		err := fmt.Errorf("expected bitfield but got %s", msg)
		return nil, err
	}

	return msg.Payload, nil
}

// New connects with a peer, completes a handshake, and receives a handshake
// returns an err if any of those fail.
func New(peer peer.Peer, infoHash, PeerID [20]byte) (*Client, error) {
	conn, err := net.DialTimeout("tcp", peer.String(), 3*time.Second)
	fmt.Println("dial peer", peer.String())
	if err != nil {
		fmt.Println("Dial peer failed. Error:", err)
		return nil, err
	}

	_, err = completeHandshake(conn, infoHash, PeerID)
	if err != nil {
		fmt.Println("Handshake peer failed. Error:", err)
		return nil, err
	}

	bf, err := recvBitField(conn)
	if err != nil {
		conn.Close()
		fmt.Println("Receive peer bitfield failed. Error:", err)
		return nil, err
	}

	return &Client{
		Conn:     conn,
		Chocked:  true,
		BitField: bf,
		peer:     peer,
		infoHash: infoHash,
		peerID:   PeerID,
	}, nil
}

func (c *Client) Read() (*message.Message, error) {
	return message.Read(c.Conn)
}

// SendRequest sends a Request message to the peer
func (c *Client) SendRequest(index, begin, length int) error {
	req := message.FormatRequest(index, begin, length)
	_, err := c.Conn.Write(req.Serialize())
	return err
}

// SendInterested sends an Interested message to the peer
func (c *Client) SendInterested() error {
	msg := message.Message{ID: message.MsgInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendNotInterested sends a NotInterested message to the peer
func (c *Client) SendNotInterested() error {
	msg := message.Message{ID: message.MsgNotInterested}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendUnchoke sends an Unchoke message to the peer
func (c *Client) SendUnchoke() error {
	msg := message.Message{ID: message.MsgUnchoke}
	_, err := c.Conn.Write(msg.Serialize())
	return err
}

// SendHave sends a Have message to the peer
func (c *Client) SendHave(index int) error {
	msg := message.FormatHave(index)
	_, err := c.Conn.Write(msg.Serialize())
	return err
}
