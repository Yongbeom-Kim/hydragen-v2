package domain

import "hydragen-v2/server/utils"

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

func (c CompoundMetadata) AddImageUrl() CompoundMetadata {
	c.ImageUrl = utils.BACKEND_URL_PREFIX + "compounds/" + c.InchiKey + "/image"
	return c
}
