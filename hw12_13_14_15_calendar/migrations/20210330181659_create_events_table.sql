-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
    id UUID NOT NULL PRIMARY KEY,
    title text,
    start_date timestamp,
    duration bigint,
    description text,
    user_id UUID,
    notify_before bigint
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd
