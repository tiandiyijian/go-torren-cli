package main

import (
	"log"
	"os"

	"github.com/tiandiyijian/go-torrent-cli/bencode"
	"github.com/tiandiyijian/go-torrent-cli/torrentfile"
)

func main() {
	inPath := os.Args[1]
	outPath := os.Args[2]
	file, _ := os.Open(inPath)

	var tf torrentfile.TorrentFile
	err := bencode.Unmarshal(file, &tf)
	if err != nil {
		log.Fatalln("invalid torrent file")
	}

	t, err := tf.ToTorrent()
	if err != nil {
		log.Fatalln("torrentfile to torrent err")
	}

	err = t.GetPeers()
	if err != nil {
		log.Fatalln("get peers error", err)
	}

	buf, err := t.Download()
	if err != nil {
		log.Fatalln("download error")
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		log.Fatalln("create file error")
	}
	defer outFile.Close()
	_, err = outFile.Write(buf)
	if err != nil {
		log.Fatalln("write file error")
	}
}
