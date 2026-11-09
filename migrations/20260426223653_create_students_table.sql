-- +goose Up
CREATE TABLE classes (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL
);

CREATE UNIQUE INDEX idx_classes_code ON classes (code);

INSERT INTO classes (code, name) VALUES 
    ('3-INFO-MANHA', '3° Ano de Informática (Manhã)'),
    ('2-INFO-MANHA', '2° Ano de Informática (Manhã)'),
    ('1-INFO-MANHA', '1° Ano de Informática (Manhã)'),
    ('3-ADM-NOITE', '3° Ano de Administração (Noite)'),
    ('2-ADM-NOITE', '2° Ano de Administração (Noite)'),
    ('1-ADM-NOITE', '1° Ano de Administração (Noite)'),
    ('1-MKT-NOITE', '1° Ano de Marketing (Tarde)');

CREATE TABLE students (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(254) NOT NULL UNIQUE,
    class_code VARCHAR(50) REFERENCES classes(code) NOT NULL,
    description TEXT,
    banner INT NOT NULL,
    avatar JSONB NOT NULL,
    password TEXT NOT NULL
);

CREATE UNIQUE INDEX idx_students_email ON students (email);

-- +goose Down
DROP TABLE students;
