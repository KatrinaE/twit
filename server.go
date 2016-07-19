// http://thenewstack.io/make-a-restful-json-api-go/
// https://github.com/golang/go/wiki/SQLInterface
package main

import (
//	"fmt"
//    "html"
	"io"
	"log"
	"net/http"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/bmizerany/pat"
	"local/twit/twitutil"
)

const (
    DB_HOST = ""
    DB_NAME = "foo"
    DB_USER = "foo"
    DB_PASS = "foo"
)

func AllTweets(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "AllTweets!")
	db, err := sql.Open("postgres", "user=DB_USER dbname=DB_NAME sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT id, user_id, text FROM users")
	if err != nil {
		log.Fatal(err)
	}
/*	for rows.Next() {
		var name string
		if err := rows.Scan(&id); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s is foo\n", id)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
*/
}

func CreateTweet(w http.ResponseWriter, req *http.Request) {
//	tweetId := twit.makeUUID()
	tweetId := 1
	userId := req.URL.Query().Get(":userId")
	tweetText := req.URL.Query().Get(":tweetText")
	//io.WriteString(w, "CreateTweet " + tweetId)
	db, err := sql.Open("postgres", "user=DB_USER dbname=DB_NAME sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	result, err := db.Exec(
		"INSERT INTO tweets (id, user_id, text) VALUES ($1, $2, $3)",
		tweetId,
		userId,
		tweetText,
	)
}

func GetTweet(w http.ResponseWriter, req *http.Request) {
	tweetId := req.URL.Query().Get(":tweetId")
	io.WriteString(w, "GetTweet " + tweetId)
	db, err := sql.Open("postgres", "user=DB_USER dbname=DB_NAME sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	row := db.QueryRow("SELECT id, user_id, text FROM users WHERE id = $1", tweetId)
//	err := row.Scan(&id)
}

func DeleteTweet(w http.ResponseWriter, req *http.Request) {
	tweetId := req.URL.Query().Get(":tweetId")
	db, err := sql.Open("postgres", "user=DB_USER dbname=DB_NAME sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	io.WriteString(w, "DeleteTweet " + tweetId)
	result, err := db.Exec("DELETE FROM tweets WHERE id = $1", tweetId)
}

func UserTweets(w http.ResponseWriter, req *http.Request) {
	userId := req.URL.Query().Get(":userId")
	io.WriteString(w, "UserTweets " + userId)
	db, err := sql.Open("postgres", "user=DB_USER dbname=DB_NAME sslmode=disable")
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
}

func FollowedTweets(w http.ResponseWriter, req *http.Request) {
	userId := req.URL.Query().Get(":userId")
	io.WriteString(w, "UserFollowingTweets " + userId)
	// will have to use redis for this
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
