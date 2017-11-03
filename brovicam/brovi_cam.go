package brovicam

/*
#cgo CPPFLAGS: -std=c++11

#include "brovi_cam.h"
*/
import "C"
import "errors"
import "unsafe"

var (
	//ErrOpenCamFail 打开摄像头失败
	ErrOpenCamFail = errors.New("failed to open camera")
	//ErrStartFail 启动视频失败
	ErrStartFail = errors.New("failed to start")
	//ErrStopFail 停止视频失败
	ErrStopFail = errors.New("failed to stop")
	//ErrGetNextBufferFail 获取下一帧失败
	ErrGetNextBufferFail = errors.New("failed to get next buffer")
	//ErrEmptyCallBack 回调函数为设置
	ErrEmptyCallBack = errors.New("callback function has been left unset")
)

//FrameCallback deal with frame bytes
type FrameCallback func([]byte)

//BroviCam handles the video capturing process
type BroviCam struct {
	exitSig  chan struct{}
	cb       FrameCallback
	broviCam *C.BroviCam
}

func newBroviCam(config *C.BroviCamConfig) (*BroviCam, error) {
	bc := &BroviCam{exitSig: make(chan struct{})}
	if bc.broviCam = (*C.BroviCam)(C.BroviCam_Open(config)); bc.broviCam == nil {
		return nil, ErrOpenCamFail
	}
	return bc, nil
}

//Close closes the camera file and destroys underlying dependency
func (bc *BroviCam) Close() {
	close(bc.exitSig)
	C.BroviCam_Close(unsafe.Pointer(bc.broviCam))
}

//Start starts the video stream
func (bc *BroviCam) Start(cb FrameCallback) error {
	bc.cb = cb
	return bc.start()
}

//Stop stops the video stream
func (bc *BroviCam) Stop() error {
	bc.exitSig <- struct{}{}
	if int(C.BroviCam_Stop(unsafe.Pointer(bc.broviCam))) < 0 {
		return ErrStopFail
	}
	return nil
}

//OneFrame is FOR TEST ONLY!!!
func (bc *BroviCam) OneFrame(cb FrameCallback) error {
	if int(C.BroviCam_Start(unsafe.Pointer(bc.broviCam))) < 0 {
		return ErrStartFail

	}
	status := C.BroviCam_NextBufferA(unsafe.Pointer(bc.broviCam))
	if status.buffer == nil {
		return ErrGetNextBufferFail
		// return errors.New("failed to dequeue buffer")
	}
	cb(C.GoBytes(status.buffer.start, C.int(status.buffer.length)))
	C.BroviCam_NextBufferB(unsafe.Pointer(bc.broviCam), status)
	if int(C.BroviCam_Stop(unsafe.Pointer(bc.broviCam))) < 0 {
		return ErrStopFail
	}
	return nil
}

func (bc *BroviCam) start() error {
	if bc.cb == nil {
		return ErrEmptyCallBack
	}
	if int(C.BroviCam_Start(unsafe.Pointer(bc.broviCam))) < 0 {
		return ErrStartFail
	}
	go bc.stream()
	return nil
}

func (bc *BroviCam) stream() {
	for {
		select {
		case <-bc.exitSig:
			return
		default:
			//fall throught
		}

		status := C.BroviCam_NextBufferA(unsafe.Pointer(bc.broviCam))
		bc.cb(C.GoBytes(status.buffer.start, C.int(status.buffer.length)))
		C.BroviCam_NextBufferB(unsafe.Pointer(bc.broviCam), status)
	}
}

//--------------------------------------------------

//Builder the builder to build new BroviCam
type Builder struct {
	config *C.BroviCamConfig
}

//NewBuilder creates a new builder
func NewBuilder(devfile string) *Builder {
	builder := new(Builder)
	builder.config = &C.BroviCamConfig{
		devfile: C.CString(devfile),
		width:   640, //DEFAULT VALUE
		height:  480, //DEFAULT VALUE
	}
	return builder
}

//Open builds the actual BroviCam object and start intializing process
func (builder *Builder) Open() (*BroviCam, error) {
	return newBroviCam(builder.config)
}

//SetWidth overwrites the width setting
func (builder *Builder) SetWidth(width int) *Builder {
	builder.config.width = C.int(width)
	return builder
}

//SetHeight overwrites the height setting
func (builder *Builder) SetHeight(height int) *Builder {
	builder.config.height = C.int(height)
	return builder
}
