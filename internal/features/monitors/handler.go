package monitors

import (
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

	err := json.ReadJson(w, r, &body)
	if err != nil {
		json.WriteJSONError(w, http.StatusBadRequest, "Bad request")
		return
	}

	h.service.Create(r.Context(), body)

	json.WriteJSON(w, http.StatusCreated, "OK")
}
