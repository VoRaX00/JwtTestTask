-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);

CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES Users(id) ON DELETE CASCADE NOT NULL ,
    refresh_token_hash TEXT UNIQUE NOT NULL,
    ip VARCHAR(45) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE refresh_tokens;
DROP TABLE users;
-- +goose StatementEnd
