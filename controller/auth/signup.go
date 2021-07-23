package auth

import (
	"context"
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
	err := db.QueryRow(context.Background(), "select id, email from users where email=$1", c.FormValue("email")).Scan(&users.ID, &users.Email)
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
	_, err = db.Exec(context.Background(),"INSERT INTO users(created_at, updated_at, name, email, phone, password, slug) values ($1,$2,$3,$4,$5,$6,$7)", time.Now(), time.Now(), users.Name, users.Email,users.Phone,users.Password, users.Slug)
	if err != nil {
		panic(err)
	}
	return c.Render("signup", fiber.Map{
		"Situation": true,
	})
}
