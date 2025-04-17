begin;
CREATE TABLE IF NOT EXISTS coin_base (
                           title VARCHAR(50) NOT NULL,
                           rate REAL NOT NULL,
                           date TIMESTAMP NOT NULL DEFAULT NOW()
);
end;
