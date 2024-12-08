-- +goose Up
CREATE TABLE customers
(
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT        NOT NULL,
    email      TEXT        NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE orders(
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    robot_model text NOT NULL,
    robot_version text NOT NULL,
    created_at timestamptz NOT NULL
);

CREATE TABLE robots(
    id BIGSERIAL PRIMARY KEY,
    model text NOT NULL,
    version text NOT NULL,
    created_at timestamptz NOT NULL
);

-- +goose Down
DROP TABLE robots;
DROP TABLE orders;
DROP TABLE customers;
