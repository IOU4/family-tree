drop table if exists person;

create table gender (id integer primary key, name text unique not null);
insert into gender values (1, 'FEMALE'), (2, 'MALE');

create table relation_type (id INTEGER primary key, name TEXT UNIQUE NOT NULL);
insert into relation_type values(1, 'parent'), (2, 'sibling');

create table person (
  id INTEGER primary key autoincrement,
  first_name varchar(50) not null,
  last_name varchar(50) not null,
  birthday DATE not null,
  gender INTEGER refrences gender not null,
  misc TEXT
);

create table if not exists relation (
  id INTEGER primary key autoincrement,
  src INTEGER refrences person not null,
  dest INTEGER refrences person not null,
  type INTEGER refrences relation_type not null
);

insert into person (first_name, last_name, birthday, gender) values ('imad', 'ouchaib', '2000-11-18', 1), ('marwa', 'ouchaib', '2012-01-03', 0);
insert into relation (src, dest, type) values(1, 2, 2);
