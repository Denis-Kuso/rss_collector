-- +goose Up 
ALTER TABLE feed_follows
	ADD CONSTRAINT fk_userID_feedFollows
		FOREIGN KEY (user_id)
		REFERENCES users(id)
		ON DELETE CASCADE,
	ADD CONSTRAINT fk_feedID_feedFollows
		FOREIGN KEY (feed_id)
		REFERENCES feeds(id)
		ON DELETE CASCADE,
	ADD CONSTRAINT UC_feedFollow UNIQUE(user_id, feed_id);

-- +goose Down
ALTER TABLE feed_follows
	DROP CONSTRAINT fk_userID_feedFollows,
	DROP CONSTRAINT fk_feedID_feedFollows,
	DROP CONSTRAINT UC_feedFollow;
