import os
from pathlib import Path

def get_root_dir() -> Path:
	return Path(__file__).resolve().parent.parent
