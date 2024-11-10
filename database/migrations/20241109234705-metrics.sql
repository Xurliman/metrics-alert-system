-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
DROP TYPE IF EXISTS metric_type;
CREATE TYPE metric_type AS ENUM ('gauge', 'counter');
CREATE TABLE IF NOT EXISTS metrics
(
    id                  uuid primary key unique not null default Uuid_generate_v4(),
    name                text             unique not null,
    metric_type         metric_type             not null,
    value               double precision            null,
    delta               int                         null,
    created_at          timestamp                        default now(),
    updated_at          timestamp,
    deleted_at          timestamp
);
-- +migrate Down
DROP TABLE IF EXISTS metrics;