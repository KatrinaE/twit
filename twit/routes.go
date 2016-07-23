package twit

import (
	"fmt"
	"github.com/bmizerany/pat" // muxer
	"io"
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

func HandleRequest() {
	m := pat.New()
	m.Get("/tweets", http.HandlerFunc(AllTweets))
	m.Post("/tweets", http.HandlerFunc(CreateTweet))
	m.Get("/tweets/:tweetId", http.HandlerFunc(GetTweet))
	m.Del("/tweets/:tweetId", http.HandlerFunc(DeleteTweet))
	m.Get("/tweets/user/:userId", http.HandlerFunc(UserTweets))
	m.Get("/tweets/followed/:userId", http.HandlerFunc(FollowedTweets))

	// Register this pat with the default serve mux so that other packages
	// may also be exported. (i.e. /debug/pprof/*)
	http.Handle("/", m)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
