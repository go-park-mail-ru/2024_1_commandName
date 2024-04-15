-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE IF NOT EXISTS auth.person
(
    id                INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    username          TEXT UNIQUE CHECK (LENGTH(username) <= 20),
    email             TEXT NOT NULL CHECK (LENGTH(email) <= 30),
    name              TEXT NOT NULL CHECK (LENGTH(name) <= 30),
    surname           TEXT NOT NULL CHECK (LENGTH(surname) <= 20),
    about             TEXT CHECK (LENGTH(about) <= 50) DEFAULT '',
    password_hash     TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    lastseen_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    avatar            TEXT DEFAULT '',
    password_salt     TEXT NOT NULL
);

CREATE SCHEMA IF NOT EXISTS chat;

CREATE TABLE IF NOT EXISTS chat.chat
(
    id          INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    type        varchar(1) NOT NULL DEFAULT 1,
    name        TEXT NOT NULL CHECK (length(name) <= 20),
    description TEXT CHECK (length(description) <= 70) DEFAULT '',
    avatar_path TEXT DEFAULT '',
    edited_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    creator_id  INT REFERENCES auth.person (id)
);

CREATE TABLE IF NOT EXISTS chat.chat_user
(
    chat_id INT REFERENCES chat.chat (id),
    user_id INT REFERENCES auth.person (id)
);

CREATE TABLE IF NOT EXISTS chat.message
(
    id              INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    user_id         INT REFERENCES auth.person (id),
    chat_id         INT REFERENCES chat.chat (id),
    message         TEXT CHECK (length(message) <= 1000) DEFAULT '',
    edited          BOOLEAN NOT NULL DEFAULT false,
    create_datetime TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS chat.contacts
(
    id              INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    user1_id INT REFERENCES auth.person (id),
    user2_id INT REFERENCES auth.person (id),
    state INT NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS auth.session
(
    id        INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    sessionid TEXT NOT NULL,
    userid    INT REFERENCES auth.person (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
