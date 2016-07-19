
-- +goose Up
CREATE TABLE t_user_timeline (
id		VARCHAR(16) PRIMARY KEY,
user_id		VARCHAR(16) REFERENCES t_user(id),
tweets		VARCHAR(160),
ctime		TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
mtime		TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL

)



-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

