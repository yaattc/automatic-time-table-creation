-- +goose Up
-- +goose StatementBegin
ALTER TABLE courses ADD COLUMN study_year_id UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000001';
ALTER TABLE courses ADD CONSTRAINT FK_study_year FOREIGN KEY (study_year_id) REFERENCES study_years(id) ON DELETE CASCADE ;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE courses DROP COLUMN study_year_id;
-- +goose StatementEnd
