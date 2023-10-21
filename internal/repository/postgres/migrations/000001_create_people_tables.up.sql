CREATE TABLE people_info(
    id uuid primary key,
    name varchar(50),
    surname varchar(80),
    patronymic varchar(120),
    age int,
    sex varchar(10)
);

CREATE TABLE person_nationality(
    id uuid primary key,
    nationality varchar(5),
    probability float4,
    person_id uuid references people_info(id)
);