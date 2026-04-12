create table gender (id TEXT PRIMARY KEY, desc TEXT);
insert into gender values ('F', 'female'), ('M', 'male');

create table person (
  id INTEGER primary key autoincrement,
  first_name varchar(50) not null,
  last_name varchar(50) not null,
  birth DATE not null,
  death DATE,
  gender INTEGER REFERENCES gender not null,
  misc TEXT
);

create table relation (
  person INTEGER primary key REFERENCES person not null,
  father INTEGER REFERENCES person not null,
  mother INTEGER REFERENCES person not null
);

insert into person (id, first_name, last_name, birth, gender) values 
  (1, 'lhoussain', 'ouchaib', '1970-10-10', 'M'),
  (2, 'hakima', 'nait-abbou', '1979-10-10', 'F'),
  (3, 'soukaina', 'ouchaib', '1998-10-01', 'F'),
  (4, 'mariam', 'ouchaib', '1999-10-01', 'F'),
  (5, 'jawad', 'ouchaib', '2000-11-18', 'M'),
  (6, 'imad', 'ouchaib', '2000-11-18', 'M');
insert into relation (person, father, mother) values(3, 1, 2);
insert into relation (person, father, mother) values(4, 1, 2);
insert into relation (person, father, mother) values(5, 1, 2);
insert into relation (person, father, mother) values(6, 1, 2);
