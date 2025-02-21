create table user
(
    id int auto_increment
    primary key,
    email varchar(40) not null,
    password varchar(40) not null,
    name varchar(40) not null,
    gender varchar(20) not null,
    age int(4) not null
)

create table match
(
    matched boolean,
    matchID int auto_increment
)