package domain

type MassSpectraRecord struct {
	ID              int64     `json:"id"`
	InchiKey        string    `json:"inchiKey"`
	MolecularWeight float64   `json:"molecularWeight"`
	ExactMass       *float64  `json:"exactMass"`
	PrecursorMz     *float64  `json:"precursorMz"`
	PrecursorType   *string   `json:"precursorType"`
	IonMode         *string   `json:"ionMode"`
	CollisionEnergy *string   `json:"collisionEnergy"`
	SpectrumType    *string   `json:"spectrumType"`
	Instrument      *string   `json:"instrument"`
	InstrumentType  *string   `json:"instrumentType"`
	Splash          *string   `json:"splash"`
	DBNumber        string    `json:"dbNumber"`
	Source          string    `json:"source"`
	Comments        *string   `json:"comments"`
	MZ              []float32 `json:"mZ"`
	Peaks           []int     `json:"peaks"`
}
