// http://thenewstack.io/make-a-restful-json-api-go/
// https://github.com/golang/go/wiki/SQLInterface
// hi
package main

import (
	"fmt"
	"strconv"
	//    "html"
	"database/sql"
	"encoding/json"
	"github.com/bmizerany/pat"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net/http"
)

func AllTweets(w http.ResponseWriter, req *http.Request) {
	dbDriver, dbOpen := GetDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT id, user_id, message FROM t_tweet")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	tweetsSlice := []Tweet{}
	for rows.Next() {
		tweet := Tweet{}
		err := rows.Scan(&tweet.Id, &tweet.UserId, &tweet.Message)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(tweet)
		tweetsSlice = append(tweetsSlice, tweet)
	}
	b, err := json.Marshal(tweetsSlice)
	w.Write(b)
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func CreateTweet(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	userIdS := req.FormValue("UserId")
	userId, err := strconv.Atoi(userIdS)
	if err != nil {
		log.Fatal(err)
	}
	tweetMsg := req.FormValue("TweetMsg")

	dbDriver, dbOpen := GetDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		log.Fatal(err)
	}

	row := db.QueryRow(
		"INSERT INTO t_tweet (user_id, message) VALUES ($1, $2) RETURNING id",
		userId,
		tweetMsg)
	if err != nil {
		log.Fatal(err)
	}

	var id int
	err2 := row.Scan(&id)
	if err2 != nil {
		log.Fatal(err)
	}

	_, err3 := db.Exec(
		"INSERT INTO t_tweet_queue (tweet_id, status) VALUES ($1, 'ready')",
		id)
	if err3 != nil {
		log.Fatal(err3)
	}

	tweet := Tweet{id, userId, tweetMsg}
	b, err := json.Marshal(tweet)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(b)
}

func GetTweet(w http.ResponseWriter, req *http.Request) {
	tweetId := req.URL.Query().Get(":tweetId")
	dbDriver, dbOpen := GetDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		log.Fatal(err)
	}
	row := db.QueryRow("SELECT id, user_id, message FROM t_tweet WHERE id = $1", tweetId)
	tweet := Tweet{}
	err2 := row.Scan(&tweet.Id, &tweet.UserId, &tweet.Message)
	if err2 != nil {
		log.Fatal(err2)
	}
	b, err := json.Marshal(tweet)
	w.Write(b)
}

func DeleteTweet(w http.ResponseWriter, req *http.Request) {
	tweetId := req.URL.Query().Get(":tweetId")
	dbDriver, dbOpen := GetDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		log.Fatal(err)
	}
	result, err := db.Exec("DELETE FROM t_tweet WHERE id = $1", tweetId)
	fmt.Println(result)
	s := "Done"
	b, err := json.Marshal(s)
	w.Write(b)
}

func UserTweets(w http.ResponseWriter, req *http.Request) {
	userId := req.URL.Query().Get(":userId")
	io.WriteString(w, "UserTweets "+userId)
	dbDriver, dbOpen := GetDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT id, user_id, text FROM tweets WHERE user_id = $1", userId)
	if err != nil {
		log.Fatal(err)
	}
	/*	for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s is bar\n", id)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
	*/
	fmt.Println(rows)
}

func FollowedTweets(w http.ResponseWriter, req *http.Request) {
	userId := req.URL.Query().Get(":userId")
	io.WriteString(w, "UserFollowingTweets "+userId)
	dbDriver, dbOpen := GetDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	// will have to use redis for this
	fmt.Println(db)
	fmt.Println(err)
}

func main() {
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
