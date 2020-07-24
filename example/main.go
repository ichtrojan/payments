package main

import (
	"fmt"
	"github.com/ichtrojan/payment"
	"log"
)

var paystack = payment.Paystack("sk_test_000000000000000000000000000000")

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
}
