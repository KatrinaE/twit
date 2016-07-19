---+ Twit

A tiny toy version of Twitter based on the architecture described at:

http://highscalability.com/blog/2013/7/8/the-architecture-twitter-uses-to-deal-with-150m-active-users.html

---++ Notes on the article

Stored in each user's timeline in Redis:
   * tweet ID of the generated tweet
   * user ID of the originator of the tweet
   * 4 bytes of bits used to mark if it’s a retweet or a reply or something else.

Your home timeline sits in a Redis cluster and is 800 entries long.

Every active user is stored in RAM

If you are inactive (no login in > 30 days), you fall out of RAM.
To reconstruct:
Query against the social graph service. Figure out who you follow.
Hit disk for every single one of them and then shove them back into Redis. 

Disk storage in MySQL

If a tweet is actually a retweet then a pointer is stored to the original tweet.

Since the timeline only contains tweet IDs they must “hydrate” those tweets,
that is find the text of the tweets. Given an array of IDs they can do a
multiget and get the tweets in parallel from T-bird.

Tweetypie has about the last month and half of tweets stored in its memcache cluster.
These are exposed to internal customers.

Each tweet also stored in "Early Bird" machine (modified Lucene) for search
"In fanout a tweet may be stored in N home timelines of how many people are
following you, in Early Bird a tweet is only stored in one Early Bird machine
(except for replication). "

Your activity information is computed on a write basis... similar to the home
timeline, it is a series of IDs of pieces of activity, so there’s favorite ID,
a reply ID, etc.

Discovery is a customized search based on what they know about you.
