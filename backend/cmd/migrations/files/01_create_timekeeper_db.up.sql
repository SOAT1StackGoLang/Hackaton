CREATE TABLE timekeeping
(
    id         uuid        not null,
    user_id    uuid        not null,
    created_at timestamptz NOT NULL
);
