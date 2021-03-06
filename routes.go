package main

import (
	"github.com/bsromr/cloneTwitter/controller/auth"
	"github.com/gofiber/fiber/v2"
)

func SetRoutes(app *fiber.App) {
	/*GETS*/
	app.Get("/", HomePage)
	app.Get("login", Login)
	app.Get("signup", Signup)
	app.Get("home", Home)
	app.Get("/:searchedUser", Profile)
	app.Get("/:searchedUser/likeTweett/:likedTweetID", LikeTweet)
	app.Get("/:searchedUser/status/:tweetID", MentionTweet)
	/*POSTS*/
	app.Post("registeruser", auth.RegisterUser)
	app.Post("loginUser", auth.LoginUser)
	app.Post("tweet", Tweet)
	app.Post("follow/:uid", Follow)
	app.Post("logout", auth.Logout)

}
