package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
)

type Image struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels"`
}

type PullReq struct {
	Image    string `json:"image"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type RemoveReq struct {
	Image string `json:"image"`
}

func (s *Server) ListImages(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	images, err := s.DockerClient.ImageList(ctx, image.ListOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(images)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
	}

}

func (s *Server) PullImage(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	req := PullReq{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authConfig := registry.AuthConfig{
		Username: req.Username,
		Password: req.Password,
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	pulling, err := s.DockerClient.ImagePull(ctx, req.Image, image.PullOptions{RegistryAuth: authStr})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer pulling.Close()
	io.Copy(w, pulling)
}

func (s *Server) RemoveImage(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	req := RemoveReq{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	removedImages, err := s.DockerClient.ImageRemove(ctx, req.Image, image.RemoveOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(removedImages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
	}

}
