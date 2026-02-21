package chemicalimageresolver

import (
	"context"
	"errors"
	"hydragen-v2/server/internal/domain"
)

type ProviderType string

type MimeType string

type Image struct {
	Bytes    []byte
	MimeType MimeType
}

var ErrNotFound = errors.New("Not Found")

type RequestCooldownStore interface {
	OnCooldown(ctx context.Context, provider ProviderType, c domain.CompoundMetadata) (bool, error)
	Add(ctx context.Context, provider ProviderType, c domain.CompoundMetadata) error
	Remove(ctx context.Context, provider ProviderType, c domain.CompoundMetadata) error
}

type CompoundMetadataStore interface {
	Get(ctx context.Context, inchiKey string) (*domain.CompoundMetadata, error)
}

type ImageCache interface {
	Fetch(ctx context.Context, provider ProviderType, c domain.CompoundMetadata) (img *Image, found bool, err error)
	Save(ctx context.Context, provider ProviderType, c domain.CompoundMetadata, img *Image, imgMimeType string) error
}

type ThirdPartyProvider interface {
	FetchImage(ctx context.Context, c domain.CompoundMetadata) (*Image, error)
}
