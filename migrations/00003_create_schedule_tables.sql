-- +goose Up
CREATE TABLE subjects(
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL UNIQUE,
    check_link text 
);
CREATE TABLE class_types(
    id int PRIMARY KEY,
    type text NOT NULL UNIQUE
);
CREATE TABLE schedule(
    group_id int,
    FOREIGN KEY (group_id) REFERENCES groups(id),
    subject_id int NOT NULL,
    FOREIGN KEY (subject_id) REFERENCES subjects(id),
    type_id int NOT NULL,
    FOREIGN KEY (type_id) REFERENCES class_types(id),
    start_time time NOT NULL,
    end_time time NOT NULL,    
    day date NOT NULL,
    PRIMARY KEY(group_id,day,start_time)
);
CREATE TABLE links(
    subject_id int,
    FOREIGN KEY (subject_id) REFERENCES subjects(id),
    type_id int,
    FOREIGN KEY (type_id) REFERENCES class_types(id),
    link text NOT NULL,
    group_id int,
    FOREIGN KEY (group_id) REFERENCES groups(id),
    PRIMARY KEY(group_id,subject_id,type_id)
);

-- +goose Down
DROP TABLE schedule;
DROP TABLE links;
DROP TABLE subjects;
DROP TABLE class_types;
