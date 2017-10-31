package brovicodec

/*
#cgo CPPFLAGS: -std=c++11
#cgo LDFLAGS: -lx264

#include <stdlib.h>
#include "brovi_codec.h"
*/
import "C"
import "unsafe"
import "io"
import "errors"

//ErrInitCodecFail -- failed to init codec
var ErrInitCodecFail = errors.New("failed to init codec")

//EncodedFrameCallback handles the returned encoded frame data
type EncodedFrameCallback func([]byte)

//BroviCodec -- the codec
type BroviCodec struct {
	codec   *C.BroviCodec
	cb      EncodedFrameCallback
	exitSig chan struct{}
	in      chan []byte
	alive   bool
}

func newCodec(config C.BroviCodecConfig, cb EncodedFrameCallback) (*BroviCodec, error) {
	codec := &BroviCodec{
		in:      make(chan []byte),
		exitSig: make(chan struct{}),
		cb:      cb,
		codec:   (*C.BroviCodec)(C.BroviCodec_New(config)),
		alive:   true,
	}
	if codec.codec == nil {
		return nil, ErrInitCodecFail
	}
	go codec.process()
	return codec, nil
}

//Write implements io.Writer interface
func (c *BroviCodec) Write(p []byte) (n int, err error) {
	if !c.alive {
		return 0, io.ErrClosedPipe
	}
	c.in <- p
	return len(p), nil
}

//Close implements the io.Closer interface
func (c *BroviCodec) Close() (err error) {
	c.exitSig <- struct{}{}
	<-c.exitSig
	close(c.in)
	close(c.exitSig)
	C.BroviCodec_Close(unsafe.Pointer(c.codec))
	return
}

func (c *BroviCodec) process() {
	ctnu := true
	for ctnu {
		var out C.H264Frame
		select {
		case <-c.exitSig:
			ctnu = false
		case in := <-c.in:
			inData := C.CBytes(in)
			ret := C.BroviCodec_EncodeFrame(unsafe.Pointer(c.codec), inData, &out)
			C.free(inData)
			if ret == C.BROVI_CODEC_ZERO_SIZE_ERR {
				continue
			}
			c.cb(C.GoBytes(out.data, out.size))
		}
	}
	for {
		var out C.H264Frame
		if ret := C.BroviCodec_FlushDelayedFrame(unsafe.Pointer(c.codec), &out); ret == 0 {
			c.exitSig <- struct{}{}
			return
		}
		c.cb(C.GoBytes(out.data, out.size))
	}
}

//-----------------------------------------------------------------

//Builder is used to build a BroviCodec object
type Builder struct {
	config C.BroviCodecConfig
	cb     EncodedFrameCallback
}

//New creates a new builder to help construct a new codec
func New(cb EncodedFrameCallback) *Builder {
	return &Builder{cb: cb}
}

//SetWidth sets the frame's width
func (b *Builder) SetWidth(width int) *Builder {
	b.config.width = C.int(width)
	return b
}

//SetHeight sets the frame's width
func (b *Builder) SetHeight(height int) *Builder {
	b.config.height = C.int(height)
	return b
}

//Build finishes the setting and build a BroviCodec instance
func (b *Builder) Build() (*BroviCodec, error) {
	return newCodec(b.config, b.cb)
}
