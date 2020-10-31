-- +goose Up
-- +goose StatementBegin
CREATE TABLE courses (
    id UUID NOT NULL,
    name TEXT NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE teachers (
    id UUID NOT NULL,
    name TEXT NOT NULL,
    surname TEXT NOT NULL,

    email TEXT NOT NULL,
    degree TEXT NOT NULL,
    about TEXT NOT NULL,
    PRIMARY KEY (id)
);


CREATE TABLE teacher_preferences (
    teacher_id UUID NOT NULL UNIQUE,
    locations JSONB,
    CONSTRAINT FK_teacher FOREIGN KEY (teacher_id) REFERENCES teachers(id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION create_prefs_on_teacher_insert()
    RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO teacher_preferences("teacher_id", "locations") VALUES (NEW.id, NULL) ON CONFLICT DO NOTHING;
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER teachers_create_preferences
    AFTER INSERT ON teachers
    FOR EACH ROW EXECUTE PROCEDURE create_prefs_on_teacher_insert();

CREATE TABLE teacher_preferences_staff (
    teacher_id UUID NOT NULL,
    staff_id UUID NOT NULL,
    CONSTRAINT ts_pk PRIMARY KEY (teacher_id, staff_id),
    CONSTRAINT teacher_staff UNIQUE (teacher_id, staff_id),
    CONSTRAINT FK_teacher FOREIGN KEY (teacher_id) REFERENCES teachers(id) ON DELETE CASCADE ,
    CONSTRAINT FK_staff FOREIGN KEY (staff_id) REFERENCES teachers(id) ON DELETE CASCADE
);

CREATE TABLE teacher_preferences_time_slots (
    teacher_id UUID NOT NULL,
    CONSTRAINT FK_preference FOREIGN KEY (teacher_id) REFERENCES teachers(id) ON DELETE CASCADE ,

    weekday INTEGER NOT NULL,
    start TIME NOT NULL,
    duration BIGINT NOT NULL,
    location TEXT NOT NULL
);

CREATE TABLE study_years (
    id UUID NOT NULL,
    name TEXT NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE groups (
    id UUID NOT NULL,
    study_year_id UUID NOT NULL,
    name TEXT NOT NULL,

    CONSTRAINT FK_study_year FOREIGN KEY (study_year_id) REFERENCES study_years(id) ON DELETE CASCADE ,
    PRIMARY KEY (id)
);

CREATE TABLE classes (
    id UUID NOT NULL,
    course_id UUID NOT NULL,
    group_id UUID NOT NULL,
    teacher_id UUID NOT NULL,

    title TEXT NOT NULL,
    location TEXT NOT NULL DEFAULT '',

    start_time TIMESTAMP NOT NULL,
    duration BIGINT NOT NULL DEFAULT 5400000000000, -- 1.5 hours

    repeats INTEGER NOT NULL,

    CONSTRAINT FK_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
    CONSTRAINT FK_group FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    CONSTRAINT FK_teacher FOREIGN KEY (teacher_id) REFERENCES teachers(id) ON DELETE CASCADE,
    PRIMARY KEY (id)
);

CREATE TABLE course_study_year_junction (
    course_id UUID NOT NULL,
    study_year_id UUID NOT NULL,

    CONSTRAINT course_sy_pk PRIMARY KEY (course_id, study_year_id),
    CONSTRAINT FK_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
    CONSTRAINT FK_study_year FOREIGN KEY (study_year_id) REFERENCES study_years(id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "teachers" CASCADE;
DROP TABLE IF EXISTS "teacher_preferences" CASCADE;
DROP TABLE IF EXISTS "teacher_preferences_time_slots" CASCADE;
DROP TABLE IF EXISTS "teacher_preferences_staff" CASCADE;

DROP TABLE IF EXISTS "study_years" CASCADE;
DROP TABLE IF EXISTS "courses" CASCADE;

DROP TABLE IF EXISTS "groups" CASCADE;
DROP TABLE IF EXISTS "classes" CASCADE;
DROP TABLE IF EXISTS "course_study_year_junction" CASCADE;
-- +goose StatementEnd
