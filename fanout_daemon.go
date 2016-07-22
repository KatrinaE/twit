// Get next tweet from queue using RPOPLPUSH (http://redis.io/commands/rpoplpush)
// Set its status to 'processing'
// update stuff in redis
// delete tweet queue record
// set queue record to 'error' if necessary
package main

import (
	"database/sql"
	"fmt"
	"github.com/golang/protobuf/proto"
	_ "github.com/lib/pq"
	"gopkg.in/redis.v4"
	pb "local/twit/protobuf"
	st "local/twit/struct"
	"local/twit/twitutil"
	"log"
	"strconv"
)

func sendToRedis(recipientId int, tweet st.Tweet) {
	// is there a better way to do this? i.e. like a tweet without a message?
	tweetLite := &pb.TweetLite{
		Id:     proto.Int(tweet.Id),
		UserId: proto.Int(tweet.UserId),
	}
	pbTweetLite, err := proto.Marshal(tweetLite)
	if err != nil {
		log.Fatalln("Failed to encode tweet:", err)

	}
	//if err := ioutil.WriteFile(fname, out, 0644); err != nil {
	//	log.Fatalln("Failed to write address book:", err)
	//}
	fmt.Println("sending to recipient ", recipientId)
	fmt.Println(tweetLite)

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	// need to check client err???

	recipientIdStr := strconv.Itoa(recipientId)
	err1 := client.LPush(recipientIdStr, pbTweetLite).Err()
	if err1 != nil {
		panic(err1)
	}
}

func getHomeTimelineFromRedis(recipientId int) []pb.TweetLite {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	recipientIdStr := strconv.Itoa(recipientId)
	result, err := client.LRange(recipientIdStr, 0, -1).Result()
	if err != nil {
		panic(err)
	}

	tweetLites := []pb.TweetLite{}
	for _, bufthing := range result {
		tweetLite := &pb.TweetLite{}
		b := []byte(bufthing)
		err := proto.Unmarshal(b, tweetLite)
		if err != nil {
			log.Fatalln("Failed to decode tweet:", err)
		}
		tweetLites = append(tweetLites, *tweetLite)
	}
	return tweetLites
}

func qryFollowerIds(tweet st.Tweet) []int {
	dbDriver, dbOpen := twitutil.GetDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		log.Fatal(err)
	}

	q := fmt.Sprintf("SELECT id FROM t_follower WHERE followed_id=%d", tweet.UserId)
	rows, err := db.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	followerIdSlice := []int{}
	for rows.Next() {
		var followerId int
		err := rows.Scan(&followerId)
		if err != nil {
			log.Fatal(err)
		}
		followerIdSlice = append(followerIdSlice, followerId)
	}

	return followerIdSlice

}

func fanout(tweet st.Tweet) {
	followerIdSlice := qryFollowerIds(tweet)
	for _, followerId := range followerIdSlice {
		sendToRedis(followerId, tweet)
	}

}

func getNextTweetId() int {
	dbDriver, dbOpen := twitutil.GetDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT tweet_id FROM t_tweet_queue WHERE status='ready' LIMIT 1")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tweetId int
	for rows.Next() {
		err := rows.Scan(&tweetId)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("got next tweet! tweet id is: ", tweetId)

	return tweetId

}

func main() {
	for {
		tweetId := getNextTweetId()

		// what to do if tweetid is null??

		// warning: reusing logic from GetTweet()
		dbDriver, dbOpen := twitutil.GetDbConfig()
		db, err := sql.Open(dbDriver, dbOpen)
		if err != nil {
			log.Fatal(err)
		}
		row := db.QueryRow("SELECT id, user_id, message FROM t_tweet WHERE id = $1", tweetId)
		tweet := st.Tweet{}
		err2 := row.Scan(&tweet.Id, &tweet.UserId, &tweet.Message)
		if err2 != nil {
			log.Fatal(err2)
		}
		fanout(tweet)

		//	getHomeTimelineFromRedis(2)
		//	fmt.Println("done getting home timeline!")

	}
}
