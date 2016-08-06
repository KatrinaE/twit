package twit

import (
	"github.com/bmizerany/pat" // muxer
	"log"
	"net/http"
	"strconv"
)

func allTweets(w http.ResponseWriter, req *http.Request) {
	tweetA, err := dbQryAllTweets()
	if err != nil {
		writeErrorResponse(w, err)
		return
	}
	writeJsonResponse(w, tweetA)
}

func createTweet(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	userIdS := req.FormValue("UserId")
	userId, err := strconv.Atoi(userIdS)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}
	tweetMsg := req.FormValue("TweetMsg")
	tweet, err := dbInsertTweet(userId, tweetMsg)
	if err != nil {
		writeErrorResponse(w, err)
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
		writeErrorResponse(w, err)
		return
	}

	tweet, err := dbGetTweet(tweetId)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}
	writeJsonResponse(w, tweet)
}

func deleteTweet(w http.ResponseWriter, req *http.Request) {
	tweetIdS := req.URL.Query().Get(":tweetId")
	tweetId, err := strconv.Atoi(tweetIdS)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	dbDelTweet(tweetId)
	response := map[string]string{"status": "ok"}
	writeJsonResponse(w, response)
}

func userTweets(w http.ResponseWriter, req *http.Request) {
	userIdS := req.URL.Query().Get(":userId")
	userId, err := strconv.Atoi(userIdS)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	tweetA, err := dbQryUserTweets(userId)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}
	writeJsonResponse(w, tweetA)
}

func followedTweets(w http.ResponseWriter, req *http.Request) {
	userIdS := req.URL.Query().Get(":userId")
	userId, err := strconv.Atoi(userIdS)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	tweetLites, err := redisGetHomeTimeline(userId)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}
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
