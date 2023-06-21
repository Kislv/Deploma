DROP TABLE IF EXISTS users          CASCADE;

CREATE TABLE users (
    chatId                                BIGINT NOT NULL PRIMARY KEY,
    age                                   BIGINT,
    isMale                                BOOLEAN,
    nationality                           VARCHAR(100)
);
