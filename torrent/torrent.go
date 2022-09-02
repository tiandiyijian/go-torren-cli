package torrent

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"io"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"time"

	"github.com/tiandiyijian/go-torrent-cli/bencode"
	"github.com/tiandiyijian/go-torrent-cli/client"
	"github.com/tiandiyijian/go-torrent-cli/peer"
)

const (
	Port         = 6881
	MaxBacklog   = 50
	MaxBlockSize = 16384
)

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

	peers, err := peer.Unmarshal([]byte(tr.Peers))
	if err != nil {
		return err
	}

	t.Peers = peers
	return nil
}

func (t *Torrent) Download() ([]byte, error) {
	log.Println("Starting download for", t.Name)

	// pieceQueue for piece to download
	pieceQueue := make(chan *pieceWork, len(t.PieceHashes))
	results := make(chan *pieceResult)
	for index, pieceHash := range t.PieceHashes {
		length := t.PieceLength
		if left := t.Length - index*t.PieceLength; left < t.PieceLength { // the last piece length may be less than piecelength
			length = left
		}
		pieceQueue <- &pieceWork{index: index, hash: pieceHash, length: length}
	}

	for _, peer := range t.Peers {
		go t.startDownloadWorker(peer, pieceQueue, results)
	}

	// Collect results
	buf := make([]byte, t.Length)
	collectedPieces := 0
	for collectedPieces < len(t.PieceHashes) {
		result := <-results

		begin := result.index * t.PieceLength
		end := begin + t.PieceLength
		if end > t.Length {
			end = t.Length
		}

		copy(buf[begin:end], result.buf)
		collectedPieces++

		percent := float32(collectedPieces) / float32(len(t.PieceHashes)) * 100
		numWorkers := runtime.NumGoroutine() - 1 // subtract 1 for main thread
		log.Printf("(%0.2f%%) Downloaded piece #%d from %d peers\n", percent, result.index, numWorkers)
	}

	close(pieceQueue)
	return buf, nil
}

func (t *Torrent) startDownloadWorker(peer peer.Peer, pieceQueue chan *pieceWork, results chan *pieceResult) {
	c, err := client.New(peer, t.InfoHash, t.PeerID)
	if err != nil {
		log.Printf("Could not handshake with %s. Error: %s. Disconnecting\n", peer.IP, err)
		return
	}
	defer c.Conn.Close()
	log.Printf("Completed handshake with %s\n", peer.IP)

	c.SendUnchoke()    // remote peer can send msg
	c.SendInterested() // want to download

	for pw := range pieceQueue {
		if !c.BitField.HasPiece(pw.index) {
			pieceQueue <- pw
			continue
		}

		buf, err := downloadPiece(c, pw)
		if err != nil {
			log.Println("exiting worker with error:", err)
			pieceQueue <- pw
			return
		}

		if !checkHash(pw.hash, buf) {
			log.Printf("Piece #%d failed integrity check\n", pw.index)
			pieceQueue <- pw
			return
		}

		c.SendHave(pw.index)
		results <- &pieceResult{index: pw.index, buf: buf}
	}
}

func downloadPiece(c *client.Client, pw *pieceWork) ([]byte, error) {
	progress := pieceProgress{
		index:  pw.index,
		client: c,
		buf:    make([]byte, pw.length),
	}

	// Setting a deadline helps get unresponsive peers unstuck.
	// 30 seconds is more than enough time to download a 256 KB piece
	c.Conn.SetDeadline(time.Now().Add(time.Second * 6000))
	defer c.Conn.SetDeadline(time.Time{})

	for progress.downloaded < pw.length {
		// If unchoked, send requests until have enough unfulfilled requests
		if !c.Chocked {
			for progress.backlog < MaxBacklog && progress.requested < pw.length {
				blockSize := MaxBlockSize
				// Last block might be shorter than the typical block
				if pw.length-progress.requested < blockSize {
					blockSize = pw.length - progress.requested
				}

				err := c.SendRequest(pw.index, progress.requested, blockSize)
				if err != nil {
					return nil, err
				}

				progress.requested += blockSize
				progress.backlog++
			}
		}

		err := progress.readMessage()
		if err != nil {
			return nil, err
		}
	}

	return progress.buf, nil
}

func checkHash(hash [20]byte, buf []byte) bool {
	bufHash := sha1.Sum(buf)
	return bytes.Equal(hash[:], bufHash[:])
}
