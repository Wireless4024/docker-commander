package main

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"os"
	"os/exec"
)

func main() {
	_ = godotenv.Load()

	app := fiber.New()

	command := os.Getenv("COMMAND")
	if len(command) == 0 {
		panic("COMMAND env is missing")
	}

	token := os.Getenv("TOKEN")
	if len(token) == 0 {
		panic("TOKEN env is missing")
	}

	app.Get("/", func(ctx *fiber.Ctx) error {
		reqToken := ctx.Query("token")
		if reqToken != token {
			if err := ctx.SendStatus(400); err != nil {
				return err
			}
			return errors.New("invalid token")
		}
		cmd := exec.Command("bash", "-c", command)

		if err := cmd.Run(); err != nil {
			return err
		}
		return ctx.SendStatus(200)
	})

	listen := os.Getenv("LISTEN")
	if len(listen) == 0 {
		listen = ":3000"
	}
	_ = app.Listen(listen)
}
