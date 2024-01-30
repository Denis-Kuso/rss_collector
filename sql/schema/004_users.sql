-- +goose Up 
CREATE TABLE feed_follows(
	user_id UUID NOT NULL,
	feed_id UUID NOT NULL,
	ID_FF UUID NOT NULL PRIMARY KEY,
	FOREIGN KEY (user_ID) REFERENCES users(id),
	FOREIGN KEY (feed_id) REFERENCES feeds(id)
);

-- +goose Down
DROP TABLE feed_follows
