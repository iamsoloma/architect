package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	containerd "github.com/containerd/containerd/v2/client"

	"architect/system"
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
		Config: &config,
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


	slog.Info("api is running", "address", s.Config.ListenAddr )
	err := server.ListenAndServe()
	if err != nil {
		slog.Error("API stoped", "error", err)
	}

}

type Health struct {
	Status          string `json:"status"`
	UTC             string `json:"utc"`
	NodeType        string `json:"type"`
	Version         string `json:"version"`
	FreeMemory      uint64 `json:"freeMemory"`
	AvailableMemory uint64 `json:"availableMemory"`
	TotalMemory     uint64 `json:"totalMemory"`
	FreeDisk        uint64 `json:"freeDisk"`
	AvailableDisk   uint64 `json:"availableDisk"`
	TotalDisk       uint64 `json:"totalDisk"`
}

func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	memory := system.ReadMemoryStats()
	disk, err := system.ReadDiskInfo("./")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := Health{
		Status:          "ok",
		UTC:             time.Now().UTC().String(),
		NodeType:        "Architect",
		Version:         "0.0.1",
		FreeMemory:      memory.FreeMem,
		AvailableMemory: memory.AvailableMem,
		TotalMemory:     memory.TotalMem,
		FreeDisk:        disk.FreeDisk,
		AvailableDisk:   disk.AvailableDisk,
		TotalDisk:       disk.TotalDisk,
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
