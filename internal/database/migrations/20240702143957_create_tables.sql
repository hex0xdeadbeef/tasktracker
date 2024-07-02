-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS workers (
    passport_serie VARCHAR(3),
    passport_num VARCHAR(3),
    surname VARCHAR(128) NOT NULL CHECK (CHAR_LENGTH(surname) > 0),
    name VARCHAR(128) NOT NULL CHECK (CHAR_LENGTH(name) > 0),
    address VARCHAR(128) NOT NULL CHECK (CHAR_LENGTH(address) > 0),
    PRIMARY KEY (passport_serie, passport_num)
);

CREATE TABLE IF NOT EXISTS tasks (
	id serial,
	name text NOT NULL CHECK(CHAR_LENGTH(name) > 0),
	started_at timestamp,
	finished_at timestamp,
	is_deleted bool NOT NULL,

	worker_passport_serie VARCHAR(3),
	worker_passport_num VARCHAR(3),

	PRIMARY KEY (id),
	FOREIGN KEY (worker_passport_serie, worker_passport_num)
	REFERENCES workers(passport_serie, passport_num)
	ON DELETE CASCADE
	ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS workers;
-- +goose StatementEnd