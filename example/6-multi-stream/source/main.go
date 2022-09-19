package main

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/yomorun/yomo"
	"github.com/yomorun/yomo/core/frame"
)

var exit chan bool

type fileInfo struct {
	Name string `json:"name"`
}

func main() {
	if len(os.Args) < 2 {
		panic("please set the dir in args, f.e. go run main.go /path/to/dir")
	}
	fp, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	fi, err := fp.Stat()
	if err != nil {
		panic(err)
	}

	if fi.IsDir() {
		panic("send files, not dir")
	}

	// init yomo-source
	client := yomo.NewSource(
		"source",
		yomo.WithZipperAddr(os.Getenv("zipper")),
		yomo.WithObserveDataTags(0x10),
	)
	defer client.Close()

	// connect to yomo-zipper
	err = client.Connect()
	if err != nil {
		panic(err)
	}
	// PROBLEM: how to wait util the connection is established?
	time.Sleep(time.Second)

	go sendFile(client, fp.Name())

	exit = make(chan bool)

	<-exit
	client.Close()
}

// sendFile sends the file to yomo-zipper.
func sendFile(client yomo.Source, fileName string) {
	log.Printf(">>>>>sending file %s to yomo-zipper...", fileName)
	videoStream, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	// open a new stream
	ctx := context.Background()
	writer, err := client.OpenStream(ctx)
	if err != nil {
		panic(err)
	}

	// send a stream frame with file name
	f := frame.NewStreamFrame([]byte{0x11})
	info := fileInfo{
		Name: filepath.Base(fileName),
	}
	meta, _ := json.Marshal(info)
	f.SetMetadata(meta)
	_, err = writer.Write(f.Encode())
	if err != nil {
		panic(err)
	}

	// calculate the md5 of the file
	pipeReader, pipeWriter := io.Pipe()
	go calculateMD5(pipeReader, fileName)

	// send video stream to yomo-zipper
	written, err := io.Copy(io.MultiWriter(pipeWriter, writer), videoStream)
	if err != nil && err != io.EOF {
		panic(err)
	}
	log.Printf("file: %s, written: %d\n", fileName, written)
	writer.Close()
	pipeWriter.Close()
	videoStream.Close()
	exit <- true
}

// calculateMD5 calculates the md5 of the file.
func calculateMD5(reader io.Reader, fileName string) {
	h := md5.New()
	if _, err := io.Copy(h, reader); err != nil {
		if err != io.EOF {
			log.Fatal(err)
		}
	}

	log.Printf("file: %s, md5: %x\n", fileName, h.Sum(nil))
}
