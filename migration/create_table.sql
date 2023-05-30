CREATE TABLE IF NOT EXISTS snippets (
id SERIAL PRIMARY KEY,
title VARCHAR(100) NOT NULL,
content TEXT NOT NULL,
created TIMESTAMP NOT NULL,
expires TIMESTAMP NOT NULL
);

CREATE INDEX idx_snippets_created ON snippets(created);

CREATE TABLE IF NOT EXISTS sessions (
    token CHAR(43) PRIMARY KEY,
    data BYTEA NOT NULL,
    expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);
