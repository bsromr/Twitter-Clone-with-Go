package auth

import (
	"context"
	"net/http"
	"strconv"
	"time"

	database "github.com/bsromr/cloneTwitter/db"
	"github.com/bsromr/cloneTwitter/db/types"
	"github.com/gofiber/fiber/v2"

	jwt "github.com/form3tech-oss/jwt-go"
)

const SecretKey = "secret"

func LoginUser(c *fiber.Ctx) error {
	var users = types.Users{}
	if err := c.BodyParser(&users); err != nil {
		return err
	}
	db := database.DB
	err := db.QueryRow(context.Background(), "select email, password from users where email=$1, password=$2", users.Email, users.Password).Scan(&users.ID, &users.Name, &users.Created_at,&users.Updated_at, &users.Deleted_at, &users.Email,&users.Password, &users.Slug)
	if users.ID == 0 {
		return c.Render("login", fiber.Map{
			"Unauthorized": true,
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(users.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Couldn't login.",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.Redirect("home", http.StatusMovedPermanently)
}
