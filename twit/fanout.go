package twit

func fanout(tweet Tweet) {
	followerA := dbQryFollows(tweet.UserId)
	for _, follower := range followerA {
		redisInsTweet(follower.Id, tweet)
	}
}

func FanoutLoop() {
	for {
		tweetId := dbGetNextQueuedTweetId()
		// what to do if tweetid is null??
		dbMarkTweetProcessing(tweetId)
		tweet := dbGetTweet(tweetId)
		fanout(tweet)
		dbDequeueTweetId(tweetId)
	}
}
