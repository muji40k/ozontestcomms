\c poster

drop schema if exists comments cascade;
create schema comments;

drop table if exists comments.comments;
create table comments.comments
(
    id uuid primary key,
    author_id uuid not null,
    commentable_id uuid not null,
    target_id uuid not null,
    content text not null,
    creation_date timestamptz not null
);

alter table comments.comments add
    constraint "fkey_comment_author_id"
    foreign key (author_id)
    references users.users(id);

alter table comments.comments add
    constraint "fkey_comment_commentable_id"
    foreign key (commentable_id)
    references commentables.commentables(id);

alter table comments.comments add
    constraint "fkey_comment_target_id"
    foreign key (target_id)
    references commentables.commentables(id);

alter table comments.comments add
    constraint "comment_content_length"
    check (char_length(content) <= 2000);

