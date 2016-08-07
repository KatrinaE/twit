package twit

import (
	"bytes"
	"database/sql"
	"github.com/bmizerany/pat" // muxer
	"net/http"
	"strconv"
	"strings"
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

func userTimelineTweets(w http.ResponseWriter, req *http.Request) {
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

func homeTimelineTweets(w http.ResponseWriter, req *http.Request) {
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
	dbDriver, dbOpen := getDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)

	// Create WHERE clause to get full tweets for all tweet ids in timeline
	var buffer bytes.Buffer
	buffer.WriteString("WHERE t_tweet.id IN (")
	for _, tweetLite := range tweetLites {
		buffer.WriteString(strconv.Itoa(int(*tweetLite.Id)))
		buffer.WriteString(",")
	}
	whereClause := buffer.String()
	whereClause = strings.TrimRight(whereClause, ",")
	whereClause += ")"

	displayTweetA, err := dbQryDisplayTweets(db, whereClause)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}
	writeJsonResponse(w, displayTweetA)
}

func RegisterRoutes() *pat.PatternServeMux {
	mux := pat.New()

	mux.Get("/tweets", http.HandlerFunc(allTweets))
	mux.Post("/tweets", http.HandlerFunc(createTweet))
	mux.Get("/tweets/:tweetId", http.HandlerFunc(getTweet))
	mux.Del("/tweets/:tweetId", http.HandlerFunc(deleteTweet))
	mux.Get("/tweets/user_timeline/:userId", http.HandlerFunc(userTimelineTweets))
	mux.Get("/tweets/home_timeline/:userId", http.HandlerFunc(homeTimelineTweets))
	return mux
}
