package main

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/sapariduo/dataconn/repository"
)

type ConfigRequest struct {
	Index string `json:"index,omitempty"`
	Type  string `json:"type,omitempty"`
	Id    string `json:"id,omitempty"`
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/config", func(r chi.Router) {
		r.Get("/all", GetAllConfig)
	})

	http.ListenAndServe(":3333", r)
}

func GetAllConfig(w http.ResponseWriter, r *http.Request) {

	if rtype := r.URL.Query().Get("type"); rtype != "" {
		res, err := repository.GetAll(rtype)
		if err != nil {
			render.Render(w, r, ErrRender(err))
			return
		}
		render.Status(r, http.StatusFound)
		render.Respond(w, r, res)

	} else {
		err := errors.New("No Parameter type found in requst of with empty value")
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

}
