-- +goose Up
ALTER table users
	ALTER COLUMN name TYPE VARCHAR(50);

-- +goose Down
ALTER table users
	ALTER COLUMN name TYPE char(50);
