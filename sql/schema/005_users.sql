-- +goose Up 
ALTER TABLE feed_follows
	ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT TO_TIMESTAMP('2017-03-31 9:30:20',
    'YYYY-MM-DD HH:MI:SS'),
	ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT TO_TIMESTAMP('2017-03-31 9:30:20','YYYY-MM-DD HH:MI:SS'),
	DROP CONSTRAINT feed_follows_feed_id_fkey,
	DROP CONSTRAINT feed_follows_user_id_fkey;

-- -goose Down
ALTER TABLE feed_follows
	DROP COLUMN created_at,
	DROP COLUMN updated_at,
	ADD FOREIGN KEY (user_id) REFERENCES users(id),
	ADD FOREIGN KEY (feed_id) REFERENCES feeds(id);

