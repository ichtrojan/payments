package stripeconnect

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ichtrojan/thoth"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v71"
	_ "github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/oauth"
	"log"
	"net/http"
	"os"
)

type CreateOAuthResponse struct {
	Success bool
}

func HandleOauthRedirect(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	stripe.Key = getStripeConfig()

	code := query.Get("code")

	params := &stripe.OAuthTokenParams{
		GrantType: stripe.String("authorization_code"),
		Code:      &code,
	}

	token, err := oauth.New(params)

	if err != nil {
		stripeErr := err.(*stripe.Error)
		if stripeErr.OAuthError == "invalid_grant" {
			logger().Log(err)
			http.Error(w, fmt.Sprintf("Invalid authorization code: %s", code), http.StatusBadRequest)
		} else {
			logger().Log(err)
			http.Error(w, "An unknown error occurred.", http.StatusInternalServerError)
		}
		return
	}

	connectedAccountId := token.StripeUserID

	saveAccountId(connectedAccountId)

	_ = json.NewEncoder(w).Encode(CreateOAuthResponse{
		Success: true,
	})
}

func saveAccountId(id string) {
	log.Println("Connected account ID: " + id)
}

func getStripeConfig() string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey, exist := os.LookupEnv("STRIPE_SECRET_KEY")

	if !exist {
		logger().Log(errors.New("STRIPE_SECRET_KEY not set in .env"))
		log.Fatal("STRIPE_SECRET_KEY not set in .env")
	}

	return secretKey
}

func logger() thoth.Config {
	logger, err := thoth.Init("log")

	if err != nil {
		log.Fatal(err)
	}

	return logger
}
