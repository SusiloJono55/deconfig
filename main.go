package main

import (
	"deconfig/repository"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func main() {
	err := repository.OpenStore()
	if err != nil {
		log.Fatal(err)
	}
	defer repository.CloseStore()
	repository.LoadFileConfiguration()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(OptionMethod)

	r.Get("/config/info/all", getInfoAll)
	r.Get("/config/info/{_type:[a-zA-Z0-9]+}", getInfo)
	r.Route("/config", func(r chi.Router) {
		r.Use(hasOwner)
		r.Post("/add", addNewConfig)
		r.Get("/get-by-name/{configName:[a-zA-Z0-9_-]+}", getConfigByName)
		r.Get("/get/{configKey:[a-zA-Z0-9_-]+}", getConfig)
		r.Put("/update/{configKey:[a-zA-Z0-9_-]+}", updateConfig)
		r.Delete("/delete/{configKey:[a-zA-Z0-9_-]+}", deleteConfig)
		r.Get("/all", getAllConfig)
	})

	http.ListenAndServe(":3333", r)
}

func OptionMethod(next http.Handler) http.Handler {
	fn := func(res http.ResponseWriter, req *http.Request) {
		if req.Method == "OPTIONS" {
			res.Header().Set("Access-Control-Allow-Origin", "*")
			res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			res.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			res.WriteHeader(http.StatusAccepted)
			return
		}
		next.ServeHTTP(res, req)
	}
	return http.HandlerFunc(fn)
}

func hasOwner(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		owner := r.Header.Get("owner")
		if owner == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if owner == "administrator" {
			owner = ""
		}
		r.Header.Set("owner", owner)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func getInfoAll(w http.ResponseWriter, r *http.Request) {
	resp := repository.GetInfoAll()
	render.Status(r, http.StatusOK)
	render.Respond(w, r, resp)
}

func getInfo(w http.ResponseWriter, r *http.Request) {
	_type := chi.URLParam(r, "_type")

	tp, ok := repository.MapType[_type]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	resp := repository.GetInfo(tp)
	render.Status(r, http.StatusOK)
	render.Respond(w, r, resp)
}

func getAllConfig(w http.ResponseWriter, r *http.Request) {
	owner := r.Header.Get("owner")
	rtype := r.URL.Query().Get("type")

	if rtype == "all" {
		res, err := repository.GetAll(owner)
		if err != nil {
			render.Render(w, r, repository.ErrRender(err))
			return
		}
		render.Status(r, http.StatusOK)
		render.Respond(w, r, res)
	} else if tp, ok := repository.MapType[rtype]; ok {
		res, err := repository.GetAllByType(tp, owner)
		if err != nil {
			render.Render(w, r, repository.ErrRender(err))
			return
		}
		render.Status(r, http.StatusOK)
		render.Respond(w, r, res)

	} else {
		err := errors.New("No Parameter type found in requst of with empty value")
		render.Render(w, r, repository.ErrInvalidRequest(err))
		return
	}

}

func getConfigByName(w http.ResponseWriter, r *http.Request) {
	owner := r.Header.Get("owner")
	configName := chi.URLParam(r, "configName")
	config := &repository.Config{}
	err := config.GetByName(configName, owner)

	var resp map[string]interface{}
	if err != nil {
		resp = map[string]interface{}{
			"found": false,
		}
	} else {
		resp = map[string]interface{}{
			"found":   true,
			"_source": config.Value,
		}
	}
	render.Status(r, http.StatusOK)
	render.Respond(w, r, resp)
}

func getConfig(w http.ResponseWriter, r *http.Request) {
	owner := r.Header.Get("owner")
	configKey := chi.URLParam(r, "configKey")
	config := &repository.Config{}
	err := config.GetOne(configKey, owner)

	var resp map[string]interface{}
	if err != nil {
		resp = map[string]interface{}{
			"found": false,
		}
	} else {
		resp = map[string]interface{}{
			"found":   true,
			"_source": config.Value,
		}
	}
	render.Status(r, http.StatusOK)
	render.Respond(w, r, resp)
}

func deleteConfig(w http.ResponseWriter, r *http.Request) {
	owner := r.Header.Get("owner")
	configKey := chi.URLParam(r, "configKey")
	config := &repository.Config{}
	err := config.GetOne(configKey, owner)
	if err != nil {
		render.Render(w, r, repository.SetRenderer(err, 500, "Failed to delete data."))
		return
	}
	err = config.Delete()
	if err != nil {
		render.Render(w, r, repository.SetRenderer(err, 500, "Failed to delete data."))
		return
	}
	render.Status(r, http.StatusOK)
	render.Respond(w, r, "OK")
}

func updateConfig(w http.ResponseWriter, r *http.Request) {
	owner := r.Header.Get("owner")

	var rBody repository.ReqBody
	json.NewDecoder(r.Body).Decode(&rBody)

	configKey := chi.URLParam(r, "configKey")
	config := new(repository.Config)
	err := config.GetOne(configKey, owner)
	if err != nil {
		render.Render(w, r, repository.ErrInvalidRequest(err))
		return
	}

	if tp, ok := repository.MapType[rBody.CType]; ok && tp == config.Type {
		config.Name = rBody.CName
		config.MapConfig(rBody)
	} else {
		err := errors.New("Undefined connection type")
		render.Render(w, r, repository.ErrInvalidRequest(err))
		return
	}

	err = config.Update()
	if err != nil {
		render.Render(w, r, repository.SetRenderer(err, 500, "Failed to update data."))
		return
	}
	render.Status(r, http.StatusOK)
	render.Respond(w, r, "OK")
}

func addNewConfig(w http.ResponseWriter, r *http.Request) {
	owner := r.Header.Get("owner")

	var rBody repository.ReqBody
	json.NewDecoder(r.Body).Decode(&rBody)

	found := new(repository.Config)
	err := found.GetByName(rBody.CName, owner)
	if err != nil && err.Error() != "No data found" {
		render.Render(w, r, repository.SetRenderer(err, 500, "Failed to create new connection."))
		return
	}
	if err == nil && repository.TypeStr[found.Type] == rBody.CType {
		err = errors.New("Error conflict")
		render.Render(w, r, repository.SetRenderer(err, 409, "Connection with this name and type has been exist."))
		return
	}

	var config *repository.Config
	if tp, ok := repository.MapType[rBody.CType]; ok {
		config = repository.NewConfiguration(rBody.CName, owner, tp, rBody)
	} else {
		err := errors.New("Undefined connection type")
		render.Render(w, r, repository.ErrInvalidRequest(err))
		return
	}

	err = config.Create()
	if err != nil {
		render.Render(w, r, repository.SetRenderer(err, 500, "Failed to create new connection."))
		return
	}
	render.Status(r, http.StatusOK)
	render.Respond(w, r, "OK")
}
