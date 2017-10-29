package main

/*
#include "brovi_cam.h"
*/
import "C"
import "errors"

var (
	//ErrOpenCamFail 打开摄像头失败
	ErrOpenCamFail = errors.New("failed to open camera")
	//ErrStartFail 启动视频失败
	ErrStartFail = errors.New("failed to start")
	//ErrStopFail 停止视频失败
	ErrStopFail = errors.New("failed to stop")
	//ErrGetNextBufferFail 获取下一帧失败
	ErrGetNextBufferFail = errors.New("failed to get next buffer")
)

//BroviCam handles the video capturing process
type BroviCam struct {
	broviCam *C.BroviCam
}

func newBroviCam(config *C.BroviCamConfig) (*BroviCam, error) {
	bc := new(BroviCam)
	if bc.broviCam = C.BroviCam_Open(config); bc.broviCam == nil {
		return nil, ErrOpenCamFail
	}
	return bc, nil
}

//Close closes the camera file and destroys underlying dependency
func (bc *BroviCam) Close() {
	C.BroviCam_Close(bc.broviCam)
}

//Start starts the video stream
func (bc *BroviCam) Start() error {
	if int(C.BroviCam_Start(bc.broviCam)) < 0 {
		return ErrStartCamFail
	}
	return nil
}

//Stop stops the video stream
func (bc *BroviCam) Stop() error {
	if int(C.BroviCam_Stop(bc.broviCam) < 0) {
		return ErrStopFail
	}
	return nil
}

//--------------------------------------------------

//BroviCamBuilder the builder to build new BroviCam
type BroviCamBuilder struct {
	config *C.BroviCamConfig
}

//NewBroviCam creates a new builder
func NewBroviCam(devfile string) *BroviCamBuilder {
	builder := new(BroviCamBuilder)
	builder.config = &C.BroviCamConfig{
		devfile: C.CString(devfile),
		width:   640, //DEFAULT VALUE
		height:  480, //DEFAULT VALUE
	}
	return builder
}

//Open builds the actual BroviCam object and start intializing process
func (builder *BroviCamBuilder) Open() (*BroviCam, error) {
	return newBroviCam(builder.config)
}

//SetWidth overwrites the width setting
func (builder *BroviCamBuilder) SetWidth(width int) *BroviCamBuilder {
	builder.config.width = C.int(width)
	return builder
}

//SetHeight overwrites the height setting
func (builder *BroviCamBuilder) SetHeight(height int) *BroviCamBuilder {
	builder.config.height = C.int(height)
	return builder
}
