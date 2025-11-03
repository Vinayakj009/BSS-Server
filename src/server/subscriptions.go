package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func (d *Server) setupSubscriptionRoutes() {
	d.router.Post("/customers/{customer_id}/subscribe", d.handleCreateSubscription)
	d.router.Get("/customers/{customer_id}/subscriptions", d.handleGetSubscriptionsByUserId)
	d.router.Post("/customers/{customer_id}/unsubscribe", d.handleCancelSubscription)
}

func (d *Server) handleCreateSubscription(w http.ResponseWriter, r *http.Request) {
	customerId := r.PathValue("customer_id")
	customerUUID, err := uuid.Parse(customerId)
	if err != nil {
		http.Error(w, "invalid customer_id", http.StatusBadRequest)
		return
	}
	var subscription Subscription
	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	_, err = d.db.GetPlan(r.Context(), subscription.PlanID.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	subscription.CustomerID = customerUUID
	createdSubscription, err := d.db.CreateSubscription(r.Context(), subscription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdSubscription)
}

func (d *Server) handleGetSubscriptionsByUserId(w http.ResponseWriter, r *http.Request) {
	customerId := r.PathValue("customer_id")
	customerUUID, err := uuid.Parse(customerId)
	if err != nil {
		http.Error(w, "invalid customer_id", http.StatusBadRequest)
		return
	}

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
	subscriptionsPage, err := d.db.GetSubscriptionsByUserId(r.Context(), pageableRequest, customerUUID.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(subscriptionsPage)

}

func (d *Server) handleCancelSubscription(w http.ResponseWriter, r *http.Request) {
	customerId := r.PathValue("customer_id")
	if customerId == "" {
		http.Error(w, "invalid customer_id", http.StatusBadRequest)
		return
	}
	subscriptionId := r.URL.Query().Get("subscription_id")
	if subscriptionId == "" {
		http.Error(w, "subscription_id is required", http.StatusBadRequest)
		return
	}
	err := d.db.CancelSubscription(r.Context(), subscriptionId, customerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
