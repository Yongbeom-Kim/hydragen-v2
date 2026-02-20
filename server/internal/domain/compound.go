package domain

import "hydragen-v2/server/internal/origin"

type CompoundMetadata struct {
	InchiKey        string   `json:"inchiKey"`
	Name            string   `json:"name"`
	Inchi           string   `json:"inchi"`
	Smiles          string   `json:"smiles"`
	Formula         string   `json:"formula"`
	MolecularWeight *float64 `json:"molecularWeight"`
	HasMassSpectrum bool     `json:"hasMassSpectrum"`
	ImageUrl        string   `json:"imageUrl"`
}

func (c *CompoundMetadata) AddImageUrl() {
	c.ImageUrl = origin.BACKEND_URL_PREFIX + "compounds/" + c.InchiKey + "/image"
}

// NewCompoundMetadata returns a CompoundMetadata instance with ImageUrl set automatically.
func NewCompoundMetadata(
	inchiKey, name, inchi, smiles, formula string,
	molecularWeight *float64,
	hasMassSpectrum bool,
) CompoundMetadata {
	c := CompoundMetadata{
		InchiKey:        inchiKey,
		Name:            name,
		Inchi:           inchi,
		Smiles:          smiles,
		Formula:         formula,
		MolecularWeight: molecularWeight,
		HasMassSpectrum: hasMassSpectrum,
	}
	c.AddImageUrl()
	return c
}
