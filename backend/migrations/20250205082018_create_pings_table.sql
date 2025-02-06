-- +goose Up
-- +goose StatementBegin
-- Таблица контейнеров
CREATE TABLE containers (
    id SERIAL PRIMARY KEY,
    id_container VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL
);

-- Таблица пингов
CREATE TABLE pings (
    id SERIAL PRIMARY KEY,
    id_container VARCHAR(255) REFERENCES containers(id_container) ON DELETE CASCADE,
    ip VARCHAR(255) NOT NULL,
    status BOOLEAN NOT NULL,
    response_time FLOAT NOT NULL,
    last_success TIMESTAMP WITH TIME ZONE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pings;
DROP TABLE IF EXISTS containers;
-- +goose StatementEnd
