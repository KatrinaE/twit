package twit

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

type Tweet struct {
	Id      int    `json:"id"`
	UserId  int    `json:"userId"`
	Message string `json:"message"`
}

type DisplayTweet struct {
	User
	Tweet
}

type Follow struct {
	Id         int `json:"id"`
	FollowerId int `json:"followerId"`
	FollowedId int `json:"followedId"`
}
