package twit

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func dbInsertTweet(userId int, tweetMsg string) {
	dbDriver, dbOpen := GetDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		log.Fatal(err)
	}
	row, err := db.QueryRow(
		"INSERT INTO t_tweet (user_id, message) VALUES ($1, $2) RETURNING id, user_id, message",
		userId,
		tweetMsg)
	if err != nil {
		log.Fatal(err)
	}
	tweet := Tweet{}
	err2 := row.Scan(&tweet.Id, &tweet.UserId, &tweet.Msg)
	if err2 != nil {
		log.Fatal(err)
	}
	return tweet

}

func dbGetTweet(tweetId int) Tweet {
	dbDriver, dbOpen := GetDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		log.Fatal(err)
	}
	row := db.QueryRow("SELECT id, user_id, message FROM t_tweet WHERE id = $1", tweetId)
	tweet := Tweet{}
	err := row.Scan(&tweet.Id, &tweet.UserId, &tweet.Message)
	if err != nil {
		log.Fatal(err)
	}
	return tweet
}

func dbQryTweets(query string) []Tweet {
	dbDriver, dbOpen := GetDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query(q, userId)
	if err != nil {
		log.Fatal(err)
	}
	tweetA := []Tweet{}
	for rows.Next() {
		var tweet Tweet
		// warning: this assumes they did a select *
		if err := rows.Scan(&tweet.Id, &tweet.UserId, &tweet.Msg); err != nil {
			log.Fatal(err)
		}
		tweetA = append(tweetA, tweet)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return tweetA
}

func dbQryAllTweets() []Tweet {
	io.WriteString(w, "AllTweets ")
	// q := "SELECT id, user_id, text FROM tweets WHERE user_id = $1"
	// how to do it with native sql like above???
	q := fmt.Sprintf("SELECT id, user_id, text FROM tweets")
	tweetA := dbQryTweets(q)
	return tweetA
}

func dbQryUserTweets(userId int) []Tweet {
	io.WriteString(w, "UserTweets "+userId)
	// q := "SELECT id, user_id, text FROM tweets WHERE user_id = $1"
	// how to do it with native sql like above???
	q := fmt.Sprintf("SELECT id, user_id, text FROM tweets WHERE user_id = %d")
	tweetA := dbQryTweets(q)
	return tweetA
}

func dbDelTweet(tweetId int) {
	dbDriver, dbOpen := GetDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		log.Fatal(err)
	}
	result, err := db.Exec("DELETE FROM t_tweet WHERE id = $1", tweetId)
	return "Done"
}

func dbEnqueueTweetId(tweetId int) {
	dbDriver, dbOpen := GetDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	_, err := db.Exec(
		"INSERT INTO t_tweet_queue (tweet_id, status) VALUES ($1, 'ready')",
		tweetId)
	if err != nil {
		log.Fatal(err)
	}
}

func dbDequeueNextTweetId() int {
	dbDriver, dbOpen := GetDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		log.Fatal(err)
	}
	row, err := db.QueryRow("SELECT tweet_id FROM t_tweet_queue WHERE status='ready' ORDER BY ctime ASC LIMIT 1")
	if err != nil {
		log.Fatal(err)
	}

	// need to update. can you update with limit??

	var tweetId int
	err := row.Scan(&tweetId)
	if err != nil {
		log.Fatal(err)
	}
	return tweetId
}

func dbQryFollows(userId int) []Follow {
	dbDriver, dbOpen := GetDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		log.Fatal(err)
	}
	query := fmt.Sprintf("SELECT id FROM t_follower WHERE followed_id=%d", userId)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	followA := []Follow{}
	for rows.Next() {
		var follow Follow
		err := rows.Scan(&follow.Id)
		if err != nil {
			log.Fatal(err)
		}
		followeA = append(followA, follow)
	}

	return followA
}
