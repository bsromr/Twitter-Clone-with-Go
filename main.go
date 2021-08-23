package main

import (
	"log"

	"github.com/bsromr/cloneTwitter/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Panicln(err)
	}
}

func main() {
	db.Connect()
	tmpEngine := html.New("./views", ".html")
	tmpEngine.Reload(true)

	app := fiber.New(fiber.Config{
		Views: tmpEngine,
	})

	SetRoutes(app)

	app.Static("/", "./views")
	//app.Use(logger.New())
	app.Listen(":3000")
}
