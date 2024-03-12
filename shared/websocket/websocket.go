package websocket

import (
	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
)

type WebSocket struct {
	Socket *websocket.Conn
}

var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize: 4096,
	WriteBufferSize: 4096,
}

func (websocket *WebSocket) Close() (err error) {
	err = websocket.Socket.Close()

	return
}

func (socket *WebSocket) Write(message []byte) (socketError error) {
	errorChannel := make(chan error)

	err := socket.Socket.WriteMessage(websocket.BinaryMessage, message)
	if err != nil {
		errorChannel <- err
	}

	socketError = <-errorChannel

	return
}

func Upgrade(c *fiber.Ctx) (socket *websocket.Conn, err error) {
	ws := make(chan *websocket.Conn)

	err = upgrader.Upgrade(c.Context(), func(conn *websocket.Conn) {
		println("Websocket upgraded.", conn)
		ws <- conn
	})
	if err != nil {
		println("Error upgrading websocket connection: ", err.Error())
		return
	}

	socket = <-ws

	return
}