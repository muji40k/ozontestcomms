\c poster

drop schema if exists commentables cascade;
create schema commentables;

drop table if exists commentables.commentables;
create table commentables.commentables
(
    id uuid primary key,
    comments_allowed boolean
);

