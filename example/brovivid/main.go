package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Zumium/brovi/brovicam"
	"github.com/Zumium/brovi/brovicodec"
)

func reportError(err error) {
	fmt.Fprintf(os.Stderr, "error occurd: %s\n", err)
}

func main() {
	file, err := os.Create("record.h264")
	if err != nil {
		reportError(err)
		return
	}
	defer file.Close()

	codec, err := brovicodec.NewBuilder(func(d []byte) {
		fmt.Printf("%d bytes from codec\n", len(d))
		file.Write(d)
	}).Build()
	if err != nil {
		reportError(err)
		return
	}
	defer codec.Close()

	broviCam, err := brovicam.NewBuilder("/dev/video0").Open()
	if err != nil {
		reportError(err)
		return
	}
	defer broviCam.Close()

	if err := broviCam.Start(func(d []byte) {
		fmt.Printf("%d bytes from camera\n", len(d))
		codec.Write(d)
	}); err != nil {
		reportError(err)
		return
	}
	time.Sleep(10 * time.Second)
	if err := broviCam.Stop(); err != nil {
		reportError(err)
		return
	}
}
