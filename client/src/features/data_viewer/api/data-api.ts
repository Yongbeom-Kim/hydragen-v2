import type {
	CompoundDetail,
	CompoundsResponse,
	MassSpectrumResponse,
} from "../types/data";

const API_BASE = "/v2/api";

async function fetchJSON<T>(url: string): Promise<T> {
	const response = await fetch(url);
	if (!response.ok) {
		throw new Error(`Request failed (${response.status})`);
	}

	return response.json() as Promise<T>;
}

export function getCompoundList(page: number, pageSize: number) {
	const params = new URLSearchParams({
		page: String(page),
		pageSize: String(pageSize),
	});
	return fetchJSON<CompoundsResponse>(
		`${API_BASE}/compounds?${params.toString()}`,
	);
}

export function getCompoundDetail(inchiKey: string) {
	return fetchJSON<CompoundDetail>(`${API_BASE}/compounds/${inchiKey}`);
}

export function getMassSpectrum(inchiKey: string) {
	return fetchJSON<MassSpectrumResponse>(
		`${API_BASE}/mass-spectra/${inchiKey}`,
	);
}
