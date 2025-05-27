-- +goose Up
-- SQL in sqlite dialect

-- User — участник Telegram-чата
CREATE TABLE users
(
    id       INTEGER PRIMARY KEY,           -- Уникальный идентификатор пользователя в Telegram
    nickname TEXT,                          -- никнейм пользователя, подпись которая будет показана в чате
    username TEXT    NOT NULL,              -- Имя пользователя в Telegram
    is_owner BOOLEAN NOT NULL DEFAULT FALSE -- Является ли пользователь суперадмином
);

-- Meeting — событие встречи
CREATE TABLE meetings
(
    id           INTEGER PRIMARY KEY,
    code         TEXT     NOT NULL UNIQUE,                    -- Уникальный код встречи
    status       TEXT     NOT NULL,                           -- Статус встречи: active, canceled, completed, draft
    published_at DATETIME,                                    -- Время публикации в чате
    closed_at    DATETIME,                                    -- Время закрытия/отмены
    updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Время последнего обновления сообщения через крон
    message      TEXT,                                        -- Текст объявления встречи
    max          INTEGER,                                     -- Лимит участников
    cost         INTEGER,                                     -- Стоимость участия
    type_pay     TEXT     NOT NULL,                           -- Способ оплаты участия: FREE, SPLIT, FIXED
    owner_id     INTEGER  NOT NULL,                           -- ID владельца встречи (User)
    FOREIGN KEY (owner_id) REFERENCES users (id)
);

-- Chat — Telegram-группа или супергруппа
CREATE TABLE chats
(
    id         INTEGER PRIMARY KEY,            -- Telegram Chat ID
    title      TEXT    NOT NULL,               -- Название чата
    is_meeting BOOLEAN NOT NULL DEFAULT TRUE,  -- Включить ли механизм создания встреч
    is_antibot BOOLEAN NOT NULL DEFAULT FALSE, -- Включить антибот-проверку для новых пользователей
    is_private BOOLEAN NOT NULL DEFAULT FALSE  -- Это приватный чат между пользователем и ботом
);

-- ChatMeeting — промежуточная таблица Chat ↔ Meeting
CREATE TABLE chat_meetings
(
    chat_id    INTEGER NOT NULL, -- ID чата
    meeting_id INTEGER NOT NULL, -- ID встречи
    message_id INTEGER NOT NULL, -- ID сообщения в чате для встречи
    PRIMARY KEY (chat_id, meeting_id),
    FOREIGN KEY (chat_id) REFERENCES chats (id),
    FOREIGN KEY (meeting_id) REFERENCES meetings (id)
);

-- UserMeeting — Дополнительные люди от одного участника
CREATE TABLE user_meetings
(
    user_id    INTEGER NOT NULL,  -- ID пользователя
    meeting_id INTEGER NOT NULL,  -- ID встречи
    count      INTEGER DEFAULT 0, -- Количество приведённых участников
    status     TEXT    NOT NULL,  -- Статус: CANCEL или YES
    PRIMARY KEY (user_id, meeting_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (meeting_id) REFERENCES meetings (id)
);


-- +goose Down
-- SQL in sqlite dialect

DROP TABLE IF EXISTS user_meetings;
DROP TABLE IF EXISTS chat_meetings;
DROP TABLE IF EXISTS chats;
DROP TABLE IF EXISTS meetings;
DROP TABLE IF EXISTS users;