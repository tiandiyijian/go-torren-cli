package torrent

import (
	"bytes"
	"crypto/rand"
	"github.com/tiandiyijian/go-torrent-cli/bencode"
	"github.com/tiandiyijian/go-torrent-cli/peer"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var Port = 1024

type Torrent struct {
	Announce    string
	Name        string
	Length      int
	InfoHash    [20]byte
	PieceLength int
	PieceHashes [][20]byte
	Peers       []peer.Peer
	PeerID      [20]byte
}

type TrackerResp struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func (t *Torrent) buildTrackerURL(peerID [20]byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}
	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(Port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Length)},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}

func (t *Torrent) GetPeers() error {
	_, err := rand.Read(t.PeerID[:])
	if err != nil {
		return err
	}

	url, err := t.buildTrackerURL(t.PeerID, uint16(Port))
	//fmt.Println(url)
	c := &http.Client{Timeout: 15 * time.Second}
	resp, err := c.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	var tr TrackerResp
	err = bencode.Unmarshal(bytes.NewBuffer(res), &tr)
	if err != nil {
		return err
	}
	//fmt.Println("tr.Interval", tr.Interval)
	peers, err := peer.Unmarshal([]byte(tr.Peers))
	if err != nil {
		return err
	}

	t.Peers = peers
	return err
}
