package third_party_cactus

import (
	"hydragen-v2/server/internal/disk"
	"hydragen-v2/server/internal/domain"
	"hydragen-v2/server/utils"
	"io"
	"net/http"
)

func GetImageUrl(identifier string) string {
	return "https://cactus.nci.nih.gov/chemical/structure/" + identifier + "/image"
}

func GetImageUrlFallbackList(compound domain.CompoundMetadata) []string {
	var identifiers = []string{
		compound.InchiKey,
		compound.Smiles,
		compound.Name,
	}

	identifiers = utils.Filter(identifiers, func(id string) bool { return id != "" })
	return utils.Map(identifiers, GetImageUrl)
}

func fetchCactusImage(compound domain.CompoundMetadata) ([]byte, string, error) {
	fallbackUrls := GetImageUrlFallbackList(compound)
	for _, url := range fallbackUrls {
		resp, err := http.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			continue
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		contentType := resp.Header.Get("Content-Type")
		return data, contentType, nil
	}
	return nil, "", io.EOF // io.EOF: no image found with any identifier
}

func GetCactusImage(compound domain.CompoundMetadata) ([]byte, string, error) {
	data, mimeType, err := disk.ReadCactusImageFromDisk(compound)
	if err == nil && data != nil && mimeType != "" {
		return data, mimeType, nil
	}

	data, mimeType, err = fetchCactusImage(compound)
	if err != nil {
		return nil, "", err
	}

	_ = disk.SaveCactusImageToDisk(data, mimeType, compound)

	return data, mimeType, nil
}
