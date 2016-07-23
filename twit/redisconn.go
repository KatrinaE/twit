package twit

import (
	"github.com/golang/protobuf/proto"
	"gopkg.in/redis.v4"
	pb "local/twit/protobuf"
	"strconv"
)

func newRedisClient() redisclient {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	// need to check client err???

	return client
}

func redisInsTweet(recipientId int, tweet Tweet) {
	// is there a better way to do this? i.e. like a tweet without a message?
	tweetLite := &pb.TweetLite{
		Id:     proto.Int(tweet.Id),
		UserId: proto.Int(tweet.UserId),
	}
	pbTweetLite, err := proto.Marshal(tweetLite)
	if err != nil {
		log.Fatalln("Failed to encode tweet:", err)

	}
	fmt.Println("sending to recipient ", recipientId)
	fmt.Println(tweetLite)

	client := newRedisClient()
	recipientIdStr := strconv.Itoa(recipientId)
	err1 := client.LPush(recipientIdStr, pbTweetLite).Err()
	if err1 != nil {
		panic(err1)
	}
}

func getHomeTimelineFromRedis(recipientId int) []pb.TweetLite {
	recipientIdStr := strconv.Itoa(recipientId)
	client := newRedisClient()
	result, err := client.LRange(recipientIdStr, 0, -1).Result()
	if err != nil {
		panic(err)
	}

	tweetLites := []pb.TweetLite{}
	for _, s := range result {
		tweetLite := &pb.TweetLite{}
		b := []byte(s)
		err := proto.Unmarshal(b, tweetLite)
		if err != nil {
			log.Fatalln("Failed to decode tweet:", err)
		}
		tweetLites = append(tweetLites, *tweetLite)
	}
	return tweetLites
}
