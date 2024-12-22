-- name: GetUser :one
select *
from users
where users.id = ?
;

-- name: ListUsers :many
select *
from users
order by time_spent_in_minutes desc
;

-- name: UpsertUser :exec
INSERT INTO users (id, name)
VALUES (?, ?)
ON CONFLICT (id) DO UPDATE SET
    name = excluded.name;


-- name: DeleteUsers :exec
delete from users
;


-- name: InsertUserJoin :exec
INSERT INTO user_joins (
    user_id,
    guild_id,
    channel_id,
    joined_at
) VALUES (?, ?, ?, CURRENT_TIMESTAMP);


-- name: UpdateUserLeave :exec
UPDATE user_joins 
SET left_at = CURRENT_TIMESTAMP
WHERE user_id = ? 
AND left_at IS NULL;


-- name: GetAllUsersWeeklyTimeSpent :many
with
    week_start as (
        -- Get the start of current week (Monday 00:00:00)
        select
            datetime(
                'now',
                'start of day',
                '-' || case
                    strftime('%w', 'now')
                    when '0'
                    then '6'  -- If Sunday, go back 6 days
                    else cast(strftime('%w', 'now') - 1 as text)  -- Otherwise, back to Monday
                end
                || ' days'
            ) as start_date
    )
select
    u.name,
    count(distinct uj.id) as joins_this_week,
    sum(
        cast(
            (
                julianday(coalesce(uj.left_at, current_timestamp))
                - julianday(uj.joined_at)
            )
            * 1440 as integer
        )
    ) as minutes_this_week
from users u
join user_joins uj on u.id = uj.user_id
cross join week_start
where uj.joined_at >= week_start.start_date and uj.joined_at <= current_timestamp
group by u.id
order by minutes_this_week desc
limit 10
;

-- name: GetAllUsersTodayTimeSpent :many
select
    u.name,
    count(distinct uj.id) as joins_today,
    sum(
        cast(
            (
                julianday(coalesce(uj.left_at, current_timestamp))
                - julianday(uj.joined_at)
            )
            * 1440 as integer
        )
    ) as minutes_today
from users u
join user_joins uj on u.id = uj.user_id
where uj.joined_at >= date('now', 'start of day') and uj.joined_at <= current_timestamp
group by u.id
order by minutes_today desc
limit 10
;

-- name: GetUserTodayTimeSpent :one
select
    u.name,
    count(distinct uj.id) as joins_today,
    sum(
        cast(
            (
                julianday(coalesce(uj.left_at, current_timestamp))
                - julianday(uj.joined_at)
            )
            * 1440 as integer
        )
    ) as minutes_today
from users u
join user_joins uj on u.id = uj.user_id
where
    u.id = ?
    and uj.joined_at >= date('now', 'start of day')
    and uj.joined_at <= current_timestamp
group by u.id
;

-- name: GetUserWeeklyTimeSpent :one
with
    week_start as (
        -- Get the start of current week (Monday 00:00:00)
        select
            datetime(
                'now',
                'start of day',
                '-' || case
                    strftime('%w', 'now')
                    when '0'
                    then '6'  -- If Sunday, go back 6 days
                    else cast(strftime('%w', 'now') - 1 as text)  -- Otherwise, back to Monday
                end
                || ' days'
            ) as start_date
    )
select
    u.name,
    count(distinct uj.id) as joins_this_week,
    sum(
        cast(
            (
                julianday(coalesce(uj.left_at, current_timestamp))
                - julianday(uj.joined_at)
            )
            * 1440 as integer
        )
    ) as minutes_this_week
from users u
join user_joins uj on u.id = uj.user_id
cross join week_start
where
    uj.joined_at >= week_start.start_date
    and uj.joined_at <= current_timestamp
    and u.id
    =  -- Added filter for specific user
    ?
group by u.id
;

-- name: GetUserTotalTimeSpent :one
select
    u.name,
    count(distinct uj.id) as total_joins,
    sum(
        cast(
            (
                julianday(coalesce(uj.left_at, current_timestamp))
                - julianday(uj.joined_at)
            )
            * 1440 as integer
        )
    ) as total_minutes,
    max(uj.joined_at) as last_join
from users u
join user_joins uj on u.id = uj.user_id
where u.id = ?
group by u.id
;


-- name: GetAllTimeStats :many
select
    u.name,
    count(distinct uj.id) as total_joins,
    sum(
        cast(
            (
                julianday(coalesce(uj.left_at, current_timestamp))
                - julianday(uj.joined_at)
            )
            * 1440 as integer
        )
    ) as total_minutes
from users u
join user_joins uj on u.id = uj.user_id
group by u.id, u.name
order by total_minutes desc
limit 10
;

-- name: UpdateActiveSessions :exec
UPDATE user_joins
SET left_at = CURRENT_TIMESTAMP,  -- End the current session
    joined_at = CURRENT_TIMESTAMP  -- Start a new session for the same duration
WHERE user_id = ? 
AND left_at IS NULL;

