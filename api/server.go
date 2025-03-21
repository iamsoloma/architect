package api

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/docker/docker/client"
)

type Config struct {
	ListenAddr                  string
	AddrOfContainerDaemonSocket string
}

type Server struct {
	*Config
	Started time.Time

	DockerClient *client.Client
}

func NewServer(config Config) (*Server, error) {

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errors.New("failed to connect to docker: " + err.Error())
	}

	return &Server{
		Config:       &config,
		DockerClient: dockerClient,
	}, nil
}

func (s *Server) Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", s.Health)
	mux.HandleFunc("GET /info", s.Info)
	mux.HandleFunc("GET /images", s.ListImages)
	mux.HandleFunc("POST /images", s.PullImage)
	mux.HandleFunc("DELETE /images", s.RemoveImage)
	server := http.Server{
		Addr:    s.Config.ListenAddr,
		Handler: mux,
	}

	s.Started = time.Now().UTC()

	slog.Info("api is running", "address", s.Config.ListenAddr)
	err := server.ListenAndServe()
	if err != nil {
		s.DockerClient.Close()
		slog.Error("API stoped", "error", err)
	}

}
