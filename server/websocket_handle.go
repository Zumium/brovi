package server

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

func liveStreamHandler(streamDuplicator *StreamDuplicator) echo.HandlerFunc {
	return func(c echo.Context) error {
		stream := streamDuplicator.NewOutput()
		conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
		if err != nil {
			return err
		}

		exitSig := make(chan struct{})

		//read loop to process ping, pong and close messages
		go func() {
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					exitSig <- struct{}{}
					<-exitSig
					close(exitSig)
					conn.Close()
				}
			}
		}()

		//writing loop to send video stream data
		go func() {
			defer stream.Close()

			buf := make([]byte, 4096)
			for {
				select {
				case <-exitSig:
					stream.Close()
					exitSig <- struct{}{}
					return
				default:
					//failover through
				}
				n, err := stream.Read(buf)
				if err != nil {
					return
				}
				if err := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
					return
				}
			}
		}()
		return nil
	}
}
