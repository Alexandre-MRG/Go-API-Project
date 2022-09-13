package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/jaswdr/faker"
)

const (
	READ_BODY_ERROR_MESSAGE   = "Failed to read body:"
	DESERIALIZE_ERROR_MESSAGE = "Failed to unserialize body"
	SERIALIZE_ERROR_MESSAGE   = "Failed to serialize response"
)

type Wallet struct {
	Blockchain string `json:"blockchain"`
	Pin_code   string `json:"pin_code"`
}

type Response struct {
	Wallet_address   string `json:"wallet_address"`
	Currency_code    string `json:"currency_code"`
	Currency_balance int    `json:"currency_balance"`
}

func HandleNewWallet(w http.ResponseWriter, r *http.Request) {
	// Try to read body
	body, err := io.ReadAll(r.Body)

	// Check read body is success
	if err != nil {
		LogError(READ_BODY_ERROR_MESSAGE, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	newWallet := Wallet{}
	if err := json.Unmarshal(body, &newWallet); err != nil {
		LogError(DESERIALIZE_ERROR_MESSAGE, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Print passed tweet
	fmt.Printf("newWallet: %s", newWallet.Blockchain)

	address := faker.Crypto.EtheriumAddress(faker.New().Crypto())
	// Create response
	res := Response{
		Wallet_address:   address,
		Currency_code:    "ETH",
		Currency_balance: 0,
	}
	// Serialize response
	jsonResponse, err := json.Marshal(res)
	// Check serialize success
	if err != nil {
		LogError(SERIALIZE_ERROR_MESSAGE, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func LogError(errorMessage string, err error) {
	log.Println(errorMessage, err)
}
