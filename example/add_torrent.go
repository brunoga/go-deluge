package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	deluge "github.com/brunoga/go-deluge"
)

// Flags.
var delugeAddress = flag.String("deluge_address", "127.0.0.1:8112", "deluge server address")
var delugePassword = flag.String("deluge_password", "", "deluge server password")
var addPaused = flag.Bool("add_paused", false, "if true, torrent is added in the paused state")

func main() {
	flag.Parse()

	if *delugeAddress == "" || *delugePassword == "" {
		fmt.Printf("-deluge_address and -deluge_password must be provided and non-empty")
		return
	}

	d, err := deluge.New("http://" + *delugeAddress + "/json", *delugePassword)
	if err != nil {
		panic(err)
	}

	options := map[string]interface{}{
		"add_paused" : *addPaused,
	}

	for _, torrent := range flag.Args() {
		var id string
		if strings.HasPrefix(torrent, "magnet:") {
			id, err = d.CoreAddTorrentMagnet(torrent, options)
			if err != nil {
				fmt.Println("Error adding torrent via magnet URL :", err)
				continue
			}
		} else if strings.Index(torrent, "://") != -1 {
			id, err = d.CoreAddTorrentUrl(torrent, options)
			if err != nil {
				fmt.Println("Error adding torrent via URL :", err)
				continue
			}
		} else {
			file, err :=  os.Open(torrent)
			if err != nil {
				fmt.Println("Error opening local file :", err)
				continue
			}

			defer file.Close()

			fileName := filepath.Base(torrent)

			data, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Println("Error reading local file :", err)
				continue
			}

			var buffer bytes.Buffer
			encoder := base64.NewEncoder(base64.StdEncoding, &buffer)
			encoder.Write(data)
			encoder.Close()

			id, err = d.CoreAddTorrentFile(fileName, buffer.String(), options)
			if err != nil {
				fmt.Println("Error adding torrent file :", err)
			}
		}

		fmt.Println("Torrent added. Id : ", id)
	}
}
