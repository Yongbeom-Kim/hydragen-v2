package chemicalimageresolver_thirdparty

import (
	"context"
	chemicalimageresolver "hydragen-v2/server/internal/chemical_image_resolver/core"
	"hydragen-v2/server/internal/domain"
	"io"
	"net/http"
)

type CactusThirdPartyProvider struct {
}

func (provider *CactusThirdPartyProvider) urlFallbackList(compound domain.CompoundMetadata) []string {

	url := func(identifier string) string {
		return "https://cactus.nci.nih.gov/chemical/structure/" + identifier + "/image"
	}

	var identifiers = []string{
		compound.InchiKey,
		compound.Smiles,
		compound.Name,
	}

	var nonEmptyIdentifiers []string
	for _, id := range identifiers {
		if id != "" {
			nonEmptyIdentifiers = append(nonEmptyIdentifiers, id)
		}
	}
	var urls []string
	for _, id := range nonEmptyIdentifiers {
		urls = append(urls, url(id))
	}
	return urls

}

func (provider *CactusThirdPartyProvider) FetchImage(ctx context.Context, compound domain.CompoundMetadata) (*chemicalimageresolver.Image, error) {
	fallbackUrls := provider.urlFallbackList(compound)
	for _, url := range fallbackUrls {
		resp, err := http.Get(url)
		if err != nil {
			continue
		}
		data, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}
		if len(data) == 0 {
			continue
		}
		mimeType := resp.Header.Get("Content-Type")
		return &chemicalimageresolver.Image{
			Bytes:    data,
			MimeType: chemicalimageresolver.MimeType(mimeType),
		}, nil
	}
	return nil, chemicalimageresolver.ErrNotFound
}
