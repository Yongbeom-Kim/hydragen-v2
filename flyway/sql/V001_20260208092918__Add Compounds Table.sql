-- Add your migration SQL here
CREATE TABLE compounds (
    inchikey  CHAR(27) PRIMARY KEY,
    name      TEXT,
    inchi     TEXT,
    smiles    TEXT,
    formula   TEXT
);

CREATE INDEX idx_compounds_name_lower   ON compounds (LOWER(name));
CREATE INDEX idx_compounds_formula      ON compounds (formula);
