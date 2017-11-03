package main

import (
	"fmt"
	"os"

	"github.com/Zumium/brovi/brovicam"
)

func reportError(err error) {
	fmt.Fprintf(os.Stderr, "error occurd: %s\n", err)
	os.Exit(1)
}

func main() {
	file, err := os.Create(os.Args[2])
	if err != nil {
		reportError(err)
	}
	defer file.Close()

	bc, err := brovicam.NewBuilder(os.Args[1]).Open()
	if err != nil {
		reportError(err)
	}
	defer bc.Close()

	if err := bc.OneFrame(func(frame []byte) {
		n, _ := file.Write(frame)
		fmt.Printf("write %d bytes\n", n)
	}); err != nil {
		reportError(err)
	}
}
