package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"project-go-/internal/config"
	"project-go-/internal/database"
	"project-go-/internal/task"
	"time"
)

type RestService struct {
	cxt      context.Context
	dbCtx    *database.Context
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

func NewRestService(cfg *config.Config, dbCtx *database.Context) *RestService {
	svc := &RestService{
		cxt:   context.Background(),
		dbCtx: dbCtx,
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
	router.HandleFunc("/scan", svc.HandleScanRequest).Methods("POST").Name("scan")
	router.HandleFunc("/result", svc.HandleCheckScanningResultRequest).Methods("POST").Name("result")
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
func (svc *RestService) HandleCheckScanningResultRequest(w http.ResponseWriter, r *http.Request) {
	type RequestData struct {
		Name string `json:"name"`
	}
	var req RequestData

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "Invalid request: 'name' field is required", http.StatusBadRequest)
		return
	}

	result, err := svc.dbCtx.PgsqlContext.GetScanResult(req.Name)
	if err != nil {
		http.Error(w, "No scan result found", http.StatusNotFound)
		return
	}

	WriteResponse(w, result, http.StatusOK)
}

func (svc *RestService) HandleScanRequest(w http.ResponseWriter, r *http.Request) {
	type RequestData struct {
		Name   string `json:"name"`
		SHA256 string `json:"sha256"`
	}
	type ResponseData struct {
		Result  string `json:"result"`
		Message string `json:"message"`
	}
	var req RequestData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	var rsp ResponseData
	if status, err := svc.dbCtx.RedisContext.GetValue(req.Name); err == nil && status == "up" { //first redis op get
		rsp = ResponseData{
			Result:  "ok",
			Message: "The image is scanning",
		}
		WriteResponse(w, rsp, http.StatusOK)
		return
	}
	task.CreateNewTask(req.Name, req.SHA256)
}
