-- +goose Up
set schema 'trading';

create table users
(
  id uuid primary key,
  username varchar not null,
  password varchar,
  type varchar not null,
  details jsonb,
  activated_at timestamp,
  locked_at timestamp,
  created_at timestamp not null,
  updated_at timestamp not null,
  deleted_at timestamp null
);

create index idx_users_un on users(username);

create table sessions
(
  id uuid primary key,
  user_id uuid not null,
  username varchar not null,
  login_at timestamp not null,
  last_activity_at timestamp,
  details jsonb,
  logout_at timestamp
);

create index idx_sess_user_id on sessions(user_id);
create index idx_sess_username on sessions(username);

insert into users (id, username, password, type, details, activated_at, created_at, updated_at)
values('e722451f-fe6e-3a23-90c0-41d12a29ab32', 'admin', '$2a$10$m9FO4sFeDAwy6DvwisYwze5W3I83OktbzJNyTB3A/Po9fEE1vtL62', 'admin', '{"groups": ["sysadmin"]}', '2021-01-01 00:00:00.000', '2021-01-01 00:00:00.000', '2021-01-01 00:00:00.000');

-- +goose Down
set schema 'trading';

drop table users;
drop table sessions;
