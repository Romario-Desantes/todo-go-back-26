CREATE TABLE IF NOT EXISTS public.tasks
(
    id              serial PRIMARY KEY,
    title           varchar(50) NOT NULL,
    description     varchar(100),
    date            timestamptz,
    status          varchar(50) NOT NULL,
    created_date    timestamptz NOT NULL,
    updated_date    timestamptz NOT NULL,
    deleted_date    timestamptz
);