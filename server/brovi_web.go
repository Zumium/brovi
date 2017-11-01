package server

import (
	"fmt"
	"io"

	"github.com/Zumium/brovi/cfg"
	"github.com/labstack/echo"
)

const defaultListenAddr = "0.0.0.0"

var e = echo.New()

//Init initializes the web service
func Init(stream io.ReadCloser) error {
	e.GET("/live", liveStreamHandler(stream))
	return nil
}

//Start starts the web service then blocks
func Start() error {
	return e.Start(fmt.Sprintf("%s:%d", defaultListenAddr, cfg.Port()))
}
