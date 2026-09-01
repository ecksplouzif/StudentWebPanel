-- +goose Up
CREATE TABLE groups(
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY ,
    name text NOT NULL UNIQUE
);
CREATE TABLE roles(
    id int PRIMARY KEY,
    type text NOT NULL UNIQUE
);
CREATE TABLE users(
    sub uuid PRIMARY KEY,
    name text NOT NULL,
    group_id int,
    FOREIGN KEY (group_id) REFERENCES groups(id),
    role_id int,
    FOREIGN KEY (role_id) REFERENCES roles(id),
    is_admin bool NOT NULL
);

-- +goose Down
DROP TABLE users;
DROP TABLE groups;
DROP TABLE roles;