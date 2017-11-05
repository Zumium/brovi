package server

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Zumium/brovi/cfg"
	"github.com/labstack/echo"
)

const defaultListenAddr = "0.0.0.0"

var e = echo.New()

//Init initializes the web service
func Init(duplicator *StreamDuplicator) error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}
	absExecPath, err := filepath.Abs(execPath)
	if err != nil {
		return err
	}
	e.Static("/", filepath.Join(absExecPath, "static"))
	e.GET("/live", liveStreamHandler(duplicator))
	return nil
}

//Start starts the web service then blocks
func Start() error {
	return e.Start(fmt.Sprintf("%s:%d", defaultListenAddr, cfg.Port()))
}
