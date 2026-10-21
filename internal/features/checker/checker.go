package checker

import (
	"context"
	"net/http"
	"time"

	"github.com/magomed066/url-checker/internal/store"
)

type CheckResult struct {
	MonitorID  string
	StatusCode int
	Duration   time.Duration
	CheckedAt  time.Time
	Err        error
}

func Check(ctx context.Context, client *http.Client, monitor store.Monitor) CheckResult {
	result := CheckResult{
		MonitorID: monitor.ID,
		CheckedAt: time.Now(),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, monitor.Url, nil)
	if err != nil {
		result.Duration = time.Since(result.CheckedAt)
		result.Err = err
		return result
	}

	response, err := client.Do(req)
	result.Duration = time.Since(result.CheckedAt)

	if err != nil {
		result.Err = err
		return result
	}

	defer response.Body.Close()

	result.StatusCode = response.StatusCode

	return result
}

func (r CheckResult) Outcome() string {
	if r.Err != nil {
		return "ERROR"
	}

	if r.StatusCode >= 200 && r.StatusCode < 400 {
		return "UP"
	}

	return "DOWN"
}
