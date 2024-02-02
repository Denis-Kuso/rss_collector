-- +goose Up
CREATE TABLE users (id UUID NOT NULL PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	name char(50) NOT NULL
	api_key VARCHAR(64) NOT NULL UNIQUE DEFAULT encode(sha256(random()::text::bytea), 'hex'));	

-- +goose Down
DROP TABLE users;
