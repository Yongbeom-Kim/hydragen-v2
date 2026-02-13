export type CompoundListItem = {
	inchiKey: string;
	name: string;
	inchi: string;
	smiles: string;
	formula: string;
	molecularWeight: number | null;
	hasMassSpectrum: boolean;
};

export type CompoundsResponse = {
	items: CompoundListItem[];
	page: number;
	pageSize: number;
	total: number;
};

export type CompoundDetail = {
	inchiKey: string;
	name: string;
	inchi: string;
	smiles: string;
	formula: string;
	hasMassSpectrum: boolean;
	massSpectrumHref?: string;
};

export type MassSpectrumItem = {
	id: number;
	inchiKey: string;
	molecularWeight: number;
	exactMass: number | null;
	precursorMz: number | null;
	precursorType: string | null;
	ionMode: string | null;
	collisionEnergy: string | null;
	spectrumType: string | null;
	instrument: string | null;
	instrumentType: string | null;
	splash: string | null;
	dbNumber: string;
	source: string;
	comments: string | null;
	mZ: number[];
	peaks: number[];
};

export type MassSpectrumResponse = {
	inchiKey: string;
	count: number;
	items: MassSpectrumItem[];
};
