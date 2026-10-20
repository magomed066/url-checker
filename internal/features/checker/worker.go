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
	results chan<- CheckResult,
) {

}
