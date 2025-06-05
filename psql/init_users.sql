\c poster

drop schema if exists users cascade;
create schema users;

drop table if exists users.users;
create table users.users
(
    id uuid primary key,
    email text not null unique,
    password text not null
);

