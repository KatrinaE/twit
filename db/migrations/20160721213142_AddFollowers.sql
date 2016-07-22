
-- +goose Up
INSERT INTO t_follower (id, follower_id, followed_id) VALUES (1, 2, 1);
INSERT INTO t_follower (id, follower_id, followed_id) VALUES (2, 3, 1);
INSERT INTO t_follower (id, follower_id, followed_id) VALUES (3, 4, 1);

-- +goose Down
DELETE FROM t_follower WHERE id IN('1', '2', '3');
