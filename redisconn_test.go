package twit

import (
	"github.com/golang/protobuf/proto"
	"strconv"
	"testing"
)

func testRedisInsertTweet(t *testing.T) {
	recipientId := 500
	recipientIdStr := strconv.Itoa(recipientId)
	userId := 1
	tweetId := 1
	tweetFixture := Tweet{
		Id:      tweetId,
		UserId:  userId,
		Message: "i have lots of followers",
	}
	tweetLiteFixture := TweetLite{
		Id:     proto.Int(tweetId),
		UserId: proto.Int(userId),
	}
	redisClient := newRedisClient()
	err := redisInsertTweet(redisClient, recipientId, tweetFixture)
	if err != nil {
		t.Logf("%+v", err)
		redisClient.Del(recipientIdStr)
		t.FailNow()
	}
	result, err := redisClient.LPop(recipientIdStr).Result()
	if err != nil {
		t.Logf("%+v", err)
		redisClient.Del(recipientIdStr)
		t.FailNow()
	}
	b := []byte(result)
	tweetLite := TweetLite{}
	proto.Unmarshal(b, &tweetLite)
	if tweetLite.UserId != tweetLiteFixture.UserId {
		t.Logf("Wanted: %+v -- Got: %+v", tweetLiteFixture, tweetLite)
		redisClient.Del(recipientIdStr)
		t.FailNow()
	}
	if tweetLite.Id != tweetLiteFixture.Id {
		t.Logf("Wanted: %+v -- Got: %+v", tweetLiteFixture, tweetLite)
		redisClient.Del(recipientIdStr)
		t.FailNow()
	}

	redisClient.Del(recipientIdStr)
}

func testRedisGetHomeTimeline(t *testing.T) {
	recipientId := 500
	recipientIdStr := strconv.Itoa(recipientId)
	tweetLiteFixtureA := []TweetLite{
		TweetLite{Id: proto.Int(1), UserId: proto.Int(1)},
		TweetLite{Id: proto.Int(2), UserId: proto.Int(1)},
	}
	redisClient := newRedisClient()
	for _, tweetLite := range tweetLiteFixtureA {
		tweetLitePb, err := proto.Marshal(&tweetLite)
		if err != nil {
			t.Logf("%+v", err)
			redisClient.Del(recipientIdStr)
			t.FailNow()
		}
		err = redisClient.LPush(recipientIdStr, tweetLitePb).Err()
		if err != nil {
			t.Logf("%+v", err)
			redisClient.Del(recipientIdStr)
			t.FailNow()
		}
	}

	tweetLiteA, err := redisGetHomeTimeline(redisClient, recipientId)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if len(tweetLiteA) != len(tweetLiteFixtureA) {
		t.Logf("tweetLiteA and tweetLiteFixtureA are not the same length")
		t.Logf("tweetLiteA: %+v", tweetLiteA)
		t.Logf("tweetLiteFixtureA: %+v", tweetLiteFixtureA)
		redisClient.Del(recipientIdStr)
		t.FailNow()
	}

	for i, _ := range tweetLiteA {
		if tweetLiteA[i].UserId != tweetLiteFixtureA[i].UserId {
			t.Logf("Wanted: %+v -- Got: %+v",
				tweetLiteFixtureA, tweetLiteA)
			redisClient.Del(recipientIdStr)
			t.FailNow()
		}
		if tweetLiteA[i].Id != tweetLiteFixtureA[i].Id {
			t.Logf("Wanted: %+v -- Got: %+v",
				tweetLiteFixtureA, tweetLiteA)
			redisClient.Del(recipientIdStr)
			t.FailNow()
		}
	}
	redisClient.Del(recipientIdStr)
}
