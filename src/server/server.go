package server

import (
	"bss/src/database"
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type PageableRequest = database.PageableRequest
type Page[V any] = database.Page[V]
type Plan = database.Plan
type Subscription = database.Subscription
type Event = database.Event

type Database interface {
	CreatePlan(ctx context.Context, plan Plan) (Plan, error)
	GetPlans(ctx context.Context, pageableRequest PageableRequest) (Page[Plan], error)
	GetPlan(ctx context.Context, id string) (Plan, error)
	UpdatePlan(ctx context.Context, plan Plan) (Plan, error)

	CreateSubscription(ctx context.Context, subscription Subscription) (Subscription, error)
	GetSubscriptionsByUserId(ctx context.Context, pageableRequest PageableRequest, userId string) (Page[Subscription], error)
	CancelSubscription(ctx context.Context, id string, custId string) error
}

type Server struct {
	router *chi.Mux
	db     Database
}

func NewServer(db Database) *Server {
	s := &Server{
		router: chi.NewRouter(),
		db:     db,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)

	s.router.Get("/hello", s.handleHello)
	s.setupPlanRoutes()
	s.setupSubscriptionRoutes()
}

func (s *Server) handleHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Hello, World!"}`))
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
