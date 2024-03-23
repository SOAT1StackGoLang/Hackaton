CREATE TABLE timekeeping
(
    id             uuid        NOT NULL,
    user_id        text        NOT NULL,
    created_at     timestamptz NOT NULL,
    reference_date date        NOT NULL,
    updated_at     timestamptz,
    worked_minutes integer     NOT NULL,
    open           boolean     NOT NULL DEFAULT true,
    details        jsonb       NOT NULL,

    CONSTRAINT timekeeping_pk PRIMARY KEY (id),
    CONSTRAINT timekeeping_user_id_reference_date_key UNIQUE (user_id, reference_date)
);

CREATE INDEX idx_user_id_reference_date ON timekeeping (user_id, created_at);

