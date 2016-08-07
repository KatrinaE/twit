# Twit

Twit is a toy implementation of Twitter's (old) tweet ingestion and delivery
backend, based on the architecture described in
[this High Scalability article](http://highscalability.com/blog/2013/7/8/the-architecture-twitter-uses-to-deal-with-150m-active-users.html).
My motivation was to learn a few technologies that were new to me
(Go, Redis, protocol buffers) and practice some skills that can
always use honing (API design, web app architecture, testing).
I enjoyed learning how Twitter solved
the technical challenges associated with tweet delivery, because I
had wrestled with similar issues while building a user notification
system at work.

Twit's only functionality is to create tweets and deliver them to
the appropriate users (anytime a user tweets, that tweet is added to his or her
followers' home timelines. A user's home timeline is a list of all the
tweets of all the people that user follows). Twit does NOT allow you
to create users, follow other users, or perform any other Twitter functions.
Instead, a few users and follow relationships are added at startup.
Twit also does not have a frontend, so the API is the only way to play with it.
Like I said, this is a toy :-).


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

Twitter's core functionality is pretty simple: save each users' tweets
and return them in their followers' home timelines. If they didn't have any
scalability issues, they could easily store and deliver tweets from
a relational database, using a schema like:

**User table**

| id | username |
|----|----------|
| 1 | bob |
| 2 | jane |
| 3 | sue |

**Follow table**

| id | follower_id | followed_id |
|----|-------------|-------------|
| 1 | 3 | 1 |
| 2 | 3 | 2 |


**Tweet table**

| id | user_id | message |
|----|---------|---------|
| 1 | 1 | hello world, i'm bob |
| 2 | 2 | hello world, i'm jane |

Then, in order to build the home timeline for user #3, Twitter would
only need to run a single SQL query:

```
SELECT user.id, user.username, tweet.message FROM user
INNER JOIN tweet ON user.id = tweet.user_id
INNER JOIN follow ON user.id = follow.follower_id
WHERE follow.follower_id = 3;
```

The API would return results like:

```
{
    {
	"userId" : "1",
	"userName" : "bob",
	"message" : "hello world, i'm bob"
    },
    {
	"userId" : "2",
	"userName" : "jane",
	"message" : "hello world, i'm jane"
    }
}
```

But a triple-join query in a relational database did not scale very well,
so Twitter decided to do something different. Instead of determining which
tweets are in a user's home timeline each time the timeline is *read*,
Twitter updates home timelines anytime a new tweet is *written*.
All users' home timelines are stored in a Redis key-value store, where the
keys are user IDs and the values are lists of tweets. When a user requests
their home timeline, Twitter already knows which tweets to show them.
This denormalized, compute-on-write system
allows them to return content faster and scale their system more effectively.

### Twit Architecture

Twit includes an HTTP server that listens for requests to create, delete, and retrieve tweets.

   * Server command: *cmd/server/main.go*
   * Endpoint URLs and definitions: *routes.go*

User, follow, and tweet message data is stored in a Postgres
database similar to the one described above.

   * Postgres functions: *pgconn.go*
   * Postgres data models: *models.go*

Home timelines are stored in
Redis lists whose keys are recipient user IDs and whose values are lists of tweets.
To save space, the tweet list does not contain the actual tweet messages - only the
tweet ID and the posting user's ID, which
are serialized into a [protocol buffer](https://developers.google.com/protocol-buffers/).

   * Redis functions: *redisconn.go*
   * Protocol buffers: *tweetlite.proto* (source), *tweetlite.pb.go* (compiled)

Home timelines are updated during a process called "fanout", which is
initiated when a tweet is created. When a user sends a tweet to Twit,
the tweet is inserted into a queue; later,
a fanout worker pulls the tweet off the queue and inserts it into the
tweeter's followers' home timelines. The queue is
implemented in Redis, using its
[reliable queue](http://redis.io/commands/rpoplpush) pattern.

   * Fanout command: *cmd/fanoutworker/main.go*
   * Fanout functions: *fanout.go*
   * Queue functions: *redisconn.go*

When a user's home timeline is requested, it is retrieved from Redis and
deserialized into an array of Go structs. Then, display-related fields like
the tweeter's username and the tweet message are retrieved from Postgres
and added to each tweet record.
(Twitter calls this process "hydration"). The entire thing is then
formatted as JSON and sent to the requester via HTTP.


*Note*: 
Twit architecture resembles the Twitter architecture, though it does not
match it perfectly. For example, Twitter used MySql - fronted with
Memcached - for disk storage, and stored follow relationships in a graph
database; Twit uses Postgres for all of this information. Twitter doesn't
describe the specifics of how they serialized their tweets before stuffing
them in Redis, only to say
that they used some extra bytes to store metadata, like whether a tweet is
a retweet. Twit uses protocol buffers without any additional metadata.

## Tests

Twit's Postgres-related functions are unittested with mocks from the
[go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) package. Its Redis-related
functions are integration tested using a test Redis instance defined in `dbconf.yml`.

## File Glossary

Executables

   * *cmd/fanoutworker/main.go* - tweet delivery worker
   * *cmd/server/main.go* - web API server

Source Code

   * *fanout.go* - functions that deliver tweets to users' home timelines
   * *models.go* - structs mirroring Postgres database schema
   * *pgconn.go* - functions that construct and execute Postgres queries
   * *redisconn.go* - functions that construct and execute Redis queries
   * *tweetlite.proto*, *tweetlite.pb.go* - tweet protocol buffer definition
   * *util.go* - utility functions, mostly related to retrieving configs
     and writing responses.

Tests

   * *pgconn_tests.go*
   * *redisconn_tests.go*

Database Management

   * *db/dbconf.yml* - Postgres and Redis configurations
   * *db/migrations* - Postgres migrations (managed by [goose](http://bitbucket.org/liamstask/goose)).
   * *db/create_test_records.sql* - SQL commands to pre-populate the database with user,
     tweet, and follow records.
