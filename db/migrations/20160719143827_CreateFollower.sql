
-- +goose Up
CREATE TABLE t_follower (
id		VARCHAR(16) PRIMARY KEY,
follower_id	VARCHAR(16) REFERENCES t_user(id),
followed_id	VARCHAR(16) REFERENCES t_user(id),
ctime		TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
mtime		TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);


-- +goose Down
DROP TABLE t_follower;
