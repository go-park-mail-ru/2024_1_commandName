CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE IF NOT EXISTS auth.person
(
    id                INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    username          TEXT CHECK (length(username) <= 20),
    email             TEXT CHECK (length(email) <= 30),
    name              TEXT CHECK (length(name) <= 30),
    surname           TEXT CHECK (length(surname) <= 20),
    aboat             TEXT CHECK (length(aboat) <= 50),
    password_hash     TEXT,
    create_date       TIMESTAMP,
    lastseen_datetime TIMESTAMP,
    avatar            TEXT,
    password_salt TEXT
);

CREATE SCHEMA IF NOT EXISTS chat;

CREATE TABLE IF NOT EXISTS chat.chat
(
    id          INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    type        varchar(1),
    name        TEXT CHECK (length(name) <= 20),
    description TEXT CHECK (length(description) <= 70),
    avatar_path TEXT,
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
    message         TEXT CHECK (length(message) <= 1000),
    edited          BOOLEAN,
    create_datetime TIMESTAMP
);

CREATE TABLE IF NOT EXISTS auth.session
(
    id        INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    sessionid TEXT,
    userid    INT REFERENCES auth.person (id)
    );