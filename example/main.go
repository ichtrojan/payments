package main

import (
	"fmt"
	"github.com/ichtrojan/payment"
	"log"
)

var paystack = payment.Paystack("sk_test_a1a08d9c363c81b67e2686c2bb7931b50aff4f02")

func main() {
	initiateCharge, err := paystack.InitiateCharge("exmaple@domain.com", "chop_life_01")

	if err != nil {
		log.Println(err)
	}

	fmt.Printf("%+v\n", initiateCharge)

	verifyCharge, err := paystack.VerifyCharge("chop_life_01")

	if err != nil {
		log.Println(err)
	}

	fmt.Printf("%+v\n", verifyCharge)

	chargeCard, err := paystack.ChargeCard("AUTH_qeut4h3xfn", "exmaple@domain.com", 9000)

	if err != nil {
		log.Println(err)
	}

	fmt.Printf("%+v\n", chargeCard)

	fetchTransaction, err := paystack.FetchTransaction(01)

	if err != nil {
		log.Println(err)
	}

	fmt.Printf("%+v\n", fetchTransaction)
}
