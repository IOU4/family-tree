create table gender (id TEXT PRIMARY KEY, desc TEXT);
insert into gender values ('F', 'female'), ('M', 'male');

create table kinship (id TEXT PRIMARY KEY, desc TEXT);
insert into kinship values('P', 'parent'), ('S', 'sibling'), ('H', 'half, partner');

create table person (
  id INTEGER primary key autoincrement,
  first_name varchar(50) not null,
  last_name varchar(50) not null,
  birthday DATE not null,
  gender INTEGER REFERENCES gender not null,
  misc TEXT
);

create table if not exists relation (
  id INTEGER primary key autoincrement,
  p1 INTEGER REFERENCES person not null,
  p2 INTEGER REFERENCES person not null,
  kinship TEXT REFERENCES kinship not null
);

insert into person (id, first_name, last_name, birthday, gender) values 
  (1, 'lhoussain', 'ouchaib', '1970-10-10', 'M'),
  (2, 'hakima', 'nait-abbou', '1979-10-10', 'F'),
  (3, 'soukaina', 'ouchaib', '1998-10-01', 'F'),
  (4, 'mariam', 'ouchaib', '1999-11-24', 'F'),
  (5, 'jawad', 'ouchaib', '2000-11-18', 'M'),
  (6, 'imad', 'ouchaib', '2000-11-18', 'M'),
  (7, 'marwa', 'ouchaib', '2012-01-03', 'F'),
  (8, 'imran', 'ouchaib', '2020-01-03', 'M');
insert into relation (p1, p2, kinship) values(1, 2, 'H');
insert into relation (p1, p2, kinship) values(1, 3, 'P');
insert into relation (p1, p2, kinship) values(3, 4, 'S');
insert into relation (p1, p2, kinship) values(3, 5, 'S');
insert into relation (p1, p2, kinship) values(3, 6, 'S');
insert into relation (p1, p2, kinship) values(3, 7, 'S');
insert into relation (p1, p2, kinship) values(3, 8, 'S');
