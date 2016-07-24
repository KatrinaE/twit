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
			m := "Could not insert tweet "
			m += fmt.Sprintf("%+v into %d's ", tweet, follower.Id)
			m += fmt.Sprintf("home timeline. Err: %+v", err)
			log.Print(m)
			return
		}
		m := fmt.Sprintf("Inserted tweet %+v into ", tweet)
		m += fmt.Sprintf("user %d's home timeline", follower.Id)
		log.Print(m)
	}
}

func FanoutLoop() {
	for {
		tweetId, err := dbGetNextQueuedTweetId()
		switch {
		case err == sql.ErrNoRows:
			time.Sleep(3000 * time.Millisecond)
			continue
		case err != nil:
			log.Printf("%+v", err)
			// do not terminate execution -- maybe
			// revisit this decision later
			continue
		}
		tweet, err := dbGetTweet(tweetId)
		if err != nil {
			m := fmt.Sprintf("Could not get tweet %d, ", tweetId)
			m += "even though it was queued. "
			m += fmt.Sprintf("Err: %+v", err)
			log.Print(m)
			dbMarkTweetErrored(tweetId)
			if err != nil {
				m := "Could not mark tweet "
				m += fmt.Sprintf("%d", tweetId)
				m += " errored. Error: "
				m += fmt.Sprintf("%+v", err)
				log.Print(m)
			}
			continue
		}
		err = dbMarkTweetProcessing(tweetId)
		if err != nil {
			m := fmt.Sprintf("Could not mark tweet %d ", tweetId)
			m += "processing. Aborting fanout."
			log.Print(m)
			continue
		}
		fanout(tweet)
		dbDequeueTweetId(tweetId)
	}
}
