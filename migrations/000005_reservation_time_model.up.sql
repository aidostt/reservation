CREATE EXTENSION IF NOT EXISTS btree_gist;

-- Replace the free-text time-of-day + date model with an explicit interval.
-- A reservation occupies its table for [start_at, ends_at); ends_at is derived
-- by the application from the restaurant's turn-time policy, so the duration is
-- policy-controlled rather than client-controlled.
ALTER TABLE reservations
    ADD COLUMN start_at   timestamptz NOT NULL,
    ADD COLUMN ends_at    timestamptz NOT NULL,
    ADD COLUMN party_size integer     NOT NULL DEFAULT 1;

-- The occupied interval; source of truth for availability and overlap.
ALTER TABLE reservations
    ADD COLUMN period tstzrange
        GENERATED ALWAYS AS (tstzrange(start_at, ends_at)) STORED;

-- Core domain invariant: a table cannot be booked for two overlapping
-- intervals. Enforced by the database, so it holds under concurrency and
-- across multiple service instances.
ALTER TABLE reservations
    ADD CONSTRAINT reservations_no_overlap
        EXCLUDE USING gist (tableid WITH =, period WITH &&);

ALTER TABLE reservations
    DROP COLUMN reservationtime,
    DROP COLUMN reservationdate;
