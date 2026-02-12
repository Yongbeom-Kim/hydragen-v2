import time

from facade.db import upsert_compounds_batch, upsert_mass_spectra_batch
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

# Scale factor for m/z to store as int4 (4 decimal places)
MZ_SCALE = 10_000


def metadata_to_compounds_table_row(metadata: dict[str, str]) -> dict[str, str] | None:
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


def data_to_mass_spectra_table_row(
    metadata: dict[str, str],
    m_z_arr: list[int],
    intensity_arr: list[int]
) -> dict[str, str | float | None]:
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
		
	row["m_z"] = m_z_arr
	row["peaks"] = intensity_arr
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
				compounds_batch: list[dict[str, str]] = []
				rows_batch: list[dict] = []
				for raw_item in self._get_dataset_raw_items():
					metadata, (m_z_arr, intensity_arr) = self._parse_raw_item(raw_item)
					compound = metadata_to_compounds_table_row(metadata)
					if compound is None:
						continue

					mass_spec = data_to_mass_spectra_table_row(metadata, m_z_arr, intensity_arr)
					if not mass_spec.get("inchikey") or mass_spec.get("molecular_weight") is None or not mass_spec.get("db_number"):
						continue
					mass_spec["source"] = self.uniq_key
					
					compounds_batch.append(compound)
					rows_batch.append(mass_spec)

					if len(rows_batch) >= self.batch_size:
						self.batch_write_to_db(cur, conn, compounds_batch, rows_batch)
						compounds_batch = []
						rows_batch = []
				if rows_batch:
					self.batch_write_to_db(cur, conn, compounds_batch, rows_batch)

	def batch_write_to_db(self, cur, conn, compounds_batch: list, rows_batch: list) -> None:
		"""Upsert one batch of compounds and mass spectra (deduped in SQL), then commit."""
		upsert_compounds_batch(cur, compounds_batch)
		upsert_mass_spectra_batch(cur, rows_batch)
		self._row_count += len(rows_batch)
		print(f"Committed {self._row_count} records so far.", flush=True)
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
	def _parse_raw_item(raw_item: bytes) -> tuple[dict[str, str], tuple[list[float], list[float]]]:
		"""Parse a MassBank record: metadata (Field: Value) and peak data (m/z intensity)."""
		text = raw_item.decode("utf-8")
		lines = text.strip().splitlines()

		metadata: dict[str, str] = {}
		m_z_arr: list[float] = []
		intensity_arr: list[int] = []

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
						m_z_arr.append(float(parts[0]))
						intensity_arr.append(float(parts[1]))
					except ValueError:
						pass

		return metadata, (m_z_arr, intensity_arr)