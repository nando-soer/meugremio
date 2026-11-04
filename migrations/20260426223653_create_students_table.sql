-- +goose Up
CREATE TABLE students (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    email VARCHAR(254) NOT NULL UNIQUE,
    class TEXT NOT NULL,
    password TEXT NOT NULL
);

CREATE UNIQUE INDEX idx_students_email ON students (email);

-- +goose Down
DROP TABLE students;
