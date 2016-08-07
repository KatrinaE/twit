package twit

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
)

func dbInsertTweet(db *sql.DB, userId int, tweetMsg string) (Tweet, error) {
	sql := "INSERT INTO t_tweet (user_id, message) " +
		"VALUES ($1, $2) RETURNING id, user_id, message"
	tweet := Tweet{}
	// Use Query instead of Exec b/c returning rows
	err := db.QueryRow(sql, userId, tweetMsg).
		Scan(&tweet.Id, &tweet.UserId, &tweet.Message)
	if err != nil {
		return tweet, err
	}
	return tweet, err
}

func dbGetTweet(db *sql.DB, tweetId int) (Tweet, error) {
	tweet := Tweet{}
	query := "SELECT id, user_id, message FROM t_tweet WHERE id = $1"
	row := db.QueryRow(query, tweetId)
	err := row.Scan(&tweet.Id, &tweet.UserId, &tweet.Message)
	if err != nil {
		return tweet, err
	}
	return tweet, err
}

func dbDelTweet(db *sql.DB, tweetId int) error {
	_, err := db.Exec("DELETE FROM t_tweet WHERE id = $1", tweetId)
	if err != nil {
		return err
	}
	return err
}

func dbQryTweets(db *sql.DB, whereClause string) ([]Tweet, error) {
	tweetA := []Tweet{}
	sA := []string{"SELECT id, user_id, message FROM t_tweet", whereClause}
	query := strings.Join(sA, " ")
	rows, err := db.Query(query)
	switch {
	case err == sql.ErrNoRows:
		// No problem
		return tweetA, nil
	case err != nil:
		return tweetA, err
	}
	defer rows.Close()
	for rows.Next() {
		var tweet Tweet
		err := rows.Scan(&tweet.Id, &tweet.UserId, &tweet.Message)
		if err != nil {
			return tweetA, err
		}
		tweetA = append(tweetA, tweet)
	}
	if err := rows.Err(); err != nil {
		return tweetA, err
	}
	return tweetA, err
}

func dbQryAllTweets(db *sql.DB) ([]Tweet, error) {
	whereClause := ""
	tweetA, err := dbQryTweets(db, whereClause)
	return tweetA, err
}

func dbQryUserTweets(db *sql.DB, userId int) ([]Tweet, error) {
	whereClause := fmt.Sprintf("WHERE user_id = %d", userId)
	tweetA, err := dbQryTweets(db, whereClause)
	return tweetA, err
}

func dbQryUserFollowers(db *sql.DB, userId int) ([]Follow, error) {
	followA := []Follow{}
	query := "SELECT id, follower_id, followed_id FROM t_follower WHERE followed_id=$1"
	rows, err := db.Query(query, userId)
	defer rows.Close()
	if err != nil {
		return followA, err
	}

	for rows.Next() {
		var follow Follow
		err := rows.Scan(&follow.Id, &follow.FollowerId, &follow.FollowedId)
		if err != nil {
			return followA, err
		}
		followA = append(followA, follow)
	}
	return followA, err
}
