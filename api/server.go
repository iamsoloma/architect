package api

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	containerd "github.com/containerd/containerd/v2/client"
	"github.com/docker/docker/client"
)

type Config struct {
	ListenAddr                  string
	AddrOfContainerDaemonSocket string
}

type Server struct {
	*Config
	Started time.Time

	ContainerDaemon *containerd.Client
	DockerClient    *client.Client
}

func NewServer(config Config) (*Server, error) {
	containerdClient, err := containerd.New(config.AddrOfContainerDaemonSocket)
	if err != nil {
		return nil, errors.New("failed to connect to containerd: " + err.Error())
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errors.New("failed to connect to docker: " + err.Error())
	}

	return &Server{
		Config:          &config,
		ContainerDaemon: containerdClient,
		DockerClient:    dockerClient,
	}, nil
}

func (s *Server) Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", s.Health)
	mux.HandleFunc("GET /namespaces", s.ListNamespaces)
	mux.HandleFunc("GET /images", s.ListImages)
	mux.HandleFunc("POST /images", s.PullImage)
	mux.HandleFunc("DELETE /namespace", s.DeleteNamespace)
	mux.HandleFunc("POST /namespace", s.CreateNamespace)

	server := http.Server{
		Addr:    s.Config.ListenAddr,
		Handler: mux,
	}

	s.Started = time.Now().UTC()

	slog.Info("api is running", "address", s.Config.ListenAddr)
	err := server.ListenAndServe()
	if err != nil {
		s.ContainerDaemon.Close()
		s.DockerClient.Close()
		slog.Error("API stoped", "error", err)
	}

}
