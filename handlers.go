package main

import (
	"context"
	"github.com/bsromr/cloneTwitter/controller/auth"
	_ "github.com/bsromr/cloneTwitter/db"
	database "github.com/bsromr/cloneTwitter/db"
	"github.com/bsromr/cloneTwitter/db/types"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
	"time"
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
	activeUser := getUserIdFromCookie(c)

	tweet.User_id = int(activeUser.ID)
	//log.Println("Tweet: ", tweet.Tweet, "User Id: ", user.ID)
	//db.Create(&tweet)
	_, err := db.Exec(context.Background(),"INSERT INTO tweets(created_at, updated_at, user_id, tweet, like_count) values ($1,$2,$3,$4,$5)", time.Now(),time.Now(), activeUser.ID ,c.FormValue("tweet"), 0)
	if err != nil{
		log.Fatal(err)
	}
	return c.Redirect("home")
}

func Profile(c *fiber.Ctx) error {
	db := database.DB
	//log.Println(utils.ImmutableString(c.Params("foo")))
	var countTweets int
	loggedInUser := getUserIdFromCookie(c)
	var user = types.Users{}
	var tweet = types.Tweets{} //Tweets
	var tweets []types.Tweets
	//var tweet_info = types.Tweet_Info{}

	rows, err := db.Query(context.Background(), "Select users.name, users.slug, tweets.id, tweets.tweet, tweets.created_at, tweets.liked_user_id FROM tweets LEFT JOIN users ON users.id = tweets.user_id WHERE users.slug = $1 order by tweets.created_at DESC;", c.Params("searchedUser"))
	if err != nil{
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next(){
		if err := rows.Scan(&user.Name, &user.Slug, &tweet.ID, &tweet.Tweet, &tweet.Created_at, &tweet.Liked_user_id); err != nil{
			log.Fatal(err)
		}
		countTweets++
		tweet.LikeCount = len(tweet.Liked_user_id)
		tweets = append(tweets, tweet)
	}
	searchedUser := types.Users{} //User Info belonging to the slug
	db.QueryRow(context.Background(),"select id,created_at,name,slug from users where slug = $1", c.Params("searchedUser")).Scan(&searchedUser.ID,&searchedUser.Created_at,&searchedUser.Name,&searchedUser.Slug)

	//WHO TO FOLLOW
	var whoToFollowUser types.Users //Who to follow
	var whoToFollowUsers []types.Users
	wtfUserRows, err := db.Query(context.Background(),"SELECT name,slug from users where id <> $1 order by RANDOM() limit 3", searchedUser.ID)
	if err != nil{
		log.Fatal(err)
	}
	defer wtfUserRows.Close()
	for wtfUserRows.Next(){
		if err := wtfUserRows.Scan(&whoToFollowUser.Name, &whoToFollowUser.Slug); err != nil{
			log.Fatal(err)
		}
		whoToFollowUsers = append(whoToFollowUsers, whoToFollowUser)
	}

	//Following/Follower Count
	var followingCount int64
	db.QueryRow(context.Background(),"SELECT COUNT(*) as count FROM relationships JOIN users ON relationships.follower_id = users.id WHERE users.slug = $1", c.Params("searchedUser")).Scan(&followingCount)
	var followerCount int64
	db.QueryRow(context.Background(),"SELECT COUNT(*) FROM relationships JOIN users ON relationships.followed_id = users.id WHERE users.slug = $1", c.Params("searchedUser")).Scan(&followerCount)

	return c.Render("profile", fiber.Map{
		"Username":       searchedUser.Name,
		"Slug":           searchedUser.Slug,
		"CreatedAt":      searchedUser.Created_at,
		"UserTweets": 	  tweets,
		"TotalTweets":    countTweets,
		"WhoToFollow":    whoToFollowUsers,
		"FollowingCount": followingCount,
		"FollowerCount":  followerCount,
		"LoggedInUser":   loggedInUser,
	})
}

func LikeTweet(c *fiber.Ctx) error {
	//fmt.Println("searched user ID: ", c.Params("searchedUser"), "seçilen tweet id: ", c.Params("likedTweetID"))
	db := database.DB
	onlineUser := getUserIdFromCookie(c)
	//	fmt.Println("Aktif kullanıcı id: ", onlineUser.ID)

	tweets := &types.Tweets{}
	//Get who liked the tweet
	db.QueryRow(context.Background(), "SELECT liked_user_id FROM tweets WHERE id = $1", c.Params("likedTweetID")).Scan(&tweets.Liked_user_id)
	for i, _ := range tweets.Liked_user_id {
		if tweets.Liked_user_id[i] == onlineUser.ID {
			tweets.Liked_user_id = append(tweets.Liked_user_id[:i], tweets.Liked_user_id[i+1:]...) //delete user from liked_user's slice
			//fmt.Println(tweets.Liked_user_id)
			_, err := db.Exec(context.Background(),"UPDATE tweets SET liked_user_id = $1 WHERE id = $2", tweets.Liked_user_id, c.Params("likedTweetID"))
			if err != nil{
				log.Fatal(err)
			}
			return c.Redirect("/" + c.Params("searchedUser"))
		}
	}
	tweets.Liked_user_id = append(tweets.Liked_user_id, onlineUser.ID)
	_,err := db.Exec(context.Background(), "UPDATE tweets SET liked_user_id = $1 WHERE id = $2", tweets.Liked_user_id,  c.Params("likedTweetID"))
	if err != nil{
		log.Fatal(err)
	}

	return c.Redirect("/" + c.Params("searchedUser"))
}

func MentionTweet(c *fiber.Ctx) error {

	return c.Redirect("/")
}

func Follow(c *fiber.Ctx) error {
	//log.Println("UUUID= ", utils.ImmutableString(c.Params("uid")))
	var err error
	db := database.DB
	follow := types.Relationships{}

	activeUser := getUserIdFromCookie(c)

	follow.FollowedId, err = strconv.Atoi(c.Params("uid"))
	if err != nil {
		log.Fatal(err)
	}
	follow.FollowerId = int(activeUser.ID)

	var exists bool
	db.QueryRow(context.Background(),"SELECT EXISTS(select follower_id, followed_id FROM relationships where follower_id = $1 AND followed_id = $2)", activeUser.ID, follow.FollowedId).Scan(&exists)
	if exists{
		_, err := db.Exec(context.Background(), "delete from relationships WHERE follower_id = $1 AND followed_id = $2", activeUser.ID, follow.FollowedId)
		if err != nil{
			log.Fatal(err)
		}
	}else{
		_, err := db.Exec(context.Background(), "INSERT INTO relationships(created_at, updated_at, follower_id, followed_id) values($1,$2,$3,$4)", time.Now(), time.Now(), activeUser.ID, follow.FollowedId)
		if err != nil{
			log.Fatal(err)
		}
	}
	return c.Redirect("/"+activeUser.Slug+"")
}

func getUserIdFromCookie(c *fiber.Ctx) types.Users {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(auth.SecretKey), nil
	})
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		c.Redirect("login")
		return types.Users{}
	}
	claims := token.Claims.(*jwt.StandardClaims)
	var user types.Users
	db := database.DB
	db.QueryRow(context.Background(),"SELECT id,created_at,updated_at,name,email,phone,slug FROM users WHERE id = $1", claims.Issuer).Scan(&user.ID,&user.Created_at,&user.Updated_at, &user.Name,&user.Email,&user.Phone,&user.Slug)
	if user.ID == 0 {
		c.Redirect("login")
		return types.Users{}
	}
	return user
}
