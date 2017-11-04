package server

import (
	"bytes"
	"io"
	"testing"
)

func TestStreamDuplicator(t *testing.T) {
	duplicator := NewStreamDuplicator()
	inputFn := duplicator.Inputer()
	out1 := duplicator.NewOutput()
	out2 := duplicator.NewOutput()
	out3 := duplicator.NewOutput()
	go inputFn([]byte("hello world"))
	check := func(out io.ReadCloser, i int) {
		buf := make([]byte, 20)
		n, err := out.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(buf[:n], []byte("hello world")) {
			t.Fatalf("not match for %d\n", i)
		}
		out.Close()
	}
	go check(out1, 1)
	go check(out2, 2)
	go check(out3, 3)
}
