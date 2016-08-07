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

func redisInsertTweet(client *redis.Client, recipientId int, tweet Tweet) error {
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

	recipientIdStr := strconv.Itoa(recipientId)
	err = client.LPush(recipientIdStr, tweetLitePb).Err()
	if err != nil {
		return err
	}
	return err
}

func redisGetHomeTimeline(client *redis.Client, recipientId int) ([]TweetLite, error) {
	recipientIdStr := strconv.Itoa(recipientId)
	result, err := client.LRange(recipientIdStr, 0, -1).Result()
	if err != nil {
		panic(err)
	}

	tweetLites := []TweetLite{}
	for _, s := range result {
		log.Print(s)
		tweetLite := &TweetLite{}
		b := []byte(s)
		err := proto.Unmarshal(b, tweetLite)
		if err != nil {
			log.Printf("Failed to decode tweet: %s Error: %+v",
				s, err)
			return tweetLites, err
		}
		tweetLites = append(tweetLites, *tweetLite)
	}
	return tweetLites, err
}

func redisEnqueueTweetId(client *redis.Client, tweetId int) error {
	tweetIdS := strconv.Itoa(tweetId)
	err := client.LPush("TweetReadyQueue", tweetIdS).Err()
	return err
}

func redisGetNextQueuedTweetId(client *redis.Client) (int, error) {
	tweetIdS, err := client.RPopLPush("TweetReadyQueue", "TweetInProcessQueue").Result()
	tweetId, err := strconv.Atoi(tweetIdS)
	return tweetId, err
}

func redisDequeueTweetId(client *redis.Client, tweetId int) error {
	tweetIdS := strconv.Itoa(tweetId)
	err := client.LRem("TweetInProcessQueue", 1, tweetIdS).Err()
	return err
}
