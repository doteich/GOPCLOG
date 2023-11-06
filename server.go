package main

import (
	"net/http"

	"github.com/doteich/OPC-UA-Logger/controller"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/gopcua/opcua/ua"
)

func server() *chi.Mux {
	r := chi.NewRouter()
	corsM := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		MaxAge:         300,
	})
	r.Use(corsM.Handler)
	r.Group(func(r chi.Router) {
		r.Mount("/api/v1", apiRoutes())
	})

	return r
}

func apiRoutes() chi.Router {
	r := chi.NewRouter()
	r.Post("/modify_node", ModifyOPC)
	return r
}

type ModifyOPCRequest struct {
	NodeId string      `json:"nodeId"`
	Value  interface{} `json:"value"`
}

func ModifyOPC(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, err)
		}
	}()

	var req ModifyOPCRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}

	resp, err := controller.WriteNode(req.NodeId, req.Value)
	if err != nil {
		//TODO: log error
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrBadRequest)
		return
	}

	for _, result := range resp.Results {
		if result != ua.StatusOK {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp)
			return
		}
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}
