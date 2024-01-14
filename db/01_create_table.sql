create table posts (
    id uuid primary key default gen_random_uuid(),
    title text not null,
    content text not null,
    created_at timestamp with time zone not null,
    updated_at timestamp with time zone not null
);
