-- +goose Up
-- +goose StatementBegin
CREATE TABLE pings (
    id SERIAL PRIMARY KEY,
    ip VARCHAR(255),
    status BOOLEAN,
    response_time FLOAT,
    last_success TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pings;
-- +goose StatementEnd
