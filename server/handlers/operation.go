package handlers

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/blakeofwilliam/alice-example/shared/helpers"
	"github.com/blakeofwilliam/alice-example/shared/operations"
	"github.com/gofiber/websocket/v2"
)

var Operation = websocket.New(func(c *websocket.Conn) {
	operation := c.Params("operation")

	if operation == "generate" {
		err := generate(c)
		if err != nil {
			log.Println("<<< Error on generate: ", err.Error())
			c.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
		}
	}
})

func generate(c *websocket.Conn) (err error) {
	peerManager, err := helpers.ConnectToClient(c)
	if err != nil {
		log.Println("<<< Error creating peer manager: ", err.Error())
		return
	}
	
	curveID := strings.ToUpper(c.Query("curve", "secp256k1"))
	service, err := operations.NewGenerateService(curveID, peerManager)
	if err != nil {
		log.Println("<<< Error creating DKG service: ", err.Error())
		return
	}

	result, err := service.Process()
	if err != nil {
		return
	}

	log.Println("<<< Result: ", result, " Error: ", err)

	response, err := json.Marshal(result)
	if err != nil {
		return
	}

	c.WriteMessage(websocket.BinaryMessage, response)

	return
}
