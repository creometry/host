package controller

import (
	"api/su"

	"github.com/gofiber/fiber/v2"
)

type Address struct {
	Number []string `json:"number"`
}

func Payment(c *fiber.Ctx) error {
	address := new(Address)
	if err := c.BodyParser(address); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	for _, v := range address.Number {
		su.Send(v)
	}
	return c.Status(500).SendString("good")

}
