// http://thenewstack.io/make-a-restful-json-api-go/
// https://github.com/golang/go/wiki/SQLInterface
package main

import (
	"fmt"
//    "html"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/bmizerany/pat"
//	"local/twit/twitutil"
	"github.com/spf13/viper"
)

type Tweet struct {
    Id int
    UserId int
    Message string
}

func getDbConfig() (string, string) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("dbconf")
	viper.AddConfigPath("./db/")   // right now dbconf is only config
	err := viper.ReadInConfig()
	if err != nil {
	    panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	env := viper.GetString("environment")
	dbDriverField := fmt.Sprintf("%s.driver", env)
	dbDriver := viper.GetString(dbDriverField)
	openField := fmt.Sprintf("%s.open", env)
	dbOpen := viper.GetString(openField)
	return dbDriver, dbOpen
}

func AllTweets(w http.ResponseWriter, req *http.Request) {
	dbDriver, dbOpen := getDbConfig()
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
//	tweetId := twit.makeUUID()
	tweetId := 1
	userId := req.URL.Query().Get(":userId")
	tweetText := req.URL.Query().Get(":tweetText")
	//io.WriteString(w, "CreateTweet " + tweetId)
	dbDriver, dbOpen := getDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		log.Fatal(err)
	}
	result, err := db.Exec(
		"INSERT INTO tweets (id, user_id, text) VALUES ($1, $2, $3)",
		tweetId,
		userId,
		tweetText,
	)
	fmt.Println(result)
}

func GetTweet(w http.ResponseWriter, req *http.Request) {
	tweetId := req.URL.Query().Get(":tweetId")
	io.WriteString(w, "GetTweet " + tweetId)
	dbDriver, dbOpen := getDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		log.Fatal(err)
	}
	row := db.QueryRow("SELECT id, user_id, text FROM users WHERE id = $1", tweetId)
//	err := row.Scan(&id)
	fmt.Println(row)
}

func DeleteTweet(w http.ResponseWriter, req *http.Request) {
	tweetId := req.URL.Query().Get(":tweetId")
	dbDriver, dbOpen := getDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		log.Fatal(err)
	}
	io.WriteString(w, "DeleteTweet " + tweetId)
	result, err := db.Exec("DELETE FROM tweets WHERE id = $1", tweetId)
	fmt.Println(result)
}

func UserTweets(w http.ResponseWriter, req *http.Request) {
	userId := req.URL.Query().Get(":userId")
	io.WriteString(w, "UserTweets " + userId)
	dbDriver, dbOpen := getDbConfig()
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
	io.WriteString(w, "UserFollowingTweets " + userId)
	dbDriver, dbOpen := getDbConfig()
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
