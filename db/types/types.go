package types

import "github.com/jinzhu/gorm"

type Users struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Slug     string `json:"slug"`
}

type Relationships struct {
	gorm.Model
	FollowerId int `json:"follower_id"`
	FollowedId int `json:"followed_id"`
}

type Tweets struct {
	gorm.Model
	User_id       int    `json:"user_id"`
	Tweet         string `json:"tweet"`
	LikeCount     int    `json:"likeCount`
	Liked_user_id int    `sql:"-"`
}

type Tweet_Info struct {
	gorm.Model
	Tweet_id      int `json:"tweet_id`
	Liked_user_id int `json:"liked_user_id`
}
