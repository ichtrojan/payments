package routes

import (
	"github.com/gorilla/mux"
	"github.com/vergly/payment/controllers/paystack"
)

func Api() *mux.Router {
	api := mux.NewRouter()

	card := api.PathPrefix("/card/").Subrouter()
	card.HandleFunc("/initiate", paystack.InitiateCharge)
	card.HandleFunc("/verify", paystack.VerifyCharge).Methods("POST")

	bank := api.PathPrefix("/bank/").Subrouter()
	bank.HandleFunc("/", paystack.AllBanks)
	bank.HandleFunc("/verify", paystack.GetAccountDetails)
	bank.HandleFunc("/initiate", paystack.InitiateBankTransfer)

	return api
}
