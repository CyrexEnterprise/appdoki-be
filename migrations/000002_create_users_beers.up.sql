CREATE TABLE IF NOT EXISTS beer_transfers (
    giver_id TEXT NOT NULL REFERENCES users(id),
    taker_id TEXT NULL REFERENCES users(id),
    beers    INT CHECK(beers > 0),
    given_at TIMESTAMP NOT NULL DEFAULT now()
);
