package handlers

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/blakeofwilliam/alice-example/shared/helpers"
	"github.com/blakeofwilliam/alice-example/shared/operations"
	"github.com/gofiber/fiber/v2"
)

func Generate(c *fiber.Ctx) (err error) {
	curveID := strings.ToUpper(c.Query("curve", "secp256k1"))

	peerManager, err := helpers.ConnectToServer(c, curveID)
	if err != nil {
		return
	}

	service, err := operations.NewGenerateService(curveID, peerManager)
	if err != nil {
		return
	}

	result, err := service.Process()
	if err != nil {
		return c.Status(500).Send([]byte(err.Error()))
	}

	log.Println("<<< Result: ", result, " Error: ", err)

	response, err := json.Marshal(result)
	if err != nil {
		return c.Status(500).Send([]byte(err.Error()))
	}

	c.Send(response)

	return
}
