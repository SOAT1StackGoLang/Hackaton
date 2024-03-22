CREATE TABLE timekeeping
(
    id             uuid        not null,
    user_id        text        not null,
    created_at     timestamptz NOT NULL,
    updated_at     timestamptz,
    worked_minutes integer     NOT NULL,
    open           boolean     NOT NULL default true,
    details        jsonb       NOT NULL,

    constraint timekeeping_pk
        PRIMARY KEY (id)
);

create index idx_user_id_reference_date on timekeeping (user_id, created_at);


