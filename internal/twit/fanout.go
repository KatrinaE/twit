package twit

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func fanout(tweet Tweet) {
	followerA, err := dbQryUserFollowers(tweet.UserId)
	if err != nil {
		m := "Could not get followers for tweet "
		m += fmt.Sprintf("%+v. Err: %+v", tweet, err)
		log.Print(m)
		return
	}
	for _, follower := range followerA {
		err := redisInsertTweet(follower.Id, tweet)
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

func FanoutLoop() {
	for {
		tweetId := redisGetNextQueuedTweetId()
		tweet, err := dbGetTweet(tweetId)
		if err != nil {
			m := fmt.Sprintf("Could not get tweet %d, ", tweetId)
			m += "although queued. "
			m += fmt.Sprintf("Err: %+v", err)
			log.Print(m)
			// Just proceed with next tweet for now. In a
			// real system, would want to log and follow up.
			continue
		}
		fanout(tweet)
		redisDequeueTweetId(tweetId)
	}
}
