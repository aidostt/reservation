ALTER TABLE reservations DROP CONSTRAINT IF EXISTS reservations_no_overlap;

ALTER TABLE reservations
    ADD COLUMN reservationtime TEXT,
    ADD COLUMN reservationdate DATE;

ALTER TABLE reservations
    DROP COLUMN period,
    DROP COLUMN party_size,
    DROP COLUMN ends_at,
    DROP COLUMN start_at;
