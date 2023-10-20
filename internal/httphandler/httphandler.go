package httphandler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rendau/s3-uploader/internal/core"
)

type handlerSt struct {
	cr *core.Core
}

func NewHttpHandler(cr *core.Core) http.Handler {
	r := chi.NewRouter()

	// healthcheck
	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// docs
	docFS := http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs")))
	r.Get("/docs/*", func(w http.ResponseWriter, r *http.Request) {
		docFS.ServeHTTP(w, r)
	})

	// --------------------

	h := &handlerSt{
		cr: cr,
	}

	// save
	r.Post("/save", h.Save)

	return r
}
