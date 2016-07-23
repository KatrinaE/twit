package twit

type User struct {
	Id       int
	Username string
}

type Tweet struct {
	Id      int
	UserId  int
	Message string
}

type Follow struct {
	Id         int
	FollowerId int
	FollowedId int
}
