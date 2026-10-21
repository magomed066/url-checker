package checker

import (
	"context"
	"net/http"

	"github.com/magomed066/url-checker/internal/store"
)

func Worker(
	ctx context.Context,
	client *http.Client,
	jobs <-chan store.Monitor,
	result chan<- CheckResult,
) {

	for {
		select {
		case <-ctx.Done():
			return

		case job, ok := <-jobs:
			if !ok {
				return
			}

			if ctx.Err() != nil {
				return
			}

			checked := Check(ctx, client, job)

			select {
			case <-ctx.Done():
				return

			case result <- checked:
			}
		}
	}
}
