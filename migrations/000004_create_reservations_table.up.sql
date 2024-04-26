CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS reservations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    userid          TEXT NOT NULL REFERENCES users ON DELETE CASCADE,
    tableid         UUID NOT NULL REFERENCES restables ON DELETE CASCADE,
    reservationtime TEXT
    );