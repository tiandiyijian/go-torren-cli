package main

import (
	"github.com/tiandiyijian/go-torrent-cli/bencode"
	"github.com/tiandiyijian/go-torrent-cli/torrentfile"
	"log"
	"os"
)

func main() {
	file, _ := os.Open("files/debian.torrent")
	var tf torrentfile.TorrentFile
	err := bencode.Unmarshal(file, &tf)
	if err != nil {
		log.Fatalln("invalid torrent file")
	}

	t, err := tf.ToTorrent()
	if err != nil {
		log.Fatalln("torrentfile to torrent err")
	}
	t.GetPeers()
}
