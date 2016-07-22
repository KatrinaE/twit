package twit

import (
    "log"
    "fmt"
    "github.com/golang/protobuf/proto"
    "local/twit/protobuf"
    "local/twit/twitutil"
)

type MahTimeline {
	protobuf.Timeline
	addTweet(protobuf.Tweet) error
}

func main() {

	userId := twit.makeUUID()
	tweetId := twit.makeUUID()

	timelineTweet := &protobuf.Timeline_Tweet {
		Id: proto.String("hello"),
	        UserId:  proto.String(userId.String()),
	}

	timeline := &protobuf.Timeline {
	   Tweets: []*protobuf.Timeline_Tweet { timelineTweet },
	}
	data, err := proto.Marshal(timeline)
	if err != nil {
	    log.Fatal("marshaling error: ", err)
	}

	timelineTweet2 := &protobuf.Timeline_Tweet {
		Id: proto.String("goodbye"),
	        UserId:  proto.String(userId.String()),
	}

	newThing := append(timeline.Tweets, timelineTweet2)

	newTimeline := &protobuf.Timeline{}
	err = proto.Unmarshal(data, newTimeline)
	if err != nil {
	    log.Fatal("unmarshaling error: ", err)
	}
	fmt.Println(timeline)
	fmt.Println(newTimeline)
	// Now test and newTest contain the same data.
//	if timeline.GetTweets() != newTimeline.GetTweets() {
//	    log.Fatalf("data mismatch %q != %q", timeline.GetTweets(), newTimeline.GetTweets())
//	}
	// etc.
//    log.Printf("Unmarshalled to: %+v", newTimeline)
}
