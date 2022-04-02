package main

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"os"
	"os/exec"
	"strings"
)

func main() {
	_ = godotenv.Load()

	app := fiber.New()

	images := make(map[string]string)

	for _, kv := range os.Environ() {
		arr := strings.Split(kv, "=")
		if len(arr) > 0 {
			key := strings.TrimSpace(arr[0])
			value := strings.TrimSpace(arr[1])

			if !strings.HasPrefix(key, "IMAGE_") || len(key) == 0 || len(value) == 0 {
				continue
			}

			images[key[6:]] = value
		}
	}

	/*j, _ := json.Marshal(images)
	println(string(j))*/

	token := os.Getenv("TOKEN")
	if len(token) == 0 {
		panic("TOKEN env is missing")
	}

	trigger := func(ctx *fiber.Ctx) error {
		reqToken := ctx.Query("token")
		if len(reqToken) == 0 {
			reqToken = ctx.GetReqHeaders()["X-Gitlab-Token"]
		}
		if reqToken != token {
			if err := ctx.SendStatus(400); err != nil {
				return err
			}
			return errors.New("invalid token")
		}

		command := images[ctx.Params("image")]
		if len(command) == 0 {
			return errors.New("invalid image name")
		}

		cmd := exec.Command("bash", "-c", command)

		if err := cmd.Run(); err != nil {
			return err
		}
		return ctx.SendStatus(200)
	}

	app.Get("/:image", trigger)
	app.Post("/:image", trigger)

	listen := os.Getenv("LISTEN")
	if len(listen) == 0 {
		listen = ":3000"
	}
	_ = app.Listen(listen)
}
