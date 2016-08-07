INSERT INTO t_user (id, username) VALUES (1, 'ben');
INSERT INTO t_user (id, username) VALUES (2, 'alyssa');
INSERT INTO t_user (id, username) VALUES (3, 'louis');
INSERT INTO t_user (id, username) VALUES (4, 'eva');
INSERT INTO t_user (id, username) VALUES (5, 'cy');

INSERT INTO t_follow (id, follower_id, followed_id) VALUES (1, 2, 1);
INSERT INTO t_follow (id, follower_id, followed_id) VALUES (2, 3, 1);
INSERT INTO t_follow (id, follower_id, followed_id) VALUES (3, 4, 1);

INSERT INTO t_tweet (id, user_id, message) VALUES (1, 1, 'i have lots of followers');
INSERT INTO t_tweet (id, user_id, message) VALUES (2, 1, 'my second tweet');
INSERT INTO t_tweet (id, user_id, message) VALUES (3, 2, 'i have no followers');

INSERT INTO t_tweet_queue (tweet_id, status) VALUES (1, 'ready');
INSERT INTO t_tweet_queue (tweet_id, status) VALUES (2, 'ready');
INSERT INTO t_tweet_queue (tweet_id, status) VALUES (3, 'ready');

-- DELETE FROM t_tweet_queue;
-- DELETE FROM t_tweet;
-- DELETE FROM t_follow;
-- DELETE FROM t_user;
