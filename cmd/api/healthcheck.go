package main

import (
	"net/http"
	"time"
)

type healthResponse struct {
	Status  string    `json:"status"`
	Version string    `json:"version"`
	Time    time.Time `json:"time"`
	Host    string    `json:"host"`
}

func (a *applicationDependencies) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := healthResponse{
		Status:  "healthy",
		Version: appVersion,
		Time:    time.Now().UTC(),
		Host:    r.Host,
	}

	err := a.writeJSON(w, http.StatusOK, envelope{"data": response}, nil)
	if err != nil {
		a.logger.Error("failed to write health check response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
