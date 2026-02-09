import time

from psycopg.types.json import Jsonb

from facade.db import upsert_compound, upsert_mass_spectrum
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
					upsert_compound(cur, compound)
					upsert_mass_spectrum(cur, row)

					self._row_count += 1
					if self._row_count % self.batch_size == 0:
						print(f"Committed {self._row_count} records so far.", flush=True)
						conn.commit()
						if self.batch_delay:
							time.sleep(self.batch_delay)
				conn.commit()

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