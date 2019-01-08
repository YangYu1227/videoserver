create table app_upload_video
(
  user_id   varchar(100) not null
    primary key,
  book_id   varchar(32)  null,
  video_vid varchar(64)  null
);

create table app_user
(
  deviceid       varchar(100) not null
    primary key,
  unlocked_books varchar(100) null,
  nick_name      varchar(64)  null
);

create table barcode
(
  id          int auto_increment
    primary key,
  barcode_id  varchar(64)                         not null,
  bar_code    varchar(64)                         not null,
  is_use      int(2)    default 0                 null,
  create_time timestamp default CURRENT_TIMESTAMP null
);

create table book
(
  id            int auto_increment
    primary key,
  book_name     varchar(32)                         not null,
  book_id       varchar(32)                         not null,
  barcode_count int                                 not null,
  create_time   timestamp default CURRENT_TIMESTAMP not null,
  constraint book_book_name_uindex
    unique (book_name)
);

create table book_bar_code
(
  id          int auto_increment
    primary key,
  bar_code_id varchar(64)                         not null,
  book_id     varchar(64)                         not null,
  all_count   int                                 not null,
  use_count   int                                 not null,
  status      int(4)                              not null,
  create_time timestamp default CURRENT_TIMESTAMP null,
  update_time timestamp default CURRENT_TIMESTAMP null,
  constraint book_bar_code_bar_code_id_uindex
    unique (bar_code_id)
);

create table book_user
(
  id       int auto_increment
    primary key,
  name     varchar(64) not null,
  password varchar(64) not null
);

create table sessions
(
  session_id varchar(64) not null
    primary key,
  TTL        tinytext    null,
  login_name varchar(64) null
);

INSERT INTO `book_server`.`book_user` (`name`, `password`) VALUES ('admin', '123456')