package server

import (
	"context"
	"encoding/json"
	"net/http"

	"Yokify/podman"

	"Yokify/ui"

	"go.podman.io/podman/v6/pkg/domain/entities/types"
)

type runtime interface {
	GetContainers(ctx context.Context) ([]types.ListContainer, error)
	RemoveContainer(ctx context.Context, container string, force bool) error
	StartContainer(ctx context.Context, container string) error
	StopContainer(ctx context.Context, container string, timeout uint) error
	KillContainer(ctx context.Context, container string) error
	RestartContainer(ctx context.Context, container string) error
	CreateContainerFromImage(ctx context.Context, image string, options podman.CreateContainerOptions) (types.ContainerCreateResponse, error)
	CreateContainerFromFile(ctx context.Context, projectRoot string, name string, containerOptions podman.CreateContainerOptions, imageOptions podman.CreateImageOptions) (types.ContainerCreateResponse, error)
	GetImages(ctx context.Context) ([]*types.ImageSummary, error)
	RemoveImages(ctx context.Context, imagesList []string) error
}

type createContainerRequest struct {
	Image string               `json:"image"`
	Name  string               `json:"name"`
	Ports []podman.PortMapping `json:"ports"`
	Env   map[string]string    `json:"env"`
	Start bool                 `json:"start"`
}

type removeImagesRequest struct {
	Images []string `json:"images"`
}

type Server struct{ rt runtime }

func New(rt runtime) *Server {
	return &Server{rt: rt}
}

func (s *Server) GetContainers(w http.ResponseWriter, r *http.Request) {
	list, err := s.rt.GetContainers(r.Context())
	writeJSON(w, list, err)
}

func (s *Server) StartContainer(w http.ResponseWriter, r *http.Request) {
	err := s.rt.StartContainer(r.Context(), r.PathValue("id"))
	writeJSON(w, nil, err)
}

func (s *Server) StopContainer(w http.ResponseWriter, r *http.Request) {
	err := s.rt.StopContainer(r.Context(), r.PathValue("id"), 10)
	writeJSON(w, nil, err)
}

func (s *Server) RemoveContainer(w http.ResponseWriter, r *http.Request) {
	err := s.rt.RemoveContainer(r.Context(), r.PathValue("id"), true)
	writeJSON(w, nil, err)
}

func (s *Server) CreateContainerFromImage(w http.ResponseWriter, r *http.Request) {
	req, err := readJSON[createContainerRequest](w, r)
	if err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.Image == "" {
		http.Error(w, "image is required", http.StatusBadRequest)
		return
	}

	response, err := s.rt.CreateContainerFromImage(r.Context(), req.Image, podman.CreateContainerOptions{
		Name:  req.Name,
		Ports: req.Ports,
		Env:   req.Env,
		Start: req.Start,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, response, err)
}

func (s *Server) GetImages(w http.ResponseWriter, r *http.Request) {
	list, err := s.rt.GetImages(r.Context())
	writeJSON(w, list, err)
}

func (s *Server) RemoveImages(w http.ResponseWriter, r *http.Request) {
	req, err := readJSON[removeImagesRequest](w, r)
	if err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	err = s.rt.RemoveImages(r.Context(), req.Images)
	writeJSON(w, nil, err)
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/containers", s.GetContainers)
	mux.HandleFunc("GET /api/containers/{id}/start", s.StartContainer)
	mux.HandleFunc("GET /api/containers/{id}/stop", s.StopContainer)
	mux.HandleFunc("DELETE /api/containers/{id}", s.RemoveContainer)
	mux.HandleFunc("POST /api/containers", s.CreateContainerFromImage)
	mux.HandleFunc("GET /api/images", s.GetImages)
	mux.HandleFunc("DELETE /api/images", s.RemoveImages)
	mux.Handle("/", spa(ui.FS()))
	return mux
}

func writeJSON(w http.ResponseWriter, v any, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func readJSON[T any](w http.ResponseWriter, r *http.Request) (T, error) {
	var v T
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	err := json.NewDecoder(r.Body).Decode(&v)
	return v, err
}
