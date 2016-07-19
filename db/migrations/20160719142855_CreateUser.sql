
-- +goose Up
CREATE TABLE t_user (
id		VARCHAR(16) PRIMARY KEY,
username	VARCHAR(30) NOT NULL,
"password"	VARCHAR(48) NOT NULL,
email		VARCHAR(100),
first_name	VARCHAR(30),
last_name	VARCHAR(30),
ctime		TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
mtime		TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL
);


-- +goose Down
DROP TABLE t_user;
