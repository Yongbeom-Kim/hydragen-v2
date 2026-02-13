package domain

type CompoundMetadata struct {
	InchiKey        string   `json:"inchiKey"`
	Name            string   `json:"name"`
	Inchi           string   `json:"inchi"`
	Smiles          string   `json:"smiles"`
	Formula         string   `json:"formula"`
	MolecularWeight *float64 `json:"molecularWeight"`
	HasMassSpectrum bool     `json:"hasMassSpectrum"`
}
