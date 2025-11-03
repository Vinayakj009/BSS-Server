package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func (s *Server) setupPlanRoutes() {
	s.router.Post("/plans", s.handleCreatePlan)
	s.router.Get("/plans", s.handleGetPlans)
	s.router.Get("/plans/{id}", s.handleGetPlan)
	s.router.Put("/plans/{id}", s.handleUpdatePlan)
}

func (s *Server) handleCreatePlan(w http.ResponseWriter, r *http.Request) {
	var plan Plan
	if err := json.NewDecoder(r.Body).Decode(&plan); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	createdPlan, err := s.db.CreatePlan(r.Context(), plan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdPlan)
}

func (s *Server) handleGetPlans(w http.ResponseWriter, r *http.Request) {
	page := 1
	pageSize := 10

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := r.URL.Query().Get("pageSize"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	pageableRequest := PageableRequest{
		Page:     page,
		PageSize: pageSize,
	}
	plansPage, err := s.db.GetPlans(r.Context(), pageableRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(plansPage)
}

func (s *Server) handleGetPlan(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	plan, err := s.db.GetPlan(r.Context(), idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(plan)
}

func (s *Server) handleUpdatePlan(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	planId, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid plan id", http.StatusBadRequest)
		return
	}
	var plan Plan
	if err := json.NewDecoder(r.Body).Decode(&plan); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	plan.ID = planId
	updatedPlan, err := s.db.UpdatePlan(r.Context(), plan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedPlan)
}
