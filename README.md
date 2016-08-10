# Twit

Twit is a toy implementation of Twitter's (old) tweet ingestion and delivery
backend, based on the architecture described in
[this High Scalability article](http://highscalability.com/blog/2013/7/8/the-architecture-twitter-uses-to-deal-with-150m-active-users.html).
My motivation was to learn a few technologies that were new to me
(Go, Redis, protocol buffers) and practice some skills that can
always use honing (API design, web app architecture, testing).

Twit's only functionality is to create tweets and deliver them to
the appropriate users. Anytime a user tweets, that tweet is added to his or her
followers' home timelines; a home timeline is a list of all the
tweets of all the people that user follows. Twit does NOT allow you
to create users, follow other users, or perform any other Twitter functions.
A few users and follow relationships are added at startup to give you something
to play with. Twit also does not have a frontend, so the API is the only way to
interact with it. Like I said, this is a toy for learning, not a real system :-).


## Installation

### Get the code

```
$ go get github.com/katrinae/twit
```

### Set up Redis

If you want to actually run Twit, you'll need to have Redis installed.
See the [Redis homepage](http://www.redis.io/) if you need help with installation.
Once it's installed, you can run it with:

```
$ redis
```

### Set up Postgres

You'll also need to have Postgres installed to run Twit. See [this tutorial](https://www.codefellows.org/blog/three-battle-tested-ways-to-install-postgresql/)
for installation help. Once you have Postgres up and running, you can create your development database with
the `createdb` command-line tool:

```
$ createdb {DB_NAME} --owner {DB_USER}
```

#### Create the Postgres schema

Twit uses [goose](https://bitbucket.org/liamstask/goose/) to manage its database schema. To install goose and create the database tables, run:

```
$ go get bitbucket.org/liamstask/goose/cmd/goose
$ goose up
```

#### Import test data

Because Twit does not allow you to create or modify users, it ships with some test data.
You can import this data into Postgres from with psql, the Postgres shell.
Replace $GOPATH with your actual gopath; environment variables are not recognized in the shell.

```
$ psql --dbname=DB_NAME --username=DB_USER
=# \i $GOPATH/github.com/katrinae/twit/db/createTestRecords.sql;
```

### Customize your conf file

Edit the file *db/dbconf.yml* with the parameters for the Redis and Postgres
instances you will be using. At the very least, you will want to change DB_NAME
and DB_USER to match your environment.

When you are done editing, remove the file from
your git worktree with:

```
$ git update-index --assume-unchanged db/dbconf.yml
```

(This is to prevent you from accidentally pushing your private credentials
somewhere).

## API

Because my goal was to learn about backend systems, Twit has no web frontend,
just a HTTP API. The endpoints are:

   * **GET /tweets** - returns a list of all tweets in the system
   * **POST /tweets** - creates a new tweet
   * **GET /tweets/:tweetId** - returns a specific tweet
   * **DEL /tweets/:tweetId** - deletes a specific tweet
   * **GET /tweets/user_timeline/:userId** - returns a list of all the tweets
     posted by the specified user
   * **GET /tweets/home_timeline/:userId** - returns a user's home timeline (posts made by the users he follows)

All endpoints return JSON. They are defined in *routes.go*.

## Internal Architecture

### Background

If Twitter wasn't facing any
scalability challenges, they could easily store and deliver tweets from
a relational database, using a schema like:

**User table**

| id | username |
|----|----------|
| 1 | Bob |
| 2 | Jane |
| 3 | Sue |

**Follow table**

| id | follower_id | followed_id |
|----|-------------|-------------|
| 1 | 3 | 1 |
| 2 | 3 | 2 |


**Tweet table**

| id | user_id | message |
|----|---------|---------|
| 1 | 1 | Hello world, i'm Bob |
| 2 | 2 | Hello world, i'm Jane |

Then, in order to build the home timeline for user 3 (with username Sue), Twitter would
only need to run a single SQL query:

```
SELECT tweet.id, user.username, tweet.message FROM user
INNER JOIN follow ON user.id = follow.followed_id
INNER JOIN tweet ON tweet.user_id = follow.followed_id
WHERE follow.follower_id = 3
```

Because Sue is following Bob and Jane, who have each made one tweet, the API would return:

```
{
    {
		"userId" : "1",
		"userName" : "Bob",
		"message" : "Hello world, i'm Bob"
    },
    {
		"userId" : "2",
		"userName" : "Jane",
		"message" : "Hello world, i'm Jane"
    }
}
```

But a triple-join query didn't scale very well,
so Twitter decided to do something different. Instead of assembling a home timeline's
tweets each time it is *read*,
Twitter updates home timelines anytime a new tweet is *written*.
All users' home timelines are stored in a Redis key-value store. When a user requests
their home timeline, Twitter fetches the tweet list from Redis with a single command.
This denormalized, compute-on-write system
allows them to return content faster and scale their system more effectively.

### Twit Architecture

Like Twitter, Twit stores home timelines in Redis and updates them on write. On read, tweets
are retrieved from Redis and additional data is added to them from Postgres.

The entry point is an HTTP server that listens for requests to create, delete, and retrieve tweets.

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

   * Home timeline endpoint handler: `homeTimelineTweets()` in *routes.go*

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
[go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) package.

   * Test file: *pgconn_test.go*
   
Its Redis-related
functions are integration tested using a test Redis instance defined in *db/dbconf.yml*.

   * Test file: *redisconn_test.go*
   
Run all tests by running `go test` from the command line.

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
   * *db/createTestRecords.sql* - SQL commands to pre-populate the database with user,
     tweet, and follow records.
