"""Database facade for dataset loaders. Connection pool and all DB access go through this module."""

import os
from contextlib import AbstractContextManager
from typing import Any, Optional

import psycopg
from psycopg import Cursor
from psycopg_pool import ConnectionPool

_POOL: Optional[ConnectionPool] = None


def _get_pool() -> ConnectionPool:
	"""Return a lazily-created module-level connection pool."""
	global _POOL
	if _POOL is None:
		_POOL = ConnectionPool(
			kwargs={
				"host": os.environ.get("POSTGRES_HOST", "postgres"),
				"port": os.environ.get("POSTGRES_PORT", "5432"),
				"dbname": os.environ.get("POSTGRES_DB", "postgres"),
				"user": os.environ.get("POSTGRES_USER", "postgres"),
				"password": os.environ.get("POSTGRES_PASSWORD", ""),
			},
			min_size=1,
			max_size=10,
			open=True,
		)
	return _POOL


def close_pool() -> None:
	"""Close the module-level connection pool. Optional cleanup for long-running processes."""
	global _POOL
	if _POOL is not None:
		_POOL.close()
		_POOL = None


def get_connection() -> AbstractContextManager[psycopg.Connection]:
	"""Return a context manager that yields a connection from the pool."""
	return _get_pool().connection()

# Column order for mass_spectra inserts/updates
MASS_SPECTRA_COLS = (
	"inchikey",
	"molecular_weight",
	"exact_mass",
	"precursor_mz",
	"precursor_type",
	"ion_mode",
	"collision_energy",
	"spectrum_type",
	"instrument",
	"instrument_type",
	"splash",
	"db_number",
	"source",
	"comments",
	"raw_data",
)


def upsert_compound(cur: Cursor[Any], compound: dict[str, str]) -> None:
	"""Upsert a single row into compounds. On conflict (inchikey), update name, inchi, smiles, formula."""
	cur.execute(
		"""
		INSERT INTO compounds (inchikey, name, inchi, smiles, formula)
		VALUES (%(inchikey)s, %(name)s, %(inchi)s, %(smiles)s, %(formula)s)
		ON CONFLICT (inchikey) DO UPDATE SET
			name = EXCLUDED.name,
			inchi = EXCLUDED.inchi,
			smiles = EXCLUDED.smiles,
			formula = EXCLUDED.formula
		""",
		compound,
	)


def upsert_mass_spectrum(cur: Cursor[Any], row: dict[str, Any]) -> None:
	"""Upsert a single row into mass_spectra. Row must include all keys in MASS_SPECTRA_COLS (e.g. raw_data as Jsonb)."""
	cur.execute(
		"""
		INSERT INTO mass_spectra (
			inchikey, molecular_weight, exact_mass, precursor_mz, precursor_type,
			ion_mode, collision_energy, spectrum_type, instrument, instrument_type,
			splash, db_number, source, comments, raw_data
		) VALUES (
			%(inchikey)s, %(molecular_weight)s, %(exact_mass)s, %(precursor_mz)s, %(precursor_type)s,
			%(ion_mode)s, %(collision_energy)s, %(spectrum_type)s, %(instrument)s, %(instrument_type)s,
			%(splash)s, %(db_number)s, %(source)s, %(comments)s, %(raw_data)s
		)
		ON CONFLICT (inchikey, db_number, source) DO UPDATE SET
			molecular_weight = EXCLUDED.molecular_weight,
			exact_mass = EXCLUDED.exact_mass,
			precursor_mz = EXCLUDED.precursor_mz,
			precursor_type = EXCLUDED.precursor_type,
			ion_mode = EXCLUDED.ion_mode,
			collision_energy = EXCLUDED.collision_energy,
			spectrum_type = EXCLUDED.spectrum_type,
			instrument = EXCLUDED.instrument,
			instrument_type = EXCLUDED.instrument_type,
			splash = EXCLUDED.splash,
			comments = EXCLUDED.comments,
			raw_data = EXCLUDED.raw_data
		""",
		{k: row.get(k) for k in MASS_SPECTRA_COLS},
	)

__all__ = [
	"MASS_SPECTRA_COLS",
	"close_pool",
	"get_connection",
	"upsert_compound",
	"upsert_mass_spectrum",
]
