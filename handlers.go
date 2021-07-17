package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/bsromr/cloneTwitter/controller/auth"
	"github.com/bsromr/cloneTwitter/db"
	_ "github.com/bsromr/cloneTwitter/db"
	database "github.com/bsromr/cloneTwitter/db"
	"github.com/bsromr/cloneTwitter/db/types"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func HomePage(c *fiber.Ctx) error {
	return c.Render("/", fiber.Map{
		"Title": "Welcome to Twitter",
	})
}

func Login(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{
		"Title": "Login to Twitter",
	})
}

func Signup(c *fiber.Ctx) error {
	return c.Render("signup", fiber.Map{
		"Title": "Signup to Twitter",
	})
}

func Home(c *fiber.Ctx) error {
	user := getUserIdFromCookie(c)

	return c.Render("home", fiber.Map{
		"UserName":  user.Name,
		"UserMail":  user.Email,
		"UserPhone": user.Phone,
		"Slug":      user.Slug,
	})
}

func Tweet(c *fiber.Ctx) error {
	tweet := types.Tweets{}
	if err := c.BodyParser(&tweet); err != nil {
		return err
	}
	db := database.DB
	user := getUserIdFromCookie(c)

	tweet.User_id = int(user.ID)
	//log.Println("Tweet: ", tweet.Tweet, "User Id: ", user.ID)
	db.Create(&tweet)
	return c.Redirect("home", fiber.StatusMovedPermanently)
}

func Profile(c *fiber.Ctx) error {
	db := db.DB
	//log.Println(utils.ImmutableString(c.Params("foo")))

	//loggedInUser := getUserIdFromCookie(c)

	var tweets = []types.Tweets{} //Tweets
	result := db.Raw("Select users.name, users.slug,tweets.id, tweets.tweet, tweets.created_at, tweets.like_count FROM tweets Left JOIN users ON users.id = tweets.user_id Left JOIN tweet_infos ON tweet_infos.tweet_id = tweets.id WHERE users.slug = '" + c.Params("searchedUser") + "' group by (users.name, users.slug, tweets.tweet,tweets.id, tweets.created_at, tweets.like_count)  order by tweets.created_at DESC;").Scan(&tweets)
	totalTweets := result.RowsAffected

	searchedUserInfo := types.Users{} //User Info belonging to the slug
	db.Where("slug = ?", c.Params("searchedUser")).First(&searchedUserInfo)

	var whoToFollowUsers = []types.Users{} //Who to follow
	db.Raw("SELECT * from users where id <> ? order by RANDOM() limit 3", searchedUserInfo.ID).Scan(&whoToFollowUsers)

	var followingCount int64
	db.Raw("SELECT COUNT(*) as count FROM relationships JOIN users ON relationships.follower_id = users.id WHERE users.slug = '" + c.Params("searchedUser") + "' ").Count(&followingCount)

	var followerCount int64
	db.Raw("SELECT COUNT(*) FROM relationships JOIN users ON relationships.followed_id = users.id WHERE users.slug = '" + c.Params("searchedUser") + "' ").Count(&followerCount)

	return c.Render("profile", fiber.Map{
		"Username":       searchedUserInfo.Name,
		"Slug":           searchedUserInfo.Slug,
		"CreatedAt":      searchedUserInfo.CreatedAt,
		"UserTweets":     tweets,
		"TotalTweets":    totalTweets,
		"WhoToFollow":    whoToFollowUsers,
		"FollowingCount": followingCount,
		"FollowerCount":  followerCount,
		//TODO: LIKED BUTONU KIRMIZI YAP
	})
}

func LikeTweet(c *fiber.Ctx) error {
	fmt.Println("searched user ID: ", c.Params("searchedUser"), "seçilen tweet id: ", c.Params("likedTweetID"))
	db := db.DB

	onlineUser := getUserIdFromCookie(c)
	//	fmt.Println("Aktif kullanıcı id: ", onlineUser.ID)

	tweet_infos := types.Tweet_Info{}
	tweet_infos.Tweet_id, _ = strconv.Atoi(c.Params("likedTweetID"))
	tweet_infos.Liked_user_id = int(onlineUser.ID)

	result := db.Where("tweet_id = ? and liked_user_id = ? ", c.Params("likedTweetID"), onlineUser.ID).First(&tweet_infos)
	if result.RowsAffected > 0 {
		fmt.Println("Silinme işlemi uygulanacak")
		db.Where("tweet_id = ? and liked_user_id = ?", c.Params("likedTweetID"), onlineUser.ID).Delete(&tweet_infos)
		db.Exec("UPDATE tweets SET like_count = tweets.like_count - 1 where id = ?", c.Params("likedTweetID")) //decrease count of liked tweet
		return c.Redirect("/" + c.Params("searchedUser"))
	}

	db.Create(&tweet_infos)
	db.Exec("UPDATE tweets SET like_count = tweets.like_count + 1 where id = ?", c.Params("likedTweetID")) //increase count of liked tweet
	return c.Redirect("/" + c.Params("searchedUser"))
}

func Follow(c *fiber.Ctx) error {
	//log.Println("UUUID= ", utils.ImmutableString(c.Params("uid")))
	var err error
	db := db.DB
	follow := types.Relationships{}

	activeUser := getUserIdFromCookie(c)

	follow.FollowedId, err = strconv.Atoi(c.Params("uid"))
	if err != nil {
		log.Fatal(err)
	}
	follow.FollowerId = int(activeUser.ID)

	result := db.Where("follower_id = ? AND followed_id = ?", activeUser.ID, follow.FollowedId).First(&follow)
	if result.RowsAffected > 0 {
		return c.Redirect("/"+activeUser.Slug+"", fiber.StatusMovedPermanently)
	}

	db.Create(&follow)

	return c.Redirect("/"+activeUser.Slug+"", fiber.StatusMovedPermanently)
}

func getUserIdFromCookie(c *fiber.Ctx) types.Users {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(auth.SecretKey), nil
	})
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		c.Redirect("login", fiber.StatusMovedPermanently)
		return types.Users{}
	}
	claims := token.Claims.(*jwt.StandardClaims)
	var user types.Users
	db := database.DB
	db.Where("id = ?", claims.Issuer).First(&user)
	if user.ID == 0 {
		log.Println("siktir çekildi")
		c.Redirect("login", fiber.StatusMovedPermanently)
		return types.Users{}
	}

	return user
}
