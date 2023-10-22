CREATE TABLE people_info(
    id uuid primary key,
    name varchar(50) CHECK (length(name) > 0),
    surname varchar(80) CHECK ( length(surname) > 0 ),
    patronymic varchar(120) CHECK ( length(patronymic) > 0 ),
    age int CHECK ( age > 0 ),
    sex varchar(10) CHECK ( sex IN ('male', 'female') )
);

CREATE TABLE person_nationality(
    id uuid primary key,
    nationality varchar(5) CHECK (length(nationality) > 0),
    probability float4,
    person_id uuid references people_info(id) ON DELETE CASCADE
);