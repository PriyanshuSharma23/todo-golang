CREATE TABLE IF NOT EXISTS todos (
    id bigserial PRIMARY KEY,
    title text NOT NULL,
    is_complete boolean NOT NULL,
    due_on timestamp NOT NULL,
    body text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version int NOT NULL DEFAULT 1
);