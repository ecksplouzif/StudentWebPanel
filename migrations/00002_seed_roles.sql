-- +goose Up
INSERT into roles(id, type) VALUES
    (1, 'class_representative'),
    (2, 'deputy_class_representative'),
    (3, 'union_representative'),
    (4, 'student')
;

-- +goose Down
DELETE FROM roles WHERE id IN (1,2,3,4);
