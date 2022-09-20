package main

import (
	"crypto/md5"
	"encoding/json"
	"io"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

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

	// init yomo-sfn
	client := yomo.NewStreamFunction(
		sfnName,
		yomo.WithZipperAddr(os.Getenv("zipper")),
		yomo.WithObserveDataTags(0x11),
	)

	// set stream handler (must set before connect)
	client.SetStreamHandler(streamHandler)

	// connect to yomo-zipper
	err := client.Connect()
	if err != nil {
		panic(err)
	}

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel

	client.Close()
}

type fileInfo struct {
	Name string `json:"name"`
}

func streamHandler(in io.Reader) io.Reader {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dir := filepath.Dir(ex)

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

	// create output file
	p := path.Join(dir, sfnName+"-"+strconv.Itoa(int(time.Now().Unix()))+"-"+info.Name)
	f, err := os.Create(p)
	if err != nil {
		panic(err)
	}

	written, err := io.Copy(f, in)
	if err != nil {
		if err == io.EOF {
			log.Println(">>EOF")
		} else {
			panic(err)
		}
	}
	log.Printf("written: %d, %s\n", written, p)

	calculateMD5(p)

	f.Close()

	return nil
}

// calculateMD5 calculates the md5 of the file.
func calculateMD5(fileName string) {
	vs, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer vs.Close()

	h := md5.New()
	if _, err := io.Copy(h, vs); err != nil {
		log.Fatal(err)
	}

	log.Printf("file: %s, md5: %x\n", fileName, h.Sum(nil))
	h.Reset()
}
