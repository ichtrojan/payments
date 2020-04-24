package paystack

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexeyco/goozzle"
	"github.com/ichtrojan/thoth"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/url"
	"os"
)

func InitiateCharge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	domain, secretKey := getPaystackConfig()

	email := "trojan@vergly.com"

	endpoint, err := url.Parse(domain + "/transaction/initialize")

	if err != nil {
		logger().Log(err)
	}

	type initiateCardCharge struct {
		Email  string `json:"email"`
		Amount int64  `json:"amount"`
	}

	response, err := goozzle.Post(endpoint).Header("Authorization", fmt.Sprintf("Bearer %s", secretKey)).JSON(initiateCardCharge{
		Email:  email,
		Amount: 5000,
	})

	if err != nil {
		logger().Log(err)
	}

	data := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}{}

	if err := response.JSON(&data); err != nil {
		logger().Log(err)
	}

	if data.Status {
		data := struct {
			Status  bool   `json:"status"`
			Message string `json:"message"`
			Data    struct {
				AccessCode string `json:"access_code"`
				Reference  string `json:"reference"`
			} `json:"data"`
		}{}

		_ = response.JSON(&data)

		if err := json.NewEncoder(w).Encode(data); err != nil {
			logger().Log(err)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger().Log(err)
	}
}

func VerifyCharge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	domain, secretKey := getPaystackConfig()

	reference := r.FormValue("reference")

	endpoint, err := url.Parse(domain + fmt.Sprintf("/transaction/verify/%s", reference))

	if err != nil {
		logger().Log(err)
	}

	response, _ := goozzle.Get(endpoint).Header("Authorization", fmt.Sprintf("Bearer %s", secretKey)).Do()

	data := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}{}

	if err := response.JSON(&data); err != nil {
		logger().Log(err)
	}

	if data.Status {
		data := struct {
			Status  bool   `json:"status"`
			Message string `json:"message"`
			Data    struct {
				Reference     string `json:"reference"`
				Authorization struct {
					AuthorizationCode string `json:"authorization_code"`
					LastFour          string `json:"last4"`
					Brand             string `json:"brand"`
				}
			} `json:"data"`
		}{}

		_ = response.JSON(&data)

		if err := json.NewEncoder(w).Encode(data); err != nil {
			logger().Log(err)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger().Log(err)
	}
}

func AllBanks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	domain, secretKey := getPaystackConfig()

	endpoint, err := url.Parse(domain + "/bank")

	if err != nil {
		logger().Log(err)
	}

	response, _ := goozzle.Get(endpoint).Header("Authorization", fmt.Sprintf("Bearer %s", secretKey)).Do()

	data := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}{}

	if err := response.JSON(&data); err != nil {
		logger().Log(err)
	}

	if data.Status {
		data := struct {
			Status  bool   `json:"status"`
			Message string `json:"message"`
			Data    []struct {
				Name string `json:"name"`
				Code string `json:"code"`
			} `json:"data"`
		}{}

		_ = response.JSON(&data)

		if err := json.NewEncoder(w).Encode(data); err != nil {
			logger().Log(err)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger().Log(err)
	}
}

func GetAccountDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	domain, secretKey := getPaystackConfig()

	accountNumber, bankCode := r.FormValue("account_number"), r.FormValue("bank_code")

	endpoint, err := url.Parse(domain + fmt.Sprintf("/bank/resolve?account_number=%s&bank_code=%s", accountNumber, bankCode))

	if err != nil {
		logger().Log(err)
	}

	response, _ := goozzle.Get(endpoint).Header("Authorization", fmt.Sprintf("Bearer %s", secretKey)).Do()

	data := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}{}

	if err := response.JSON(&data); err != nil {
		logger().Log(err)
	}

	if data.Status {
		data := struct {
			Status  bool   `json:"status"`
			Message string `json:"message"`
			Data    struct {
				AccountNumber string `json:"account_number"`
				AccountName   string `json:"account_name"`
			} `json:"data"`
		}{}

		_ = response.JSON(&data)

		if err := json.NewEncoder(w).Encode(data); err != nil {
			logger().Log(err)
		}

		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger().Log(err)
	}
}

func InitiateBankTransfer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	domain, secretKey := getPaystackConfig()

	endpoint, err := url.Parse(domain + "/transferrecipient")

	accountNumber, bankCode := r.FormValue("account_number"), r.FormValue("bank_code")

	if err != nil {
		logger().Log(err)
	}

	type transferRecipient struct {
		Type          string `json:"type"`
		Description   string `json:"description"`
		AccountNumber string `json:"account_number"`
		BankCode      string `json:"bank_code"`
		Currency      string `json:"currency"`
	}

	response, err := goozzle.Post(endpoint).Header("Authorization", fmt.Sprintf("Bearer %s", secretKey)).JSON(transferRecipient{
		Type:          "nuban",
		Description:   "Vergly Transfer",
		AccountNumber: accountNumber,
		BankCode:      bankCode,
		Currency:      "NGN",
	})

	if err != nil {
		logger().Log(err)
	}

	data := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}{}

	if err := response.JSON(&data); err != nil {
		logger().Log(err)
	}

	if data.Status {
		data := struct {
			Status  bool   `json:"status"`
			Message string `json:"message"`
			Data    struct {
				RecipientCode string `json:"recipient_code"`
			} `json:"data"`
		}{}

		if err := response.JSON(&data); err != nil {
			logger().Log(err)
		}

		endpoint, err := url.Parse(domain + "/transfer")

		if err != nil {
			logger().Log(err)
		}

		amount := r.FormValue("amount")

		amountInInt, _ := strconv.ParseInt(amount, 10, 64)

		type transfer struct {
			Source    string `json:"source"`
			Reason    string `json:"reason"`
			Amount    int64  `json:"amount"`
			Recipient string `json:"recipient"`
			Currency  string `json:"currency"`
		}

		response, err := goozzle.Post(endpoint).Header("Authorization", fmt.Sprintf("Bearer %s", secretKey)).JSON(transfer{
			Source:    "balance",
			Reason:    "Vergly Transfer",
			Amount:    amountInInt,
			Recipient: data.Data.RecipientCode,
			Currency:  "NGN",
		})

		if err != nil {
			logger().Log(err)
		}

		output := struct {
			Status  bool   `json:"status"`
			Message string `json:"message"`
		}{}

		if err := response.JSON(&output); err != nil {
			logger().Log(err)
		}

		if output.Status {
			if err := json.NewEncoder(w).Encode(output); err != nil {
				logger().Log(err)
			}

			return
		}

		if err := json.NewEncoder(w).Encode(output); err != nil {
			logger().Log(err)
		}

		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger().Log(err)
	}
}

func getPaystackConfig() (string, string) {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	domain, exist := os.LookupEnv("PAYSTACK_DOMAIN")

	if !exist {
		logger().Log(errors.New("PAYSTACK_DOMAIN not set in .env"))
		log.Fatal("PAYSTACK_DOMAIN not set in .env")
	}

	secretKey, exist := os.LookupEnv("PAYSTACK_SECRET_KEY")

	if !exist {
		logger().Log(errors.New("PAYSTACK_SECRET_KEY not set in .env"))
		log.Fatal("PAYSTACK_SECRET_KEY not set in .env")
	}

	return domain, secretKey
}

func logger() thoth.Config {
	logger, err := thoth.Init("log")

	if err != nil {
		log.Fatal(err)
	}

	return logger
}
