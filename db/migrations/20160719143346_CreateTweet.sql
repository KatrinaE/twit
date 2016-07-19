
-- +goose Up
CREATE TABLE t_tweet (
id		VARCHAR(16) PRIMARY KEY,
user_id		VARCHAR(16) REFERENCES t_user(id),
message		VARCHAR(160),
ctime		TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
mtime		TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);


-- +goose Down
DROP TABLE t_tweet;
