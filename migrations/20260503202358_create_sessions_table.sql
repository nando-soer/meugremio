-- +goose Up
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    student_id BIGSERIAL NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX ON sessions (student_id);

-- +goose Down
DROP TABLE sessions
