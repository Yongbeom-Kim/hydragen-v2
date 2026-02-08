-- Add your migration SQL here

CREATE TABLE mass_spectra (
    id               BIGSERIAL PRIMARY KEY,
    inchikey         CHAR(27) REFERENCES compounds(inchikey) NOT NULL,

    molecular_weight DOUBLE PRECISION NOT NULL,
    exact_mass       DOUBLE PRECISION,

    precursor_mz     DOUBLE PRECISION,
    precursor_type   TEXT,
    ion_mode         TEXT,
    collision_energy TEXT,

    spectrum_type    TEXT,
    instrument       TEXT,
    instrument_type  TEXT,

    splash           TEXT,
    db_number        TEXT NOT NULL,
    source           TEXT NOT NULL,
    comments         TEXT,

    raw_data         JSONB NOT NULL,

    UNIQUE (inchikey, db_number, source)
);

CREATE INDEX idx_spectra_source ON mass_spectra (source);
CREATE INDEX idx_mass_spectra_db_number ON mass_spectra (db_number);
CREATE INDEX idx_mass_spectra_molecular_weight ON mass_spectra (molecular_weight);