package twit

func fanout(tweet Tweet) {
	followerA := dbQryFollows(tweet)
	for _, follower := range followerA {
		redisInsTweet(follower.Id, tweet)
	}
}

func FanoutLoop() {
	for {
		tweetId := dbDequeueNextTweetId()
		// what to do if tweetid is null??
		tweet := dbGetTweet(tweetId)
		fanout(tweet)
	}
}
