-- +goose Up
CREATE TABLE password_recovery (
    id SERIAL PRIMARY KEY,
    student_id BIGSERIAL NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX ON password_recovery (token_hash);

-- +goose Down
DROP TABLE password_recovery