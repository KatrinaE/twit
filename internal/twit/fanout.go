package twit

import (
	"database/sql"
	"log"
	"time"
)

func fanout(tweet Tweet) {
	followerA, err := dbQryUserFollowers(tweet.UserId)
	if err != nil {
		log.Print("Could not insert tweet %+v. Err: %+v", tweet, err)
		return
	}
	for _, follower := range followerA {
		redisInsertTweet(follower.Id, tweet)
	}
}

func FanoutLoop() {
	for {
		tweetId, err := dbGetNextQueuedTweetId()
		switch {
		case err == sql.ErrNoRows:
			time.Sleep(300 * time.Millisecond)
		case err != nil:
			log.Print("%+v", err)
			// do not terminate execution -- maybe
			// revisit this decision later
		}
		dbMarkTweetProcessing(tweetId)
		tweet, err := dbGetTweet(tweetId)
		if err != nil {
			// todo: what about situations where you can't
			// connect to the database at all?
			dbMarkTweetErrored(tweetId)
			continue
		}
		fanout(tweet)
		dbDequeueTweetId(tweetId)
	}
}
