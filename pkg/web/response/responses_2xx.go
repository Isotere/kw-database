package response

import (
	"net/http"

	"github.com/go-chi/render"
)

func NoContent(w http.ResponseWriter, r *http.Request) {
	render.NoContent(w, r)
}

func OK(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, data)
}

func Created(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusCreated)
}
