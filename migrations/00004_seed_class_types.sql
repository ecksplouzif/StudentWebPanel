-- +goose Up
INSERT into class_types(id, type) VALUES
    (1, 'Лк'),
    (2, 'Пз'),
    (3, 'Лб')
;

-- +goose Down
DELETE FROM class_types WHERE id IN (1,2,3);
