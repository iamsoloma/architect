package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"architect/system"
)

type Health struct {
	Status          string `json:"status"`
	CurrentTime     string `json:"currentTime"`
	Uptime          string `json:"upTime"`
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
		CurrentTime:     time.Now().UTC().String(),
		Uptime:          time.Now().UTC().Sub(s.Started).String(),
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

func (s *Server) Info(w http.ResponseWriter, r *http.Request) {
	info, err := s.DockerClient.Info(context.Background())

	jsonResp, err := json.Marshal(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
	}
}