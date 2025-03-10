package api

import (
	"context"
	"encoding/json"
	"net/http"

	//"os"

	"github.com/containerd/containerd/v2/pkg/namespaces"
)

type Namespace struct {
	Name   string `json:"name"`
	Labels map[string]string `json:"labels"`
}

func (s *Server) ListNamespaces(w http.ResponseWriter, r *http.Request) {
	resp := []Namespace{}

	ctx := context.Background()

	namespaces, err := s.ContainerDaemon.NamespaceService().List(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, name := range namespaces {
		labels, err := s.ContainerDaemon.NamespaceService().Labels(ctx, name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp = append(resp, Namespace{Name: name, Labels: labels})
	}

	/*dir, err := os.ReadDir("/var/lib/containerd/io.containerd.runtime.v2.task")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, directory := range dir {
		if directory.IsDir() {
			resp = append(resp, Namespace{Name: directory.Name()})
		}
	}*/

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

func (s *Server) DeleteNamespace(w http.ResponseWriter, r *http.Request) {
	req := Namespace{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	ctx := namespaces.WithNamespace(context.Background(), req.Name)

	err := s.ContainerDaemon.NamespaceService().Delete(ctx, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("ok"))
}

func (s *Server) CreateNamespace(w http.ResponseWriter, r *http.Request) {
	req := Namespace{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	ctx := namespaces.WithNamespace(context.Background(), req.Name)

	err := s.ContainerDaemon.NamespaceService().Create(ctx, req.Name, req.Labels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("ok"))
}
