-- +goose Up
-- +goose StatementBegin
CREATE TABLE coin_base (
                           title VARCHAR(10) NOT NULL,
                           rate DECIMAL(15,2) NOT NULL,
                           date DATE NOT NULL
);
-- +goose StatementEnd
