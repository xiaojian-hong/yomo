package main

import (
	"crypto/md5"
	"encoding/json"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/yomorun/yomo"
	"github.com/yomorun/yomo/core"
	"github.com/yomorun/yomo/core/frame"
)

var sfnName string

// where files are stored
var dir string

var exit chan bool

func main() {
	if len(os.Args) < 2 {
		panic("please set the sfn name in args, f.e. go run main.go sfn-1")
	}
	sfnName = os.Args[1]
	if sfnName == "" {
		sfnName = "sink-1"
	}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dir = filepath.Dir(ex)

	// init yomo-sfn
	client := yomo.NewStreamFunction(
		sfnName,
		yomo.WithZipperAddr(os.Getenv("zipper")),
		yomo.WithObserveDataTags(0x11),
	)
	defer client.Close()

	// set stream handler (must set before connect)
	client.SetStreamHandler(streamHandler)

	// connect to yomo-zipper
	err = client.Connect()
	if err != nil {
		panic(err)
	}

	exit = make(chan bool)

	<-exit
	client.Close()
}

type fileInfo struct {
	Name string `json:"name"`
}

func streamHandler(in io.Reader) io.Reader {
	// get stream frame
	streamFrame, err := core.ParseFrame(in)
	if err != nil {
		log.Println("read frame error:", err)
		return nil
	}
	if streamFrame.Type() != frame.TagOfStreamFrame {
		log.Println("the frame type is not TagOfStreamFrame")
		return nil
	}

	// get file info from stream frame
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
	p := path.Join(dir, sfnName+"-"+info.Name)
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

	// signal to exit
	exit <- true
	return nil
}

// calculateMD5 calculates the md5 of the file.
func calculateMD5(reader io.Reader, fileName string) {
	h := md5.New()
	if _, err := io.Copy(h, reader); err != nil {
		log.Fatal(err)
	}

	log.Printf("file: %s, md5: %x\n", fileName, h.Sum(nil))
}
