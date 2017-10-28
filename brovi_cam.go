package main

/*
#include "brovi_cam.h"
*/
import "C"

//BroviCam handles the video capturing process
type BroviCam struct {
	broviCam *C.BroviCam
}

func newBroviCam(config *C.BroviCamConfig) *BroviCam {
	bc := new(BroviCam)
	bc.broviCam = C.BroviCam_Open(config)
	return bc
}

//Close closes the camera file and destroys underlying dependency
func (bc *BroviCam) Close() {
	C.BroviCam_Close(bc.broviCam)
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
func (builder *BroviCamBuilder) Open() *BroviCam {
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
