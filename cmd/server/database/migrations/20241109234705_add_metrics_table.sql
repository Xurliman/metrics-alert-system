-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS metrics
(
    id                  uuid primary key unique not null default uuid_generate_v4(),
    name                text             	    not null,
    metric_type         text             		not null,
    value               double precision            null,
    delta               bigint                      null,
    created_at          timestamp                        default now(),
    updated_at          timestamp,
    deleted_at          timestamp
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS metrics;
-- +goose StatementEnd