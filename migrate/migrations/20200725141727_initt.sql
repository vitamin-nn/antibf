-- +goose Up
CREATE TABLE ip_white_list
(
    id serial primary key,
    ip_network inet,
    EXCLUDE USING GIST (ip_network inet_ops WITH &&)
);

CREATE TABLE ip_black_list
(
    id serial primary key,
    ip_network inet,
    EXCLUDE USING GIST (ip_network inet_ops WITH &&)
);

-- +goose Down
DROP TABLE ip_white_list;
DROP TABLE ip_black_list;
