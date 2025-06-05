\c postgres

drop database if exists poster;
create database poster;

\i /scripts/init_users.sql
\i /scripts/init_commentables.sql
\i /scripts/init_posts.sql
\i /scripts/init_comments.sql

