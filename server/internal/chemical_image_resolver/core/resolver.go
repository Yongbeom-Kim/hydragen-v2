package chemicalimageresolver

import (
	"context"
	"log/slog"
)

type Resolver struct {
	cooldowns     RequestCooldownStore
	metadata      CompoundMetadataStore
	cache         ImageCache
	providers     map[ProviderType]ThirdPartyProvider
	providerOrder []ProviderType
}

func New(cooldowns RequestCooldownStore, metadata CompoundMetadataStore, cache ImageCache, providers map[ProviderType]ThirdPartyProvider, providerOrder []ProviderType) *Resolver {
	return &Resolver{cooldowns: cooldowns, metadata: metadata, cache: cache, providers: providers, providerOrder: providerOrder}
}

func (r *Resolver) Image(ctx context.Context, inchiKey string) (image *Image, ok bool) {
	compound_ptr, err := r.metadata.Get(ctx, inchiKey)
	if err != nil {
		slog.Error("[ChemicalImageResolver.Image] Failed to retrieve compound metadata", "inchiKey", inchiKey, "error", err)
		return nil, false
	}
	compound := *compound_ptr

	for _, providerType := range r.providerOrder {
		image, found, err := r.cache.Fetch(ctx, providerType, compound)
		if found {
			return image, true
		}
		if err != nil {
			slog.Info("[ChemicalImageResolver.Resolve] Failed to fetch image from cache", "providerType", providerType, "error", err)
		}

		onCooldown, err := r.cooldowns.OnCooldown(ctx, providerType, compound)
		if onCooldown {
			slog.Info("[ChemicalImageResolver.Resolve] Provider on cooldown, skipping", "providerType", providerType, "inchiKey", compound.InchiKey)
			continue
		}
		if err != nil {
			slog.Error("[ChemicalImageResolver.Resolve] Error checking provider cooldown", "providerType", providerType, "inchiKey", compound.InchiKey, "error", err)
			continue
		}

		provider := r.providers[providerType]
		img, err := provider.FetchImage(ctx, compound)
		if err != nil {
			slog.Info("[ChemicalImageResolver.Resolve] Error fetching image from provider", "providerType", providerType, "inchiKey", compound.InchiKey, "error", err)
			if addErr := r.cooldowns.Add(ctx, providerType, compound); addErr != nil {
				slog.Error("[ChemicalImageResolver.Resolve] Failed to add cooldown after fetch error", "providerType", providerType, "inchiKey", compound.InchiKey, "error", addErr)
			}
			continue
		}

		if err := r.cooldowns.Remove(ctx, providerType, compound); err != nil {
			slog.Error("[ChemicalImageResolver.Resolve] Failed to remove cooldown on success", "providerType", providerType, "inchiKey", compound.InchiKey, "error", err)
		}
		if err := r.cache.Save(ctx, providerType, compound, img, string(img.MimeType)); err != nil {
			slog.Error("[ChemicalImageResolver.Resolve] Failed to save image to cache", "providerType", providerType, "inchiKey", compound.InchiKey, "error", err)
		}
		return img, true
	}
	return nil, false
}
