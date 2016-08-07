package twit

import (
	"database/sql"
	"github.com/bmizerany/pat" // muxer
	"net/http"
	"strconv"
)

func allTweets(w http.ResponseWriter, req *http.Request) {
	dbDriver, dbOpen := getDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}
	tweetA, err := dbQryAllTweets(db)
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
	dbDriver, dbOpen := getDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}
	tweetMsg := req.FormValue("TweetMsg")
	tweet, err := dbInsertTweet(db, userId, tweetMsg)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}
	redisClient := newRedisClient()
	redisEnqueueTweetId(redisClient, tweet.Id)
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
	dbDriver, dbOpen := getDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}
	tweet, err := dbGetTweet(db, tweetId)
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

	dbDriver, dbOpen := getDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}
	dbDelTweet(db, tweetId)
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

	dbDriver, dbOpen := getDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}
	tweetA, err := dbQryUserTweets(db, userId)
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
	redisClient := newRedisClient()
	tweetLites, err := redisGetHomeTimeline(redisClient, userId)
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
