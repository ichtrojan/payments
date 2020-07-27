package payment

import (
	"errors"
	"fmt"
	"github.com/alexeyco/goozzle"
	"net/url"
)

const paystackDomain = "https://api.paystack.co"

type paystackConfig struct {
	secretKey string
}

type initiateCharge struct {
	AccessCode  string
	Reference   string
	CheckoutURL string
}

type verifyCharge struct {
	Reference         string
	AuthorizationCode string
	FirstSix          string
	LastFour          string
	Brand             string
	Month             string
	Year              string
	Bank              string
}

type chargeCard struct {
	Message   string
	Amount    int64
	Reference string
}

type transaction struct {
	Status    string
	Message   string
	Id		  int64	
	Amount    int64
	Reference string      
	Bank      string
	CardType  string
}

func Paystack(secretKey string) paystackConfig {
	return paystackConfig{
		secretKey: secretKey,
	}
}

func (config paystackConfig) InitiateCharge(email string, reference string) (initiateCharge, error) {
	endpoint, err := url.Parse(paystackDomain + "/transaction/initialize")

	if err != nil {
		return initiateCharge{}, err
	}

	response, err := goozzle.Post(endpoint).Header("Authorization", fmt.Sprintf("Bearer %s", config.secretKey)).JSON(struct {
		Email     string `json:"email"`
		Amount    int64  `json:"amount"`
		Reference string `json:"reference"`
	}{
		Email:     email,
		Amount:    5000,
		Reference: reference,
	})

	if err != nil {
		return initiateCharge{}, err
	}

	data := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}{}

	if err := response.JSON(&data); err != nil {
		return initiateCharge{}, err
	}

	if data.Status && response.Status() == 200 {
		data := struct {
			Data struct {
				AccessCode string `json:"access_code"`
				Reference  string `json:"reference"`
			} `json:"data"`
		}{}

		_ = response.JSON(&data)

		return initiateCharge{
			AccessCode:  data.Data.AccessCode,
			Reference:   data.Data.Reference,
			CheckoutURL: "https://checkout.paystack.com/" + data.Data.AccessCode,
		}, nil
	}

	return initiateCharge{}, errors.New(data.Message)
}

func (config paystackConfig) VerifyCharge(reference string) (verifyCharge, error) {
	endpoint, err := url.Parse(paystackDomain + fmt.Sprintf("/transaction/verify/%s", reference))

	if err != nil {
		return verifyCharge{}, err
	}

	response, _ := goozzle.Get(endpoint).Header("Authorization", fmt.Sprintf("Bearer %s", config.secretKey)).Do()

	data := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}{}

	if err := response.JSON(&data); err != nil {
		return verifyCharge{}, err
	}

	if data.Status && response.Status() == 200 {
		data := struct {
			Status  bool   `json:"status"`
			Message string `json:"message"`
			Data    struct {
				Reference     string `json:"reference"`
				Authorization struct {
					AuthorizationCode string `json:"authorization_code"`
					FirstSix          string `json:"bin"`
					LastFour          string `json:"last4"`
					Brand             string `json:"brand"`
					Month             string `json:"exp_month"`
					Year              string `json:"exp_year"`
					Bank              string `json:"bank"`
				}
			} `json:"data"`
		}{}

		_ = response.JSON(&data)

		return verifyCharge{
			AuthorizationCode: data.Data.Authorization.AuthorizationCode,
			FirstSix:          data.Data.Authorization.FirstSix,
			LastFour:          data.Data.Authorization.LastFour,
			Brand:             data.Data.Authorization.Brand,
			Month:             data.Data.Authorization.Month,
			Year:              data.Data.Authorization.Year,
			Bank:              data.Data.Authorization.Bank,
		}, nil
	}

	return verifyCharge{}, errors.New(data.Message)
}

func (config paystackConfig) ChargeCard(authorization string, email string, amount int64) (chargeCard, error) {
	endpoint, err := url.Parse(paystackDomain + "/transaction/charge_authorization")

	if err != nil {
		return chargeCard{}, err
	}

	response, err := goozzle.Post(endpoint).Header("Authorization", fmt.Sprintf("Bearer %s", config.secretKey)).JSON(struct {
		Email             string `json:"email"`
		Amount            int64  `json:"amount"`
		AuthorizationCode string `json:"authorization_code"`
	}{
		Email:             email,
		Amount:            amount,
		AuthorizationCode: authorization,
	})

	if err != nil {
		return chargeCard{}, err
	}

	data := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}{}

	if err := response.JSON(&data); err != nil {
		return chargeCard{}, err
	}

	if data.Status && response.Status() == 200 {
		data := struct {
			Data struct {
				Id              int64  `json:"id"`
				Status          string `json:"status"`
				Amount          int64  `json:"amount"`
				Reference       string `json:"reference"`
				GatewayResponse string `json:"gateway_response"`
			} `json:"data"`
		}{}

		_ = response.JSON(&data)

		if data.Data.Status == "failed" {
			return chargeCard{}, errors.New(data.Data.GatewayResponse)
		}

		return chargeCard{
			Message:   data.Data.GatewayResponse,
			Amount:    data.Data.Amount,
			Reference: data.Data.Reference,
		}, nil
	}

	return chargeCard{}, errors.New(data.Message)
}


func (config paystackConfig) FetchTransaction(transactionId int64) (transaction, error) {
	endpoint, err := url.Parse(paystackDomain + fmt.Sprintf("/transaction/%d", transactionId))

	if err != nil {
		return transaction{}, err
	}

	response, _ := goozzle.Get(endpoint).Header("Authorization", fmt.Sprintf("Bearer %s", config.secretKey)).Do()

	data := struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}{}

	if err := response.JSON(&data); err != nil {
		return transaction{}, err
	}

	if data.Status && response.Status() == 200 {
		data := struct {
			Data struct {
				Id              int64  `json:"id"`
				Status          string `json:"status"`
				Amount          int64  `json:"amount"`
				Reference       string `json:"reference"`
				GatewayResponse string `json:"gateway_response"`
				Authorization struct {
					CardType          string `json:"brand"`
					Bank              string `json:"bank"`
				}
			} `json:"data"`
		}{}


		_ = response.JSON(&data)


		return transaction{
			CardType:            data.Data.Authorization.CardType,
			Bank:                data.Data.Authorization.Bank,
			Message:   			 data.Data.GatewayResponse,
			Reference:           data.Data.Reference,
			Status:              data.Data.Status,
			Amount:              data.Data.Amount,
			Id:                  data.Data.Id,
		}, nil
	}

	return transaction{}, errors.New(data.Message)
}