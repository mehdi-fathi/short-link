create table links
(
    id        serial,
    link      char(32),
    short_key varchar(20)
);

create table users
(
    id        integer,
    name      char(32),
    username  char(32)
);
