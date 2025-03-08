package api

import (
	"encoding/json"
	"net/http"
	"os"
)

type Namespace struct {
	Name string `json:"name"`
}

func (s *Server) ListNamespaces(w http.ResponseWriter, r *http.Request) {
	resp := []Namespace{}

	dir, err := os.ReadDir("/var/lib/containerd/io.containerd.runtime.v2.task")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, directory := range dir {
		if directory.IsDir() {
			resp = append(resp, Namespace{Name: directory.Name()})
		}
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
	}
}
