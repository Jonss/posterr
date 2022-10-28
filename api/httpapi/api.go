package httpapi

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Jonss/posterr/config"
	"github.com/gorilla/mux"
)

type Services struct {
}

type HttpServer struct {
	router        *mux.Router
	config        config.Config
	restValidator *Validator
	services      Services
}

func NewHttpServer(
	r *mux.Router,
	cfg config.Config,
	s Services,
) *HttpServer {
	h := &HttpServer{
		router:        r,
		config:        cfg,
		services:      s,
	}
	return h
}

func (h *HttpServer) Start() {
	v, err := NewValidator()
	if err != nil {
		log.Fatalf("error on create validator. error=(%v)", err)
	}
	h.restValidator = v
	h.routes()
}

func (h *HttpServer) routes() {
	api := h.router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Posterr!")
	})
}