-- +goose Up
-- +goose StatementBegin
CREATE TABLE note (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE note;
-- +goose StatementEnd
