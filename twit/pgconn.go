package twit

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

func dbInsertTweet(userId int, tweetMsg string) Tweet {
	dbDriver, dbOpen := GetDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		log.Fatal(err)
	}
	sql := "INSERT INTO t_tweet (user_id, message) " +
		"VALUES ($1, $2) RETURNING id, user_id, message"
	row, err := db.Query(sql, userId, tweetMsg) // use query b/c no lastInsertId
	if err != nil {
		log.Fatal(err)
	}
	tweet := Tweet{}
	err = row.Scan(&tweet.Id, &tweet.UserId, &tweet.Message)
	if err != nil {
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
	query := "SELECT id, user_id, message FROM t_tweet WHERE id = $1"
	row := db.QueryRow(query, tweetId)
	tweet := Tweet{}
	err = row.Scan(&tweet.Id, &tweet.UserId, &tweet.Message)
	if err != nil {
		log.Fatal(err)
	}
	return tweet
}

func dbQryTweets(whereClause string) []Tweet {
	dbDriver, dbOpen := GetDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		log.Fatal(err)
	}
	sA := []string{"SELECT id, user_id, text FROM tweets", whereClause}
	query := strings.Join(sA, " ")
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	tweetA := []Tweet{}
	for rows.Next() {
		var tweet Tweet
		err := rows.Scan(&tweet.Id, &tweet.UserId, &tweet.Message)
		if err != nil {
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
	whereClause := ""
	tweetA := dbQryTweets(whereClause)
	return tweetA
}

func dbQryUserTweets(userId int) []Tweet {
	whereClause := fmt.Sprintf("WHERE user_id = %d", userId)
	tweetA := dbQryTweets(whereClause)
	return tweetA
}

func dbDelTweet(tweetId int) {
	dbDriver, dbOpen := GetDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("DELETE FROM t_tweet WHERE id = $1", tweetId)
}

func dbEnqueueTweetId(tweetId int) {
	dbDriver, dbOpen := GetDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	sql := "INSERT INTO t_tweet_queue (tweet_id, status) " +
		"VALUES ($1, 'ready')"
	_, err = db.Exec(sql, tweetId)
	if err != nil {
		log.Fatal(err)
	}
}

func dbGetNextQueuedTweetId() int {
	dbDriver, dbOpen := GetDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		log.Fatal(err)
	}
	query := "SELECT tweet_id FROM t_tweet_queue WHERE status='ready' " +
		"ORDER BY ctime ASC LIMIT 1"
	row := db.QueryRow(query)
	if err != nil {
		log.Fatal(err)
	}
	var tweetId int
	err = row.Scan(&tweetId)
	if err != nil {
		log.Fatal(err)
	}
	return tweetId
}

func dbMarkTweetProcessing(tweetId int) {
	dbDriver, dbOpen := GetDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	sql := "UPDATE t_tweet_queue SET status='processing' " +
		"WHERE tweet_id=$1"
	_, err = db.Exec(sql, tweetId)
	if err != nil {
		log.Fatal(err)
	}
}

func dbDequeueTweetId(tweetId int) {
	dbDriver, dbOpen := GetDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		log.Fatal(err)
	}
	sql := "DELETE from t_tweet_queue WHERE tweet_id=$1"
	_, err = db.Exec(sql, tweetId)
	if err != nil {
		log.Fatal(err)
	}
}

func dbQryFollows(userId int) []Follow {
	dbDriver, dbOpen := GetDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		log.Fatal(err)
	}
	query := "SELECT id FROM t_follower WHERE followed_id=$1"
	rows, err := db.Query(query, userId)
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
		followA = append(followA, follow)
	}
	return followA
}
