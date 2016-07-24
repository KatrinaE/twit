package twit

import (
	"github.com/bmizerany/pat" // muxer
	"log"
	"net/http"
	"strconv"
)

func allTweets(w http.ResponseWriter, req *http.Request) {
	userIdS := req.URL.Query().Get(":userId")
	userId, err := strconv.Atoi(userIdS)
	if err != nil {
		log.Fatal(err)
	}
	tweetA, err := dbQryUserTweets(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	writeJsonResponse(w, tweetA)
}

func createTweet(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	userIdS := req.FormValue("UserId")
	userId, err := strconv.Atoi(userIdS)
	if err != nil {
		log.Fatal(err)
	}
	tweetMsg := req.FormValue("TweetMsg")
	tweet, err := dbInsertTweet(userId, tweetMsg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	dbEnqueueTweetId(tweet.Id)
	w.WriteHeader(http.StatusCreated)
	writeJsonResponse(w, tweet)
}

func getTweet(w http.ResponseWriter, req *http.Request) {
	tweetIdS := req.URL.Query().Get(":tweetId")
	tweetId, err := strconv.Atoi(tweetIdS)
	if err != nil {
		log.Fatal(err)
	}

	tweet, err := dbGetTweet(tweetId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	writeJsonResponse(w, tweet)
}

func deleteTweet(w http.ResponseWriter, req *http.Request) {
	tweetIdS := req.URL.Query().Get(":tweetId")
	tweetId, err := strconv.Atoi(tweetIdS)
	if err != nil {
		log.Fatal(err)
	}

	dbDelTweet(tweetId)
	response := map[string]string{"status": "ok"}
	writeJsonResponse(w, response)
}

func userTweets(w http.ResponseWriter, req *http.Request) {
	userIdS := req.URL.Query().Get(":userId")
	userId, err := strconv.Atoi(userIdS)
	if err != nil {
		log.Fatal(err)
	}

	tweetA, err := dbQryUserTweets(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	writeJsonResponse(w, tweetA)
}

func followedTweets(w http.ResponseWriter, req *http.Request) {
	userIdS := req.URL.Query().Get(":userId")
	userId, err := strconv.Atoi(userIdS)
	if err != nil {
		log.Fatal(err)
	}

	tweetLites := redisGetHomeTimeline(userId)
	writeJsonResponse(w, tweetLites)
	// will have to hydrate tweets
}

func RegisterRoutes() *pat.PatternServeMux {
	mux := pat.New()
	mux.Get("/tweets", http.HandlerFunc(allTweets))
	mux.Post("/tweets", http.HandlerFunc(createTweet))
	mux.Get("/tweets/:tweetId", http.HandlerFunc(getTweet))
	mux.Del("/tweets/:tweetId", http.HandlerFunc(deleteTweet))
	mux.Get("/tweets/user/:userId", http.HandlerFunc(userTweets))
	mux.Get("/tweets/followed/:userId", http.HandlerFunc(followedTweets))
	return mux
}
