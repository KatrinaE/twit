
-- +goose Up
CREATE TABLE t_tweet_queue (
tweet_id	INTEGER REFERENCES t_tweet(id) ON DELETE CASCADE,
status		VARCHAR(16),
ctime		TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
mtime		TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);


-- +goose Down
DROP TABLE t_tweet_queue;

