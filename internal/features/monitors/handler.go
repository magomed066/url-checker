package monitors

import (
	"errors"
	"net/http"

	monitorsDTO "github.com/magomed066/url-checker/internal/dto/monitors"
	"github.com/magomed066/url-checker/internal/json"
)

type handler struct {
	service service
}

func NewHandler(service service) handler {
	return handler{
		service: service,
	}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var body monitorsDTO.CreateMonitorRequest

	if err := json.ReadJson(w, r, &body); err != nil {
		json.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	monitor, err := h.service.Create(r.Context(), body)
	if err != nil {
		if errors.Is(err, errors.New("invalid monitor input")) {
			json.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		json.WriteJSONError(
			w,
			http.StatusInternalServerError,
			"could not create monitor",
		)
		return
	}

	json.WriteJSON(w, http.StatusCreated, monitor)
}
