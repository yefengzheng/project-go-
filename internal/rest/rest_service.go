package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"project-go-/internal/config"
	"time"
)

type RestService struct {
	cxt      context.Context
	listener *RestServiceListener
}
type RestServiceListener struct {
	config   RestServiceCofig
	listener *http.Server
}
type RestServiceCofig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewRestService(cfg *config.Config) *RestService {
	svc := &RestService{
		cxt: context.Background(),
		listener: &RestServiceListener{
			config: RestServiceCofig{
				Port:         cfg.Rest.Port,
				ReadTimeout:  cfg.Rest.ReadTimeout,
				WriteTimeout: cfg.Rest.WriteTimeout,
			},
		},
	}
	return svc
}

func (svc *RestService) Start() error {
	log.Printf("start new rest service")
	svc.listener.listener = &http.Server{
		Addr:         fmt.Sprintf(":%d", svc.listener.config.Port),
		Handler:      svc.Router(),
		ReadTimeout:  svc.listener.config.ReadTimeout,
		WriteTimeout: svc.listener.config.WriteTimeout,
	}
	return svc.listener.listener.ListenAndServe()
}

func (svc *RestService) Stop() error {
	log.Printf("stop new rest service")
	return svc.listener.listener.Shutdown(svc.cxt)
}

func (svc *RestService) Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/healthcheck", svc.HealthCheck).Methods("GET").Name("healthcheck")
	return router
}

func WriteResponse(respWriter http.ResponseWriter, response interface{}, status int) {
	header := respWriter.Header()
	header.Set("Content-Type", "application/json")
	respWriter.WriteHeader(status)
	err := json.NewEncoder(respWriter).Encode(response)
	if err != nil {
		switch err.(type) {
		case *json.MarshalerError, *json.UnsupportedTypeError, *json.UnsupportedValueError:
			panic("Failed to marshal response" + err.Error())
		default:
			log.Printf("Failed to marshal response")
		}
	}
}
