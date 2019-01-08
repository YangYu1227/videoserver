create table app_upload_video
(
	user_id varchar(100) null,
	book_id varchar(32) null,
	video_vid varchar(64) null,
	constraint app_upload_video_pk
		primary key (user_id)
);

create table app_user
(
	deviceid varchar(100) not null,
	unlocked_books varchar(100) null,
	nick_name varchar(64) null,
	constraint app_user_pk
		primary key (deviceid)
);

create table barcode
(
	id int auto_increment,
	barcode_id varchar(64) not null,
	bar_code varchar(64) not null,
	is_use int(2) default 0 null,
	create_time timestamp default current_timestamp null,
	constraint barcode_pk
		primary key (id)
);

create table book
(
	id int auto_increment,
	book_name varchar(32) not null,
	book_id varchar(32) null,
	barcode_count int null,
	create_time timestamp default current_timestamp null,
	constraint book_pk
		primary key (id)
);

create unique index book_book_name_uindex
	on book (book_name);


create table book_bar_code
(
	id int auto_increment,
	bar_code_id varchar(64) not null,
	book_id varchar(64) not null,
	all_count int not null,
	use_count int not null,
	status int(4) not null,
	create_time timestamp default current_timestamp null,
	update_time timestamp default current_timestamp null,
	constraint book_bar_code_pk
		primary key (id)
);

create unique index book_bar_code_bar_code_id_uindex
	on book_bar_code (bar_code_id);


create table book_user
(
	id int auto_increment,
	name varchar(64) not null,
	password varchar(64) not null,
	constraint book_user_pk
		primary key (id)
);

create table sessions
(
	session_id varchar(64) not null,
	TTL tinytext null,
	login_name varchar(64) null,
	constraint sessions_pk
		primary key (session_id)
);

INSERT INTO `book_server`.`book_user` (`name`, `password`) VALUES ('admin', '123456')