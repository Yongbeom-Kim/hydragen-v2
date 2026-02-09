from abc import ABC, abstractmethod
from contextlib import AbstractContextManager
from dataclasses import dataclass
from typing import Optional
import urllib.request
import hashlib

import psycopg

from facade.db import get_connection
from utils.dir import get_root_dir


@dataclass
class DatasetLoaderState:
	dataset_key: str
	dataset_version: str
	status: str
	checksum: str
	updated_at: str


class DatasetLoaderBase(ABC):

	def __init__(self, uniq_key: str, source_url: str, dataset_version: str = "1"):
		self.uniq_key = uniq_key
		self.source_url = source_url
		self.dataset_version = dataset_version

	def load(self):
		self._download_dataset()
		state = self._read_loader_state()
		if state is not None and state.status == "success":
			checksum = self._get_dataset_checksum()
			if checksum is not None and checksum == state.checksum:
				return

		self._set_start_upload_state()
		try:
			self.upload_to_db()
		except Exception:
			self._set_complete_upload_state(status="failed")
			raise
		self._set_complete_upload_state()

	@abstractmethod
	def upload_to_db(self):
		"""Handle uploading the dataset to the target database when reload is needed."""
		pass

	def _download_dataset(self):
		dataset_dir = get_root_dir() / '.datasets'
		dataset_dir.mkdir(exist_ok=True)

		target_path = dataset_dir / self.uniq_key
		urllib.request.urlretrieve(self.source_url, target_path)
		self.dataset_path = target_path

	def _get_connection(self) -> AbstractContextManager[psycopg.Connection]:
		"""Return a context manager that yields a connection from the pool."""
		return get_connection()

	def _read_loader_state(self) -> Optional[DatasetLoaderState]:
		"""Read current state for this loader from dataset_loader table."""
		with self._get_connection() as conn:
			with conn.cursor() as cur:
				cur.execute(
					"""
					SELECT dataset_key, dataset_version, status, checksum, updated_at
					FROM dataset_loader
					WHERE dataset_key = %s
					""",
					(self.uniq_key,),
				)
				row = cur.fetchone()
		if row is None:
			return None
		return DatasetLoaderState(
			dataset_key=row[0],
			dataset_version=row[1],
			status=row[2],
			checksum=row[3],
			updated_at=row[4].isoformat() if hasattr(row[4], "isoformat") else str(row[4]),
		)
	
	def _set_start_upload_state(self):
		"""Set status to 'running', update checksum and timestamp before uploading to db.
		Upserts a row if none exists (first run).
		"""
		checksum = self._get_dataset_checksum() or ""
		with self._get_connection() as conn:
			with conn.cursor() as cur:
				cur.execute(
					"""
					INSERT INTO dataset_loader (dataset_key, dataset_version, status, checksum, updated_at)
					VALUES (%s, %s, 'running', %s, NOW())
					ON CONFLICT (dataset_key) DO UPDATE SET
						updated_at = NOW(),
						checksum = EXCLUDED.checksum,
						status = 'running'
					""",
					(self.uniq_key, self.dataset_version, checksum),
				)
				conn.commit()


	def _set_complete_upload_state(self, status: str = "success"):
		"""Set status to 'success' (or other), update checksum and timestamp after uploading."""
		checksum = self._get_dataset_checksum() or ""
		with self._get_connection() as conn:
			with conn.cursor() as cur:
				cur.execute(
					"""
					UPDATE dataset_loader
					SET updated_at = NOW(),
						checksum = %s,
						status = %s
					WHERE dataset_key = %s
					""",
					(checksum, status, self.uniq_key),
				)
				conn.commit()

	def _get_dataset_checksum(self) -> Optional[str]:
		if not hasattr(self, "dataset_path") or not self.dataset_path or not self.dataset_path.exists():
			return None

		hash_sha256 = hashlib.sha256()
		with open(self.dataset_path, "rb") as f:
			for chunk in iter(lambda: f.read(4096), b""):
				hash_sha256.update(chunk)
		return hash_sha256.hexdigest()