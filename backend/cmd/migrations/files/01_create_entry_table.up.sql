CREATE TABLE entries
(
    id         uuid       not null,
    user_id    uuid       not null,
    created_at timestamptz NOT NULL
);

create index idx_user_created on entries (user_id, created_at);
