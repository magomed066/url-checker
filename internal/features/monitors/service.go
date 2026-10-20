package monitors

import (
	"context"

	"github.com/google/uuid"
	monitorsDTO "github.com/magomed066/url-checker/internal/dto/monitors"
	"github.com/magomed066/url-checker/internal/store"
)

type service struct {
	store *store.Store
}

func NewService(store *store.Store) service {
	return service{
		store: store,
	}
}

func (s *service) Create(ctx context.Context, payload monitorsDTO.CreateMonitorRequest) {
	id := uuid.New()

	s.store.Set(id.String(), store.Monitor{
		ID:              id.String(),
		Url:             payload.Url,
		IntervalSeconds: payload.IntervalSeconds,
		Enabled:         true,
	})
}
