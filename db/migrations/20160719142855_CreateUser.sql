
-- +goose Up
CREATE TABLE t_user (
id		VARCHAR(16) PRIMARY KEY,
username	VARCHAR(30) NOT NULL,
ctime		TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
mtime		TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);


-- +goose Down
DROP TABLE t_user;
