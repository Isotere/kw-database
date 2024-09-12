package echo

import (
	"net/http"

	"github.com/Isotere/kw-database/pkg/web/response"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	_ = r.Context()

	response.NoContent(w, r)
}
