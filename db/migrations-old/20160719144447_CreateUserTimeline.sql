
-- +goose Up
CREATE TABLE t_user_timeline (
id		SERIAL PRIMARY KEY,
user_id		INTEGER REFERENCES t_user(id) ON DELETE CASCADE,
tweets		VARCHAR(160),
ctime		TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
mtime		TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);



-- +goose Down
DROP TABLE t_user_timeline;

