# Twit

Twit is a toy implementation of Twitter's (old) tweet ingestion and delivery
backend, based on the architecture described in
[this High Scalability article](http://highscalability.com/blog/2013/7/8/the-architecture-twitter-uses-to-deal-with-150m-active-users.html).
My motivation was to learn a few technologies that were new to me
(Go, Redis, protocol buffers) and practice some concepts that I already
was familiar with (API design, web app architecture, testing,
deployment). I also enjoyed learning how Twitter solved
the technical challenges associated with tweet delivery, because I
had wrestled with similar issues while building a user notification
system at work.

Twit's only functionality is to create tweets and deliver them to
the appropriate users (anytime a user tweets, that tweet is added to his or her
followers' home timelines. A user's home timeline is a list of all the
tweets of all the people that user follows). Twit assumes that users
and their 'follow'
relationships were defined at startup. Like I said, this is a toy :-).


## Installation

TODO

## API

Because my goal was to learn about backend systems, Twit has no web frontend,
just a HTTP API. The endpoints are:

   * **GET /tweets** - returns a list of all tweets in the system
   * **POST /tweets** - creates a new tweet
   * **GET /tweets/:tweet_id** - returns a specific tweet
   * **DEL /tweets/:tweet_id** - deletes a specific tweet
   * **GET /tweets/user/:userId** - returns a list of all the tweets
     posted by the specified user
   * **GET /tweets/followed/:userId** - returns a user's home timeline.

All endpoints return JSON.

The endpoints are defined in *routes.go*.

## Internal Architecture

### Background

Twitter's core functionality is pretty simple. If they didn't have any
scalability issues, they could easily store and deliver tweets from
a relational database, using a schema like:

**User table**

| id | username |
-----------------
| | |

**Follow table**

| id | follower_id | followed_id |
---------------------------------
| | |


**Tweet table**

| id | user_id | message |
-------------------------
| | |

Then, in order to build the home timeline for user #100, Twitter would
only need to run a single SQL query:

```
SELECT user.id, user.username, tweet.message FROM user
INNER JOIN tweet ON user.id = tweet.user_id
INNER JOIN follow user.id = follow.follower_id
WHERE follow.follower_id = 100;
```

The API would return results like:

```
{
    {
	"userId" : "12345",
	"userName" : "bob",
	"message" : "hello world, i'm bob"
    },
    {
	"userId" : "54321",
	"userName" : "jane",
	"message" : "hello world, i'm jane"
    }
}
```

But a triple-join query in a relational database did not scale very well,
so Twitter decided to do something different. Instead of determining which
tweets are in a user's home timeline each time the timeline is read,
Twitter stores the tweet ids of each user's home timeline in Redis, updating
the timelines every time a tweet is created. This compute-on-write system
allows them to return content faster and scale their system more effectively.

Twit architecture resembles the Twitter architecture, though it does not
match it perfectly. For example, Twitter used MySql - fronted with
Memcached - for disk storage, and stored follow relationships in a graph
database; Twit uses Postgres for all of this information. Twitter doesn't
describe the specifics of how they serialized their tweets, only to say
that they used some extra bytes to store metadata, like whether a tweet is
a retweet. Twit uses protocol buffers without any additional metadata.


### Twit Architecture

In Twit, user, follow, and tweet message data is stored in a Postgres
database similar to the one described above. Home timelines are stored in
Redis lists, in the following format:

   * Redis key: recipient ID
   * Redis value: list of tweets

To save space, the tweet list only contains tweet IDs and user IDs, which
are serialized into [protocol buffers](https://developers.google.com/protocol-buffers/).

Home timelines are updated during a process called "fanout". Fanout is
initiated when a tweet is created, but the actual delivery happens
asynchronously. Upon creation, the tweet is inserted into a queue; later,
a fanout worker pulls the tweet off the queue and inserts it into the
tweeter's followers' home timelines.

When a user's home timeline is requested, it is retrieved from Redis and
deserialized into an array of Go structs. Then, display-related fields
(tweeter username, the actual tweet texts) are retrieved from Postgres
(Twitter calls this process "hydration"). The entire thing is then
formatted as JSON and sent to the requester via HTTP.

## Tests

TODO

## File Glossary

Commands

   * *cmd/fanoutworker/fanoutworker.go* - tweet delivery worker
   * *cmd/server/server.go* - API web server

Source Code

   * *internal/fanout.go* - functions that deliver tweets to users' home timelines
   * *internal/models.go* - structs mirroring Postgres database schema
   * *internal/pgconn.go* - functions that construct and execute Postgres queries
   * *internal/redisconn.go* - functions that construct and execute Redis queries
   * *internal/tweetlite.pb.go*, *tweetlite.proto* - tweet protocol buffer definition
   * *internal/util.go* - utility functions, mostly related to retrieving configs
     and writing responses.

Tests

   * *internal/pgconn_tests.go*
   * *internal/redisconn_tests.go*

Database Management

   * *db/dbconf.yml* - specifies Postgres and Redis configurations
   * *db/migrations* - directory holding Postgres migrations (managed by [goose](http://bitbucket.org/liamstask/goose)).
   * *db/create_test_records.sql* - SQL commands to pre-populate the database with user,
     tweet, and follow records.
