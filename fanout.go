package twit

import (
	"database/sql"
	"fmt"
	"log"
)

func fanout(db *sql.DB, tweet Tweet) {
	followerA, err := dbQryUserFollowers(db, tweet.UserId)
	if err != nil {
		m := "Could not get followers for tweet "
		m += fmt.Sprintf("%+v. Err: %+v", tweet, err)
		log.Print(m)
		return
	}
	for _, follower := range followerA {
		redisClient := newRedisClient()
		err := redisInsertTweet(redisClient, follower.Id, tweet)
		if err != nil {
			m := fmt.Sprintf("Could not insert tweet %d ", tweet.Id)
			m += fmt.Sprintf("for user %d", follower.Id)
			m += fmt.Sprintf("Err: %+v", err)
			log.Print(m)
			return
		}
		m := fmt.Sprintf("Inserted tweet %d ", tweet.Id)
		m += fmt.Sprintf("for user user %d", follower.Id)
		log.Print(m)
	}
}

func FanoutLoop(db *sql.DB) {
	for {
		redisClient := newRedisClient()
		tweetId, err := redisGetNextQueuedTweetId(redisClient)
		if err != nil {
			log.Fatal(err)
		}

		// Note: need to do something if no tweetId is found. Not
		// sure what format the result will be in in this situation.

		tweet, err := dbGetTweet(db, tweetId)
		if err != nil {
			m := fmt.Sprintf("Could not get tweet %d, ", tweetId)
			m += "although queued. "
			m += fmt.Sprintf("Err: %+v", err)
			log.Print(m)
			// Just proceed with next tweet for now. In a
			// real system, would want to log and follow up.
			continue
		}
		fanout(db, tweet)
		redisDequeueTweetId(redisClient, tweetId)
	}
}
