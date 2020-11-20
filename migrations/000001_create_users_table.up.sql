CREATE TABLE IF NOT EXISTS status (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL
);

INSERT INTO status (name) VALUES
    ('Work @ office'),
    ('Work @ home'),
    ('Out of office'),
    ('Holidays'),
    ('Away');

CREATE TABLE IF NOT EXISTS users (
    id                  TEXT PRIMARY KEY,
    email               VARCHAR(255) NOT NULL UNIQUE,
    name                VARCHAR(32) NOT NULL,
    picture             VARCHAR(255) NULL,
    status_id           INT NULL REFERENCES status(id),
    oidc_refresh_token  TEXT NULL
);
