package twit

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"testing"
)

var databaseFixtureFile string

func init() {
	flag.StringVar(&databaseFixtureFile, "dbfixtures", "dbfixtures.sql",
		"database fixtures")
}

func setUpFixtures() {
	tearDownFixtures() // just in case
	importTestData()
}

func tearDownFixtures() {
	dbDriver, dbOpen := getDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		panic(fmt.Errorf("Could not connect to database for teardown"))
	}
	db.Exec("DELETE FROM t_user")
}

func importTestData() {
	dbDriver, dbOpen := getDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		panic(fmt.Errorf("%+v", err))
	}

	sqlFixturesBuf, err := os.Open(databaseFixtureFile)
	if err != nil {
		panic(fmt.Errorf("Failed to open database fixtures file %s\n",
			databaseFixtureFile))
	}
	scanner := bufio.NewScanner(sqlFixturesBuf)
	for scanner.Scan() {
		sql := scanner.Text()
		_, err = db.Exec(sql)
		if err != nil {
			tearDownFixtures()
			panic(fmt.Errorf("SQL statement failed: %s", sql))
		}
	}
}

func TestMain(m *testing.M) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("../../")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	os.Exit(m.Run())
}

func TestDbInsertTweet(t *testing.T) {
	setUpFixtures()
	msg := "My first tweet!"
	tweetFixture := Tweet{
		UserId:  3,
		Message: msg,
	}
	tweet, err := dbInsertTweet(3, msg)
	if err != nil {
		t.Logf("%+v", err)
		tearDownFixtures()
		t.FailNow()
	}
	tweetFixture.Id = tweet.Id
	if tweet != tweetFixture {
		t.Logf("Wanted: %+v -- Got: %+v", tweetFixture, tweet)
		tearDownFixtures()
		t.FailNow()
	}
	tearDownFixtures()
}

func TestDbGetTweet(t *testing.T) {
	setUpFixtures()
	tweetId := 1
	tweetFixture := Tweet{
		Id:      tweetId,
		UserId:  1,
		Message: "i have lots of followers",
	}
	tweet, err := dbGetTweet(tweetId)
	if err != nil {
		t.Logf("%+v", err)
		tearDownFixtures()
		t.FailNow()
	}
	if tweet != tweetFixture {
		t.Logf("Wanted: %+v -- Got: %+v", tweetFixture, tweet)
		tearDownFixtures()
		t.FailNow()
	}
	tearDownFixtures()
}

func TestDbGetTweetNone(t *testing.T) {
	setUpFixtures()
	tweetId := 500 // not in database
	_, err := dbGetTweet(tweetId)
	if err == sql.ErrNoRows {
		tearDownFixtures()
		return
	}
	if err != nil {
		t.Logf("%+v", err)
		tearDownFixtures()
		t.FailNow()
	}
	t.Logf("sql.ErrNoRows not thrown")
	tearDownFixtures()
	t.FailNow()

}

// func TestDbGetTweetMalformed(t *testing.T)    {}
// func TestDbGetTweetSqlInjection(t *testing.T) {}

func TestDbQryAllTweets(t *testing.T) {
	setUpFixtures()
	tweetFixtureA := []Tweet{
		Tweet{Id: 1, UserId: 1, Message: "i have lots of followers"},
		Tweet{Id: 2, UserId: 2, Message: "i have no followers"},
	}
	tweetA, err := dbQryAllTweets()
	if err != nil {
		t.Logf("%+v", err)
		tearDownFixtures()
		t.FailNow()
	}

	if len(tweetA) != len(tweetFixtureA) {
		t.Logf("tweetA and tweetFixtureA are not the same length")
		t.Logf("tweetA: %+v", tweetA)
		t.Logf("tweetFixtureA: %+v", tweetFixtureA)
		tearDownFixtures()
		t.FailNow()
	}

	for i, _ := range tweetA {
		if tweetA[i] != tweetFixtureA[i] {
			t.Logf("Wanted: %+v -- Got: %+v",
				tweetFixtureA[i], tweetA[i])
			tearDownFixtures()
			t.FailNow()
		}
	}
	tearDownFixtures()
}
func TestDbQryAllTweetsNone(t *testing.T) {
	// no fixture setup b/c don't want anything in DB
	// setUpFixtures()
	tweetA, err := dbQryAllTweets()
	if err != nil {
		t.Logf("%+v", err)
		t.FailNow()
	}
	if len(tweetA) != 0 {
		t.Logf("Too many tweets returned: %+v", tweetA)
		t.FailNow()
	}
}

