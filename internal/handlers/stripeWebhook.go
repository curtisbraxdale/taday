package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/curtisbraxdale/taday/internal/database"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

func (cfg *ApiConfig) StripeWebhookHandler(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading webhook body: %v", err)
		http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
		return
	}

	// Use your Stripe webhook secret (get from Dashboard)
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"), endpointSecret)
	if err != nil {
		log.Printf("⚠️ Webhook signature verification failed: %v", err)
		http.Error(w, "Unauthorized", http.StatusBadRequest)
		return
	}
	switch event.Type {
	case "customer.subscription.created":
		var sub stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
			log.Printf("Error parsing webhook subscription: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		user, err := cfg.Queries.GetUserByStripeID(r.Context(), sql.NullString{Valid: true, String: sub.Customer.ID})
		_, err = cfg.Queries.CreateSubscription(r.Context(), database.CreateSubscriptionParams{
			UserID:               user.ID,
			StripeCustomerID:     sub.Customer.ID,
			StripeSubscriptionID: sub.ID,
			Plan:                 "pro",
			Status:               string(sub.Status),
			CurrentPeriodStart:   time.Unix(sub.Items.Data[0].CurrentPeriodStart, 0),
			CurrentPeriodEnd:     time.Unix(sub.Items.Data[0].CurrentPeriodEnd, 0),
			CancelAtPeriodEnd:    sub.CancelAtPeriodEnd,
			CanceledAt:           TimeOrNil(sub.CanceledAt),
			TrialStart:           TimeOrNil(sub.TrialStart),
			TrialEnd:             TimeOrNil(sub.TrialEnd),
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		})
		if err != nil {
			log.Printf("Error creating subscription: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "customer.subscription.updated", "customer.subscription.deleted":
		var sub stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
			log.Printf("Error parsing webhook subscription: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		user, err := cfg.Queries.GetUserByStripeID(r.Context(), sql.NullString{Valid: true, String: sub.Customer.ID})
		_, err = cfg.Queries.UpdateSubscription(r.Context(), database.UpdateSubscriptionParams{
			Status:             string(sub.Status),
			CurrentPeriodStart: time.Unix(sub.Items.Data[0].CurrentPeriodStart, 0),
			CurrentPeriodEnd:   time.Unix(sub.Items.Data[0].CurrentPeriodEnd, 0),
			CancelAtPeriodEnd:  sub.CancelAtPeriodEnd,
			CanceledAt:         TimeOrNil(sub.CanceledAt),
			TrialStart:         TimeOrNil(sub.TrialStart),
			TrialEnd:           TimeOrNil(sub.TrialEnd),
			UserID:             user.ID,
		})
		if err != nil {
			log.Printf("Error updating subscription: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func TimeOrNil(ts int64) sql.NullTime {
	if ts == 0 {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: time.Unix(ts, 0), Valid: true}
}
