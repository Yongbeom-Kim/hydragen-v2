package chemicalimageresolver_thirdparty

import (
	"context"
	chemicalimageresolver "hydragen-v2/server/internal/chemical_image_resolver/core"
	"hydragen-v2/server/internal/domain"
	"io"
	"net/http"
	"strings"
)

type ChemblThirdPartyProvider struct {
}

func (provider *ChemblThirdPartyProvider) imageURL(inchiKey string) string {
	return "https://www.ebi.ac.uk/chembl/api/data/image/" + inchiKey + "?format=svg"
}

func normalizeMimeType(contentType string) string {
	if contentType == "" {
		return ""
	}
	return strings.TrimSpace(strings.Split(contentType, ";")[0])
}

func (provider *ChemblThirdPartyProvider) FetchImage(ctx context.Context, compound domain.CompoundMetadata) (*chemicalimageresolver.Image, error) {
	if compound.InchiKey == "" {
		return nil, chemicalimageresolver.ErrNotFound
	}
	url := provider.imageURL(compound.InchiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, chemicalimageresolver.ErrNotFound
	}
	if len(data) == 0 {
		return nil, chemicalimageresolver.ErrNotFound
	}
	contentType := normalizeMimeType(resp.Header.Get("Content-Type"))
	if contentType == "" {
		contentType = "image/svg+xml"
	}
	return &chemicalimageresolver.Image{
		Bytes:    data,
		MimeType: chemicalimageresolver.MimeType(contentType),
	}, nil
}