func TestDbQryUserTweets(t *testing.T) {
	setUpFixtures()
	tweetFixtureA := []Tweet{
		Tweet{Id: 1, UserId: 1, Message: "i have lots of followers"},
	}
	userId := 1
	tweetA, err := dbQryUserTweets(userId)
	if err != nil {
		t.Logf("%+v", err)
		tearDownFixtures()
		t.FailNow()
	}

	if len(tweetA) != len(tweetFixtureA) {
		t.Logf("tweetA and tweetFixtureA are not the same length")
		t.Logf("tweetA: %+v", tweetA)
		t.Logf("tweetFixtureA: %+v", tweetFixtureA)
		tearDownFixtures()
		t.FailNow()
	}

	for i, _ := range tweetA {
		if tweetA[i] != tweetFixtureA[i] {
			t.Logf("Wanted: %+v -- Got: %+v",
				tweetFixtureA[i], tweetA[i])
			tearDownFixtures()
			t.FailNow()
		}
	}
	tearDownFixtures()
}

func TestDbQryUserTweetsNone(t *testing.T) {
	// no fixture setup b/c don't want anything in DB
	// setUpFixtures()
	userId := 500
	tweetA, err := dbQryUserTweets(userId)
	if err != nil {
		t.Logf("%+v", err)
		t.FailNow()
	}
	if len(tweetA) != 0 {
		t.Logf("Too many tweets returned: %+v", tweetA)
		t.FailNow()
	}
}

func TestDbDelTweet(t *testing.T) {
	setUpFixtures()
	tweetId := 1
	err := dbDelTweet(tweetId)
	if err != nil {
		t.Logf("%+v", err)
		tearDownFixtures()
		t.FailNow()
	}

	// attempt to re-get
	dbDriver, dbOpen := getDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		t.Logf("%+v", err)
		tearDownFixtures()
		t.FailNow()
	}
	tweet := Tweet{}
	err = db.QueryRow("SELECT * FROM t_tweet WHERE id=$1", tweetId).Scan(&tweet)
	if err == sql.ErrNoRows {
		return
	}
	if err != nil {
		t.Logf("%+v", err)
		tearDownFixtures()
		t.FailNow()
	}
	t.Logf("sql.ErrNoRows not thrown after DELETE")
	tearDownFixtures()
}
func TestDbDelTweetNone(t *testing.T) {
	tweetId := 500
	err := dbDelTweet(tweetId)
	if err != nil {
		t.Logf("%+v", err)
		t.FailNow()
	}

	// attempt to re-get
	dbDriver, dbOpen := getDbConfig()
	db, err := sql.Open(dbDriver, dbOpen)
	if err != nil {
		t.Logf("%+v", err)
		t.FailNow()
	}
	tweet := Tweet{}
	err = db.QueryRow("SELECT * FROM t_tweet WHERE id=$1", tweetId).Scan(&tweet)
	if err == sql.ErrNoRows {
		return
	}
	if err != nil {
		t.Logf("%+v", err)
		t.FailNow()
	}
	t.Logf("sql.ErrNoRows not thrown after DELETE")
}

/*
func TestDbEnqueueTweetId(t *testing.T)               {}
func TestDbEnqueueTweetIdMalformed(t *testing.T)      {}
func TestDbGetNextQueuedTweetId(t *testing.T)         {}
func TestDbMarkTweetProcessing(t *testing.T)          {}
func TestDbMarkTweetProcessingNone(t *testing.T)      {}
func TestDbMarkTweetProcessingMalformed(t *testing.T) {}
func TestDbDequeueTweetId(t *testing.T)               {}
func TestDbDequeueTweetIdNone(t *testing.T)           {}
func TestDbDequeueTweetIdMalformed(t *testing.T)      {}
*/

func TestDbQryUserFollowers(t *testing.T) {
	setUpFixtures()
	userId := 1
	followFixtureA := []Follow{
		Follow{Id: 1, FollowerId: 2, FollowedId: 1},
		Follow{Id: 2, FollowerId: 3, FollowedId: 1},
		Follow{Id: 3, FollowerId: 4, FollowedId: 1},
	}
	followA, err := dbQryUserFollowers(userId)
	if err != nil {
		t.Logf("%+v", err)
		tearDownFixtures()
		t.FailNow()
	}
	if len(followA) != len(followFixtureA) {
		t.Logf("followA and followFixtureA are not the same length")
		t.Logf("followA: %+v", followA)
		t.Logf("followFixtureA: %+v", followFixtureA)
		tearDownFixtures()
		t.FailNow()
	}
	for i, _ := range followA {
		if followA[i] != followFixtureA[i] {
			t.Logf("Wanted: %+v -- Got: %+v",
				followFixtureA, followA)
			tearDownFixtures()
			t.FailNow()
		}
	}
	tearDownFixtures()
}

func TestDbQryUserFollowersNone(t *testing.T) {
	// no fixture setup b/c don't want anything in DB
	// setUpFixtures()
	userId := 500
	followA, err := dbQryUserFollowers(userId)
	if err != nil {
		t.Logf("%+v", err)
		t.FailNow()
	}
	if len(followA) != 0 {
		t.Logf("Too many followers returned: %+v", followA)
		t.FailNow()
	}
}

// func TestDbQryFollowsMalformed(t *testing.T) {}
