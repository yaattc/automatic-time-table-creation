-- +goose Up
-- +goose StatementBegin

ALTER TABLE courses ADD COLUMN edu_program TEXT NOT NULL DEFAULT 'bachelor';
-- noinspection SqlAddNotNullColumn
ALTER TABLE courses ADD COLUMN primary_lector_id UUID NOT NULL;
ALTER TABLE courses ADD CONSTRAINT FK_primary_lector FOREIGN KEY (primary_lector_id) REFERENCES teachers(id) ON DELETE NO ACTION;
-- noinspection SqlAddNotNullColumn
ALTER TABLE courses ADD COLUMN assistant_lector_id UUID;
ALTER TABLE courses ADD CONSTRAINT FK_assistant_lector FOREIGN KEY (assistant_lector_id) REFERENCES teachers(id) ON DELETE NO ACTION;

CREATE TABLE courses_teacher_assistants (
    course_id UUID NOT NULL,
    assistant_id UUID NOT NULL,
    CONSTRAINT FK_course_id FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
    CONSTRAINT FK_assistant_id FOREIGN KEY (assistant_id) REFERENCES teachers(id) ON DELETE CASCADE,
    CONSTRAINT cta_pk PRIMARY KEY (course_id, assistant_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE courses DROP COLUMN edu_program;
ALTER TABLE courses DROP COLUMN primary_lector_id;
ALTER TABLE courses DROP COLUMN assistant_lector_id;
DROP TABLE courses_teacher_assistants;
-- +goose StatementEnd
