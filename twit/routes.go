package twit

import (
	"github.com/bmizerany/pat" // muxer
	"log"
	"net/http"
	"strconv"
)

func AllTweets(w http.ResponseWriter, req *http.Request) {
	userId := req.URL.Query().Get(":userId")
	tweetA := dbQryUserTweets(userId)
	writeJson(w, tweetA)
}

func CreateTweet(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	userIdS := req.FormValue("UserId")
	userId, err := strconv.Atoi(userIdS)
	if err != nil {
		log.Fatal(err)
	}
	tweetMsg := req.FormValue("TweetMsg")
	tweet := dbInsertTweet(userId, tweetMsg)
	dbEnqueueTweet(tweet)
	writeJson(w, tweet)
}

func GetTweet(w http.ResponseWriter, req *http.Request) {
	tweetId := req.URL.Query().Get(":tweetId")
	tweet := dbGetTweet(tweetId)
	writeJson(w, tweet)
}

func DeleteTweet(w http.ResponseWriter, req *http.Request) {
	tweetId := req.URL.Query().Get(":tweetId")
	s := dbDelTweet(tweetId)
	writeJson(w, tweet)
}

func UserTweets(w http.ResponseWriter, req *http.Request) {
	userId := req.URL.Query().Get(":userId")
	tweetA := dbQryUserTweets(userId)
	writeJson(w, tweetA)
}

func FollowedTweets(w http.ResponseWriter, req *http.Request) {
	userId := req.URL.Query().Get(":userId")
	tweetLites := getHomeTimelineFromRedis(userId)
	// will have to hydrate tweets
}

func RegisterRoutes() *PatternServeMux {
	mux := pat.New()
	mux.Get("/tweets", http.HandlerFunc(AllTweets))
	mux.Post("/tweets", http.HandlerFunc(CreateTweet))
	mux.Get("/tweets/:tweetId", http.HandlerFunc(GetTweet))
	mux.Del("/tweets/:tweetId", http.HandlerFunc(DeleteTweet))
	mux.Get("/tweets/user/:userId", http.HandlerFunc(UserTweets))
	mux.Get("/tweets/followed/:userId", http.HandlerFunc(FollowedTweets))
	return mux
}
