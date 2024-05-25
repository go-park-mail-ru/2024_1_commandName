CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE IF NOT EXISTS auth.person
(
    id            INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    username      TEXT UNIQUE CHECK (LENGTH(username) <= 20),
    email         TEXT CHECK (LENGTH(email) <= 30),
    name          TEXT CHECK (LENGTH(name) <= 30),
    surname       TEXT CHECK (LENGTH(surname) <= 20),
    about         TEXT CHECK (LENGTH(about) <= 50) DEFAULT '',
    password_hash TEXT        NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL             DEFAULT CURRENT_TIMESTAMP,
    lastseen_at   TIMESTAMPTZ NOT NULL             DEFAULT CURRENT_TIMESTAMP,
    avatar_path   TEXT                             DEFAULT '',
    password_salt TEXT        NOT NULL,
    language      TEXT                             DEFAULT 'ru'
);

CREATE SCHEMA IF NOT EXISTS chat;

CREATE TABLE IF NOT EXISTS chat.chat_type
(
    id   TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS chat.chat
(
    id          INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    type_id     TEXT REFERENCES chat.chat_type (id) NOT NULL DEFAULT '1',
    name        TEXT                                NOT NULL CHECK (length(name) <= 20),
    description TEXT CHECK (length(description) <= 70)       DEFAULT '',
    avatar_path TEXT                                         DEFAULT '',
    created_at  TIMESTAMPTZ                         NOT NULL DEFAULT current_TIMESTAMP,
    edited_at   TIMESTAMPTZ                         NOT NULL DEFAULT current_TIMESTAMP,
    creator_id  INT REFERENCES auth.person (id)
);

CREATE TABLE IF NOT EXISTS chat.message
(
    id          INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    user_id     INT REFERENCES auth.person (id),
    chat_id     INT REFERENCES chat.chat (id),
    message     TEXT CHECK (length(message) <= 1000) DEFAULT '',
    edited      BOOLEAN     NOT NULL                 DEFAULT false,
    edited_at   TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL                 DEFAULT current_TIMESTAMP,
    file_exists BOOLEAN                              DEFAULT false
);

CREATE TABLE IF NOT EXISTS chat.file
(
    id         INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    message_id INT REFERENCES chat.message (id),
    type       TEXT default 'file',
    file_path  TEXT DEFAULT '',
    originalName TEXT DEFAULT ''
);

CREATE TABLE IF NOT EXISTS chat.sticker
(
    id          INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    description TEXT,
    type        TEXT,
    file_path   TEXT DEFAULT ''
);

CREATE TABLE IF NOT EXISTS chat.chat_user
(
    chat_id             INT REFERENCES chat.chat (id),
    user_id             INT REFERENCES auth.person (id),
    lastseen_message_id INT
);

CREATE TABLE IF NOT EXISTS chat.contact_state
(
    id   INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS chat.contacts
(
    id       INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    user1_id INT REFERENCES auth.person (id),
    user2_id INT REFERENCES auth.person (id),
    state_id INT REFERENCES chat.contact_state (id)
);

CREATE TABLE IF NOT EXISTS auth.session
(
    id        INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    sessionid TEXT NOT NULL,
    userid    INT REFERENCES auth.person (id)
);