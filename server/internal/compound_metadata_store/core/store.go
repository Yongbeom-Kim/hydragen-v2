package compoundmetadatastore

import (
	"context"
	"hydragen-v2/server/internal/domain"
)

type MetadataStore interface {
	List(ctx context.Context, page int, pageSize int) ([]domain.CompoundMetadata, error)
	Count(ctx context.Context) (int, error)
	Get(ctx context.Context, inchiKey string) (*domain.CompoundMetadata, error)
}

type Service struct {
	store MetadataStore
}

func NewService(store MetadataStore) *Service {
	return &Service{store: store}
}

func (s *Service) List(ctx context.Context, page int, pageSize int) ([]domain.CompoundMetadata, error) {
	compounds, err := s.store.List(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}
	for i := range compounds {
		compounds[i].AddImageUrl()
	}
	return compounds, nil
}

func (s *Service) Count(ctx context.Context) (int, error) {
	return s.store.Count(ctx)
}

func (s *Service) Get(ctx context.Context, inchiKey string) (*domain.CompoundMetadata, error) {
	record, err := s.store.Get(ctx, inchiKey)
	if err != nil {
		return nil, err
	}
	record.AddImageUrl()
	return record, nil
}
