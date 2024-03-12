package helpers

import (
	"log"
	"net/url"
	"os"

	pm "github.com/blakeofwilliam/alice-example/shared/peerManager"
	ws "github.com/blakeofwilliam/alice-example/shared/websocket"
	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
	fiberws "github.com/gofiber/websocket/v2"
)

func ConnectToServer(c *fiber.Ctx, curveId string) (peerManager *pm.PeerManager, err error) {
	u, err := url.Parse(os.Getenv("SERVER_HOST") + c.Path())
	if err != nil {
		log.Println("<<< Error getting server host: ", err)
		return
	}

	connection, _, err := websocket.DefaultDialer.Dial(u.String() + "?curve=" + curveId, nil)
	if err != nil {
		log.Println("<<< Error on Dial(): ", err.Error())
		return
	}

	log.Println("<<< Connected to server at: ", u.String())

	peerManager = pm.NewPeerManager("client", &ws.WebSocket{
		Socket: connection,
	}, "dkg")

	return
}

func ConnectToClient(socket *fiberws.Conn) (peerManager *pm.PeerManager, err error) {
	peerManager = pm.NewPeerManager("server", &ws.WebSocket{
		Socket: socket.Conn,
	}, "dkg")

	return
}