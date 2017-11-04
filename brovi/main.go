package main

import (
	"fmt"
	"os"

	"github.com/Zumium/brovi/brovicam"
	"github.com/Zumium/brovi/brovicodec"
	"github.com/Zumium/brovi/cfg"
	"github.com/Zumium/brovi/server"
)

func reportErr(err error) {
	fmt.Fprintf(os.Stderr, "error occured: %s\n", err)
}

func main() {
	if err := cfg.Init(); err != nil {
		reportErr(err)
		os.Exit(1)
	}

	broviCam, err := brovicam.NewBuilder("/dev/video0").SetWidth(640).SetHeight(480).Open()
	if err != nil {
		reportErr(err)
		os.Exit(1)
	}
	defer broviCam.Close()
	defer broviCam.Stop()

	streamDuplicator := server.NewStreamDuplicator()

	broviCodec, err := brovicodec.NewBuilder(streamDuplicator.Inputer()).SetWidth(640).SetHeight(480).Build()
	if err != nil {
		reportErr(err)
		os.Exit(1)
	}
	defer broviCodec.Close()

	if server.Init(streamDuplicator); err != nil {
		reportErr(err)
		os.Exit(1)
	}

	broviCam.Pipe(broviCodec)
	if err := broviCam.Start(); err != nil {
		reportErr(err)
		os.Exit(1)
	}

	if err := server.Start(); err != nil {
		reportErr(err)
		return
	}
}
