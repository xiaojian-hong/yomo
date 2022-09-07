package main

import (
	"crypto/md5"
	"encoding/json"
	"io"
	"log"
	"os"
	"path"

	"github.com/yomorun/yomo"
	"github.com/yomorun/yomo/core"
	"github.com/yomorun/yomo/core/frame"
)

var sfnName string

func main() {
	if len(os.Args) < 2 {
		panic("please set the sfn name in args, f.e. go run main.go sfn-1")
	}
	sfnName = os.Args[1]
	if sfnName == "" {
		sfnName = "sink-1"
	}
	// init yomo-source
	client := yomo.NewStreamFunction(
		sfnName,
		yomo.WithZipperAddr("localhost:9000"),
		yomo.WithObserveDataTags(0x11),
	)
	defer client.Close()

	// set stream handler (must set before connect)
	client.SetStreamHandler(streamHandler)

	// connect to yomo-zipper
	err := client.Connect()
	if err != nil {
		panic(err)
	}

	select {}
}

type fileInfo struct {
	Name string `json:"name"`
	Dir  string `json:"dir"`
}

func streamHandler(in io.Reader) io.Reader {
	streamFrame, err := core.ParseFrame(in)
	if err != nil {
		log.Println("read frame error:", err)
		return nil
	}
	if streamFrame.Type() != frame.TagOfStreamFrame {
		log.Println("the frame type is not TagOfStreamFrame")
		return nil
	}

	var info fileInfo
	err = json.Unmarshal(streamFrame.(*frame.StreamFrame).Metadata(), &info)
	if err != nil {
		log.Println("unmarshal the metadata in stream frame error:", err)
		return nil
	}

	log.Println("receiving file:", info.Name)

	// calculate the md5 of the file
	pipeReader, pipeWriter := io.Pipe()

	go calculateMD5(pipeReader, info.Name)

	// create output file
	p := path.Join(info.Dir, sfnName+"-"+info.Name)
	f, err := os.Create(p)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	written, err := io.Copy(io.MultiWriter(pipeWriter, f), in)
	if err != nil && err != io.EOF {
		panic(err)
	}
	pipeWriter.Close()
	log.Printf("written: %d, %s\n", written, p)
	return nil
}

// calu;ateMD5 calculates the md5 of the file.
func calculateMD5(reader io.Reader, fileName string) {
	h := md5.New()
	if _, err := io.Copy(h, reader); err != nil {
		log.Fatal(err)
	}

	log.Printf("file: %s, md5: %x\n", fileName, h.Sum(nil))
}
