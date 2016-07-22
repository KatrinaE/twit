
-- +goose Up
CREATE TABLE t_tweet_queue (
tweet_id	INTEGER REFERENCES t_tweet(id),
status		VARCHAR(16)
);


-- +goose Down
DROP TABLE t_tweet_queue;

