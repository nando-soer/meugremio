-- +goose Up
CREATE TABLE pending_registers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL,
    data JSONB NOT NULL,
    attempts SMALLINT NOT NULL,
    expires_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE pending_registers