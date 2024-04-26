CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS restables (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    numberofseats integer,
    isreserved boolean,
    tablenumber integer,
    restaurantID UUID NOT NULL REFERENCES restaurants ON DELETE CASCADE
    );