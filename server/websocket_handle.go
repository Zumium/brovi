package server

import (
	"io"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

func liveStreamHandler(stream io.ReadCloser) echo.HandlerFunc {
	return func(c echo.Context) error {
		conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
		if err != nil {
			return err
		}

		go func() {
			defer conn.Close()

			buf := make([]byte, 4096)
			for {
				n, err := stream.Read(buf)
				if err != nil {
					return
				}
				if err := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
					return
				}
			}
		}()
	}
}
