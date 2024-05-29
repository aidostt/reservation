CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS reservations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    userid          TEXT NOT NULL,
    tableid         UUID NOT NULL REFERENCES restables ON DELETE CASCADE,
    reservationtime TEXT,
    reservationdate DATE,
    confirmed       BOOLEAN
    );