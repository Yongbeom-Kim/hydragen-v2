-- Add your migration SQL here

CREATE TABLE dataset_loader (
	dataset_key 		TEXT PRIMARY KEY,
	dataset_version TEXT NOT NULL,
  status					text NOT NULL DEFAULT 'idle', -- idle|running|success|failed
	checksum 				TEXT NOT NULL,
	updated_at      timestamp NOT NULL DEFAULT now()
);
