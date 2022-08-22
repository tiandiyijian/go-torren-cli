package main

import (
	"fmt"
	"github.com/tiandiyijian/go-torrent-cli/bencode"
	"github.com/tiandiyijian/go-torrent-cli/torrentfile"
	"log"
	"os"
)

func main() {

	file, _ := os.Open("files/test.torrent")
	var tf torrentfile.TorrentFile
	err := bencode.Unmarshal(file, &tf)
	if err != nil {
		log.Fatalln("invalid torrent file")
	}
	fmt.Printf("%+v", tf)

}
