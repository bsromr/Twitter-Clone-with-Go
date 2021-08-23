package main

import (
	"github.com/bsromr/cloneTwitter/controller/auth"
	"github.com/gofiber/fiber/v2"
)

func SetRoutes(app *fiber.App) {
	/*GETS*/
	app.Get("/", HomePage)
	app.Get("login", LoginPage)
	app.Get("signup", SignupPage)
	app.Get("home", Home)
	app.Get("/:searchedUser", ProfilePage)
	app.Get("/:searchedUser/likeTweett/:likedTweetID", LikeTweet)
	app.Get("/:searchedUser/status/:tweetID", MentionPage)
	app.Get("/:searchedUser/status/:tweetID/mention", MentionTweet)
	/*POSTS*/
	app.Post("registeruser", auth.RegisterUser)
	app.Post("loginUser", auth.LoginUser)
	app.Post("tweet", Tweet)
	app.Post("follow/:uid", Follow)
	app.Post("logout", auth.Logout)

}
