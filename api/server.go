package api

import (
	"errors"
	"log/slog"
	"net/http"

	containerd "github.com/containerd/containerd/v2/client"
)

type Config struct {
	ListenAddr                  string
	AddrOfContainerDaemonSocket string
}

type Server struct {
	*Config

	ContainerDaemon *containerd.Client
}

func NewServer(config Config) (*Server, error) {
	client, err := containerd.New(config.AddrOfContainerDaemonSocket)
	if err != nil {
		return nil, errors.New("failed to connect to containerd: " + err.Error())
	}

	return &Server{
		Config:          &config,
		ContainerDaemon: client,
	}, nil
}

func (s *Server) Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", s.Health)
	mux.HandleFunc("GET /namespaces", s.ListNamespaces)

	server := http.Server{
		Addr:    s.Config.ListenAddr,
		Handler: mux,
	}

	slog.Info("api is running", "address", s.Config.ListenAddr)
	err := server.ListenAndServe()
	if err != nil {
		slog.Error("API stoped", "error", err)
	}

}
