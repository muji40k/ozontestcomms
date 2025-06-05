\c poster

drop schema if exists posts cascade;
create schema posts;

drop table if exists posts.posts;
create table posts.posts
(
    id uuid primary key,
    author_id uuid not null,
    commentable_id uuid not null,
    title text not null,
    content text not null,
    creation_date timestamptz not null
);

alter table posts.posts add
    constraint "fkey_post_author_id"
    foreign key (author_id)
    references users.users(id);

alter table posts.posts add
    constraint "fkey_post_commentable_id"
    foreign key (commentable_id)
    references commentables.commentables(id);

alter table posts.posts add
    constraint "post_title_length"
    check (char_length(title) <= 1000);

alter table posts.posts add
    constraint "post_content_length"
    check (char_length(content) <= 4000);

