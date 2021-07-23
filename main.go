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
	app := fiber.New(
		fiber.Config{
			Views: tmpEngine,
		})
	app.Static("/", "./views")
	//app.Use(logger.New())
	SetRoutes(app)
	app.Listen(":3000")
}
