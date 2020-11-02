-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS "course_study_year_junction" CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE course_study_year_junction (
    course_id UUID NOT NULL,
    study_year_id UUID NOT NULL,

    CONSTRAINT course_sy_pk PRIMARY KEY (course_id, study_year_id),
    CONSTRAINT FK_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
    CONSTRAINT FK_study_year FOREIGN KEY (study_year_id) REFERENCES study_years(id) ON DELETE CASCADE
);
-- +goose StatementEnd
