-- name: CreateMeeting :one
INSERT INTO meetings
(max, cost, message, owner_id, type_pay, status, code)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, max, cost, message, owner_id, type_pay, status, code;

-- name: GetMeeting :one
SELECT id,
       max,
       cost,
       message,
       owner_id,
       type_pay,
       status,
       code
FROM meetings
WHERE id = @meeting_id;

-- name: GetMeetingByCode :one
SELECT id,
       max,
       cost,
       message,
       owner_id,
       type_pay,
       status,
       code
FROM meetings
WHERE code = @code;

-- name: UpdateMeetingStatus :exec
UPDATE meetings
SET status = @status
WHERE code = @code;

-- name: UpdateMeetingUpdate :exec
UPDATE meetings
SET updated_at = @updated_at
WHERE id = @where_meeting_id;

-- name: GetMeetingsWithStatus :many
SELECT code
FROM meetings
WHERE status = @status;

-- name: GetUser :one
SELECT id, username, is_owner,nickname
FROM users
WHERE id = @user_id;

-- name: GetUsers :many
SELECT id, username, is_owner,nickname
FROM users
ORDER BY id;

-- name: CreateUser :one
INSERT INTO users (id, username, is_owner, nickname)
VALUES (@id, @username, @is_owner, @nickname)
RETURNING id, username, is_owner,nickname;

-- name: UpdateUsername :exec
UPDATE users
SET username =  @username
WHERE id = @user_id;


-- name: GetChat :one
SELECT id, title, is_meeting, is_antibot
FROM chats
WHERE id = @chat_id;

-- name: CreateChat :one
INSERT INTO chats (id, title, is_meeting, is_antibot)
VALUES (@id, @title, @is_meeting, @is_antibot)
RETURNING id, title, is_meeting, is_antibot;


-- name: GetUserMeeting :one
SELECT user_id, meeting_id, status, count
FROM user_meetings
WHERE user_id = @user_id
  AND meeting_id = @meeting_id;

-- name: CreateUserMeeting :one
INSERT INTO user_meetings (user_id, meeting_id, status, count)
VALUES (@user_id, @meeting_id, @status, @count)
RETURNING user_id, meeting_id, status, count;

-- name: UpdateUserMeetingStatus :exec
UPDATE user_meetings
SET status = @status
WHERE user_id = @user_id
  AND meeting_id = @meeting_id;

-- name: UpdateUserMeetingCount :exec
UPDATE user_meetings
SET count = @count
WHERE user_id = @user_id
  AND meeting_id = @meeting_id;

-- name: GetUsersMeetings :many
SELECT um.user_id,
       um.meeting_id,
       um.status,
       um.count,
       u.username,
       u.nickname,
       u.is_owner
FROM user_meetings um
         JOIN users u ON u.id = um.user_id
WHERE um.meeting_id = @meeting_id
ORDER BY um.user_id;


-- name: GetChatMeeting :one
SELECT chat_id, meeting_id, message_id
FROM chat_meetings
WHERE chat_id = @chat_id
  AND meeting_id = @meeting_id;


-- name: GetChatMeetingAllChatWithMeeting :many
SELECT chat_id, meeting_id, message_id
FROM chat_meetings
WHERE meeting_id = @meeting_id;

-- name: CreateChatMeeting :one
INSERT INTO chat_meetings (chat_id, meeting_id, message_id)
VALUES (@chat_id, @meeting_id, @message_id)
RETURNING chat_id, meeting_id, message_id;


-- name: UpdateChatMeeting :one
UPDATE chat_meetings
SET message_id=COALESCE(sqlc.arg(message_id), message_id)
WHERE meeting_id = @where_meeting_id
  and chat_id = @where_chat_id
RETURNING *;

-- name: GetMeetingsForUpdateTime :many
SELECT m.id, m.code,
       m.status, m.published_at,
       m.closed_at, m.updated_at,
       m.message, m.max, m.cost,
       m.type_pay, m.owner_id, cm.chat_id,
       cm.meeting_id, cm.message_id,
       c.is_private,
       u.is_owner
FROM meetings m
join chat_meetings cm on m.id = cm.meeting_id
join users u on u.id = m.owner_id
join chats c on c.id = cm.chat_id
WHERE closed_at IS NULL
  AND updated_at <= datetime('now', '-2 day');