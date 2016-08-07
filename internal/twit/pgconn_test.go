package twit

import (
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestDbInsertTweet(t *testing.T) {
	tweetId := 1
	userId := 1
	tweetMsg := "Tweet tweet!"
	tweetFixture := Tweet{
		Id:      tweetId,
		UserId:  userId,
		Message: tweetMsg,
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error when opening stub db connection: %+v", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "user_id", "message"}).
		AddRow(tweetId, userId, tweetMsg)
	mock.ExpectQuery("INSERT INTO t_tweet").
		WithArgs(userId, tweetMsg).WillReturnRows(rows)

	tweet, err := dbInsertTweet(db, userId, tweetMsg)
	if err != nil {
		t.Fatalf("tweet: %+v, err: %+v", tweet, err)
	}
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if tweet != tweetFixture {
		t.Fatalf("Tweet %+v does not match fixture %+v", tweet, tweetFixture)
	}
}
