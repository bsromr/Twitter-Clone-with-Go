package types

import "time"

type Users struct {
	ID     		int `json:"id"`
	Name     	string `json:"name"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Deleted_at time.Time `json:"deleted_at"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Slug     string `json:"slug"`
}

type Relationships struct {
	FollowerId int `json:"follower_id"`
	FollowedId int `json:"followed_id"`
}

type Tweets struct {
	User_id       int    `json:"user_id"`
	Tweet         string `json:"tweet"`
	LikeCount     int    `json:"likeCount"`
	Liked_user_id int    `sql:"-"`
}

type Tweet_Info struct {
	Tweet_id      int `json:"tweet_id`
	Liked_user_id int `json:"liked_user_id`
}
