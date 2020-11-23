-- +goose Up
-- +goose StatementBegin
CREATE TABLE time_slots (
    id UUID NOT NULL,
    weekday INTEGER NOT NULL,
    start TIME NOT NULL,
    duration BIGINT NOT NULL,
    CONSTRAINT time_slots_pk PRIMARY KEY (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE time_slots;
-- +goose StatementEnd
