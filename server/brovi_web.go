package server

import (
	"fmt"

	"github.com/Zumium/brovi/cfg"
	"github.com/labstack/echo"
)

const defaultListenAddr = "0.0.0.0"

var e = echo.New()

//Init initializes the web service
func Init(duplicator *StreamDuplicator) error {
	e.GET("/live", liveStreamHandler(duplicator))
	return nil
}

//Start starts the web service then blocks
func Start() error {
	return e.Start(fmt.Sprintf("%s:%d", defaultListenAddr, cfg.Port()))
}
