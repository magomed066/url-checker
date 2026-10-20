package monitors

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	monitorsDTO "github.com/magomed066/url-checker/internal/dto/monitors"
	"github.com/magomed066/url-checker/internal/store"
)

type service struct {
	store  *store.Store
	client *http.Client
}

func NewService(store *store.Store, client *http.Client) service {
	return service{
		store:  store,
		client: client,
	}
}

func (s *service) Create(ctx context.Context, payload monitorsDTO.CreateMonitorRequest) (store.Monitor, error) {
	if err := ctx.Err(); err != nil {
		return store.Monitor{}, err
	}

	target := strings.TrimSpace(payload.Url)
	parsed, err := url.Parse(target)

	if err != nil ||
		(parsed.Scheme != "http" && parsed.Scheme != "https") ||
		parsed.Hostname() == "" {
		return store.Monitor{}, fmt.Errorf(
			"%w: URL must use http or https and include a hostname",
			errors.New("invalid monitor input"),
		)
	}

	if payload.IntervalSeconds < 5 || payload.IntervalSeconds > 86400 {
		return store.Monitor{}, fmt.Errorf(
			"%w: interval must be between 5 and 86400 seconds",
			errors.New("invalid monitor input"),
		)
	}

	monitor := store.Monitor{
		ID:              uuid.NewString(),
		Url:             target,
		IntervalSeconds: payload.IntervalSeconds,
		Enabled:         true,
	}

	if _, err := s.store.Set(monitor.ID, monitor); err != nil {
		return store.Monitor{}, fmt.Errorf("create monitor: %w", err)
	}

	return monitor, nil
}

func (s *service) GetAll(ctx context.Context) ([]store.Monitor, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	monitors := s.store.List()

	return monitors, nil
}
