-- +goose Up
-- +goose StatementBegin
ALTER TABLE teacher_preferences_time_slots DROP COLUMN weekday;
ALTER TABLE teacher_preferences_time_slots DROP COLUMN start;
ALTER TABLE teacher_preferences_time_slots DROP COLUMN duration;
ALTER TABLE teacher_preferences_time_slots DROP COLUMN location;
ALTER TABLE teacher_preferences_time_slots ADD COLUMN slot_id UUID NOT NULL;
ALTER TABLE teacher_preferences_time_slots ADD CONSTRAINT FK_slot FOREIGN KEY (slot_id) REFERENCES time_slots(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE teacher_preferences_time_slots ADD COLUMN weekday INTEGER NOT NULL;
ALTER TABLE teacher_preferences_time_slots ADD COLUMN start TIME NOT NULL;
ALTER TABLE teacher_preferences_time_slots ADD COLUMN duration BIGINT NOT NULL;
ALTER TABLE teacher_preferences_time_slots ADD COLUMN location TEXT NOT NULL;
ALTER TABLE teacher_preferences_time_slots DROP COLUMN slot_id;
-- +goose StatementEnd
