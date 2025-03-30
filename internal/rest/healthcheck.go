package rest

import (
	"context"
	"net/http"
	"time"
)

type HealthCheckResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

func (s *RestService) HealthCheck(w http.ResponseWriter, r *http.Request) {
	rsp := HealthCheckResponse{
		Status: "SUCCESS",
		Msg:    "Scanning Service Is Helth.",
	}
	cxt, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(s.listener.config.WriteTimeout))
	defer cancel()
	r = r.WithContext(cxt)
	WriteResponse(w, rsp, http.StatusOK)
}
