
-- +goose Up
CREATE TABLE t_follower (
id		SERIAL PRIMARY KEY,
follower_id	INTEGER REFERENCES t_user(id),
followed_id	INTEGER REFERENCES t_user(id),
ctime		TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
mtime		TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);


-- +goose Down
DROP TABLE t_follower;
