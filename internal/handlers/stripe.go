package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/curtisbraxdale/taday/internal/auth"
	"github.com/curtisbraxdale/taday/internal/database"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/subscription"
)

func (cfg *ApiConfig) CreateCheckoutSession(w http.ResponseWriter, req *http.Request) {
	accessCookie, err := req.Cookie("access_token")
	if err != nil {
		log.Printf("Access token not found in cookies: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	accessToken := accessCookie.Value

	userID, err := auth.ValidateAccessToken(accessToken, cfg.Secret)
	if err != nil {
		log.Printf("Access token invalid: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userStripeID, err := cfg.Queries.GetStripeID(req.Context(), userID)
	if err != nil {
		log.Printf("Error getting StripeID: %s", err)
		w.WriteHeader(500)
		return
	}

	domain := "https://taday.io"

	if !userStripeID.Valid {
		userEmail, err := cfg.Queries.GetEmail(req.Context(), userID)
		if err != nil {
			log.Printf("Error getting email: %s", err)
			w.WriteHeader(500)
			return
		}
		stripeCustomer, err := customer.New(&stripe.CustomerParams{
			Email: stripe.String(userEmail),
		})
		if err != nil {
			log.Printf("Error creating Stripe user: %s", err)
			w.WriteHeader(500)
			return
		}
		userStripeID = sql.NullString{String: stripeCustomer.ID, Valid: true}
		err = cfg.Queries.UpdateStripeCustomerID(req.Context(), database.UpdateStripeCustomerIDParams{ID: userID, StripeCustomerID: userStripeID})
		if err != nil {
			log.Printf("Error updating Stripe customer ID: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String("price_1RgrzhP4x5qmnWDOq8VmMjTy"),
				Quantity: stripe.Int64(1),
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),

		SuccessURL: stripe.String(domain + "/success"),
		CancelURL:  stripe.String(domain + "/cancel"),

		AutomaticTax: &stripe.CheckoutSessionAutomaticTaxParams{
			Enabled: stripe.Bool(true),
		},

		Customer: stripe.String(userStripeID.String),
	}

	s, err := session.New(params)
	if err != nil {
		log.Printf("Stripe session creation failed: %v", err)
		http.Error(w, "Stripe error", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"url": s.URL})
}

func (cfg *ApiConfig) CancelSub(w http.ResponseWriter, req *http.Request) {
	accessCookie, err := req.Cookie("access_token")
	if err != nil {
		log.Printf("Access token not found in cookies: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	accessToken := accessCookie.Value

	userID, err := auth.ValidateAccessToken(accessToken, cfg.Secret)
	if err != nil {
		log.Printf("Access token invalid: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dbSubscription, err := cfg.Queries.GetActiveSubscriptionByUserID(req.Context(), userID)
	if err != nil {
		http.Error(w, "No active subscription found", http.StatusNotFound)
		return
	}

	// Call Stripe to update subscription
	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}
	_, err = subscription.Update(dbSubscription.StripeSubscriptionID, params)
	if err != nil {
		log.Printf("Error updating subscription in Stripe: %v", err)
		http.Error(w, "Failed to cancel subscription", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
