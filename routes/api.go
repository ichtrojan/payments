package routes

import (
	"github.com/gorilla/mux"
	"github.com/vergly/payment/controllers/paystack"
	"github.com/vergly/payment/controllers/stripeconnect"
)

func Api() *mux.Router {
	api := mux.NewRouter()

	card := api.PathPrefix("/paystack/card/").Subrouter()
	card.HandleFunc("/initiate", paystack.InitiateCharge)
	card.HandleFunc("/verify", paystack.VerifyCharge).Methods("POST")

	bank := api.PathPrefix("/paystack/bank/").Subrouter()
	bank.HandleFunc("/", paystack.AllBanks)
	bank.HandleFunc("/verify", paystack.GetAccountDetails)
	bank.HandleFunc("/initiate", paystack.InitiateBankTransfer)

	stripe := api.PathPrefix("/stripe/").Subrouter()
	stripe.HandleFunc("/connect/oauth", stripeconnect.HandleOauthRedirect)

	return api
}
