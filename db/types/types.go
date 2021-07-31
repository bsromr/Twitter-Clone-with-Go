package types

import "time"

type Users struct {
	ID     		int 		`json:"id"`
	Created_at 	time.Time 	`json:"created_at"`
	Updated_at 	time.Time 	`json:"updated_at"`
	Deleted_at 	time.Time 	`json:"deleted_at"`
	Name     	string 		`json:"name"`
	Email    	string 		`json:"email"`
	Phone    	string 		`json:"phone"`
	Password 	string 		`json:"password"`
	Slug     	string 		`json:"slug"`
}

type Relationships struct {
	FollowerId int `json:"follower_id"`
	FollowedId int `json:"followed_id"`
}

type Tweets struct {
	ID int `json:"id"`
	Created_at time.Time `json:"created_at"`
	Deleted_at time.Time `json:"deleted_at"`
	User_id       int    `json:"user_id"`
	Tweet         string `json:"tweet"`
	Liked_user_id []int  `json:"liked_user_id"`
	LikeCount int 		 `json:"like_count"`
	RedLikeButton uint8  `json:"red_like_button"`
}

