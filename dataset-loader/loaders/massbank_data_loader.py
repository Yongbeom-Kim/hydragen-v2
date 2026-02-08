from psycopg.types.json import Jsonb
import time

from .dataset_loader_base import DatasetLoaderBase

# MassBank metadata key -> mass_spectra table column name
METADATA_KEY_TO_TABLE_KEY = {
	"Collision_energy": "collision_energy",
	"Comments": "comments",
	"DB#": "db_number",
	"ExactMass": "exact_mass",
	"InChIKey": "inchikey",
	"Instrument": "instrument",
	"Instrument_type": "instrument_type",
	"Ion_mode": "ion_mode",
	"MW": "molecular_weight",
	"PrecursorMZ": "precursor_mz",
	"Precursor_type": "precursor_type",
	"Splash": "splash",
	"Spectrum_type": "spectrum_type",
}

def compound_row_from_metadata(metadata: dict[str, str]) -> dict[str, str] | None:
	"""Build a compounds table row from MassBank metadata. Returns None if any required field is missing."""
	row = {
		"inchikey": metadata.get("InChIKey", "").strip(),
		"name": metadata.get("Name", "").strip(),
		"inchi": metadata.get("InChI", "").strip(),
		"smiles": metadata.get("SMILES", "").strip(),
		"formula": metadata.get("Formula", "").strip(),
	}
	if not all(row.values()) or len(row["inchikey"]) != 27:
		return None
	return row


def map_metadata_to_table_row(metadata: dict[str, str]) -> dict[str, str | float | None]:
	"""Map MassBank metadata keys to mass_spectra table columns. Drops unmapped keys."""
	row: dict[str, str | float | None] = {}
	for meta_key, value in metadata.items():
		table_key = METADATA_KEY_TO_TABLE_KEY.get(meta_key)
		if table_key is None:
			continue
		if table_key in ("precursor_mz", "molecular_weight", "exact_mass"):
			try:
				row[table_key] = float(value)
			except ValueError:
				row[table_key] = None
		else:
			row[table_key] = value
	return row


MASS_SPECTRA_COLS = (
	"inchikey", "molecular_weight", "exact_mass", "precursor_mz", "precursor_type",
	"ion_mode", "collision_energy", "spectrum_type", "instrument", "instrument_type",
	"splash", "db_number", "source", "comments", "raw_data",
)


def insert_compound(cur, compound: dict[str, str]) -> None:
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


def insert_mass_spectrum(cur, row: dict) -> None:
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


class MassbankDataLoader(DatasetLoaderBase):
	def __init__(self, uniq_key: str, source_url: str, batch_size: int = 100, batch_delay: float = 0.0):
		super().__init__(uniq_key, source_url)
		self._row_count = 0
		self.batch_size = batch_size
		self.batch_delay = batch_delay

	def upload_to_db(self):
		with self._get_connection() as conn:
			with conn.cursor() as cur:
				for raw_item in self._get_dataset_raw_items():
					parsed = self._parse_raw_item(raw_item)
					metadata = parsed["metadata"]
					compound = compound_row_from_metadata(metadata)
					if compound is None:
						continue
					row = map_metadata_to_table_row(metadata)
					if not row.get("inchikey") or row.get("molecular_weight") is None or not row.get("db_number"):
						continue
					row["source"] = self.uniq_key
					row["raw_data"] = Jsonb(parsed["peaks"])
					insert_compound(cur, compound)
					insert_mass_spectrum(cur, row)

					self._row_count += 1
					if self._row_count % self.batch_size == 0:
						print(f"Committed {self._row_count} records so far.")
						conn.commit()
						if self.batch_delay:
							time.sleep(self.batch_delay)

	def _get_dataset_raw_items(self):
		with open(self.dataset_path, "rb", buffering=1024 * 1024) as f:
			item_raw = []
			for line_raw in f:
				if (len(item_raw) != 0 and line_raw.startswith(b"Name:")):
					yield b"".join(item_raw)
				item_raw.append(line_raw)

	@staticmethod
	def _parse_raw_item(raw_item: bytes) -> dict:
		"""Parse a MassBank record: metadata (Field: Value) and peak data (m/z intensity)."""
		text = raw_item.decode("utf-8")
		lines = text.strip().splitlines()

		metadata: dict[str, str] = {}
		peaks: list[dict[str, float]] = []

		for line in lines:
			line = line.strip()
			if not line:
				continue

			if ":" in line:
				key, _, value = line.partition(":")
				metadata[key.strip()] = value.strip()
			
			else:
				parts = line.split()
				if len(parts) >= 2:
					try:
						mz = float(parts[0])
						intensity = float(parts[1])
						peaks.append({"m/z": mz, "intensity": intensity})
					except ValueError:
						pass

		return {"metadata": metadata, "peaks": peaks}