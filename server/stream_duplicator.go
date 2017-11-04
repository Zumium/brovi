package server

import (
	"io"
	"sync"
)

//InputCallback is a function which is to be called to input data
type InputCallback func([]byte)

//StreamDuplicator duplicates data from single source to multiple outputs
type StreamDuplicator struct {
	in   chan []byte
	outs sync.Map

	exitSig chan struct{}
}

//NewStreamDuplicator creates a new StreamDuplicator
func NewStreamDuplicator() *StreamDuplicator {
	duplicator := &StreamDuplicator{
		in:      make(chan []byte, 1),
		exitSig: make(chan struct{}),
	}
	go duplicator.process()
	return duplicator
}

//Inputer returns the input callback function
func (dup *StreamDuplicator) Inputer() InputCallback {
	return func(d []byte) {
		dup.in <- d
	}
}

//NewOutput returns a new output as ReadCloser
func (dup *StreamDuplicator) NewOutput() io.ReadCloser {
	output := newStreamDuplicatorOutput(dup)
	dup.outs.Store(output, output.c)
	return output
}

//Close closes all outputs
func (dup *StreamDuplicator) Close() error {
	close(dup.in)
	dup.exitSig <- struct{}{}
	close(dup.exitSig)
	keys := make([]*streamDuplicatorOutput, 0)
	dup.outs.Range(func(key, value interface{}) bool {
		k := key.(*streamDuplicatorOutput)
		keys = append(keys, k)
		return true
	})
	for _, key := range keys {
		key.Close()
	}
	return nil
}

func (dup *StreamDuplicator) process() {
	for {
		select {
		case <-dup.exitSig:
			return
		default:
			//fallover
		}
		d := <-dup.in
		dup.outs.Range(func(key interface{}, value interface{}) bool {
			c := value.(chan []byte)
			c <- d
			return true
		})
	}
}

//----------------------------------------

type streamDuplicatorOutput struct {
	duplicator *StreamDuplicator
	c          chan []byte
}

func newStreamDuplicatorOutput(duplicator *StreamDuplicator) *streamDuplicatorOutput {
	return &streamDuplicatorOutput{
		duplicator: duplicator,
		c:          make(chan []byte, 1),
	}
}

//Read implements io.ReadCloser interface
func (output *streamDuplicatorOutput) Read(p []byte) (int, error) {
	d := <-output.c
	len := copy(p, d)
	return len, nil
}

//Close implements io.ReadCloser interface
func (output *streamDuplicatorOutput) Close() error {
	close(output.c)
	output.duplicator.outs.Delete(output)
	return nil
}
