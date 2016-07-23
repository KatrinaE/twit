
-- +goose Up
INSERT INTO t_user (id, username) VALUES (1, 'ben');
INSERT INTO t_user (id, username) VALUES (2, 'alyssa');
INSERT INTO t_user (id, username) VALUES (3, 'louis');
INSERT INTO t_user (id, username) VALUES (4, 'eva');
INSERT INTO t_user (id, username) VALUES (5, 'cy');


-- +goose Down
DELETE FROM t_user WHERE id IN('1', '2', '3', '4', '5');

