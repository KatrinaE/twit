package twit

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"gopkg.in/redis.v4"
	"log"
	"strconv"
)

func newRedisClient() *redis.Client {
	address, password, db := getRedisConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	return client
}

func redisInsertTweet(recipientId int, tweet Tweet) error {
	// Use lightweight version of Tweet (no message) in Redis
	tweetLite := &TweetLite{
		Id:     proto.Int(tweet.Id),
		UserId: proto.Int(tweet.UserId),
	}
	tweetLitePb, err := proto.Marshal(tweetLite)
	if err != nil {
		return err
	}
	fmt.Println("sending to recipient ", recipientId)
	fmt.Println(tweetLite)

	client := newRedisClient()
	recipientIdStr := strconv.Itoa(recipientId)
	err = client.LPush(recipientIdStr, tweetLitePb).Err()
	if err != nil {
		return err
	}
	return err
}

func redisGetHomeTimeline(recipientId int) []TweetLite {
	recipientIdStr := strconv.Itoa(recipientId)
	client := newRedisClient()
	result, err := client.LRange(recipientIdStr, 0, -1).Result()
	if err != nil {
		panic(err)
	}

	tweetLites := []TweetLite{}
	for _, s := range result {
		tweetLite := &TweetLite{}
		b := []byte(s)
		err := proto.Unmarshal(b, tweetLite)
		if err != nil {
			log.Fatalln("Failed to decode tweet:", err)
		}
		tweetLites = append(tweetLites, *tweetLite)
	}
	return tweetLites
}
