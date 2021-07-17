package auth

import (
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	database "github.com/bsromr/cloneTwitter/db"
	"github.com/bsromr/cloneTwitter/db/types"
	"github.com/gofiber/fiber/v2"
)

func RegisterUser(c *fiber.Ctx) error {
	if c.FormValue("name") == "" || c.FormValue("phone") == "" || c.FormValue("email") == "" || c.FormValue("password") == "" {
		return c.Redirect("/signup")
	}
	db := database.DB
	users := new(types.Users)

	//check if the user email exist on db.
	db.Where("email = ? ", c.FormValue("email")).Find(&users)
	if users.ID != 0 {
		return c.Render("signup", fiber.Map{
			"CantRegister": true,
			"Email":        c.FormValue("email"),
		})
	}

	//Parsing body and creating user on db.
	if err := c.BodyParser(users); err != nil {
		return err
	}
	rand.Seed(time.Now().UnixNano())
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	users.Slug = reg.ReplaceAllString(strings.ToLower(strings.ReplaceAll(users.Name, " ", "")+strconv.Itoa(rand.Intn(10000))), "")

	db.Create(&users)
	return c.Render("signup", fiber.Map{
		"Situation": true,
	})
}
