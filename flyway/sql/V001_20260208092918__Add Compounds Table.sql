-- Add your migration SQL here
CREATE TABLE compounds (
    inchikey  CHAR(27) PRIMARY KEY,
    name      TEXT NOT NULL,
    inchi     TEXT NOT NULL,
    smiles    TEXT NOT NULL,
    formula   TEXT NOT NULL
);

CREATE INDEX idx_compounds_name_lower   ON compounds (LOWER(name));
CREATE INDEX idx_compounds_formula      ON compounds (formula);
