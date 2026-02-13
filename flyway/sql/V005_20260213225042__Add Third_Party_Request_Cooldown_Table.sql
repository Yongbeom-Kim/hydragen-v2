-- Add your migration SQL here
CREATE TABLE third_party_cooldown (
	origin TEXT NOT NULL,
	unique_key TEXT NOT NULL, -- some unique key, may differ based on format + chemical identifier.
	last_requested TIMESTAMPTZ NOT NULL,
	current_cooldown_duration_hours INTEGER NOT NULL CHECK (current_cooldown_duration_hours >= 0),
	earliest_next_request TIMESTAMPTZ NOT NULL,

	PRIMARY KEY (origin, unique_key)
);
