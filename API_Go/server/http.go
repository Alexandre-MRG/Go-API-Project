package server

import (
	"blockchain_backend/dbPostgres"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
)

const (
	READ_BODY_ERROR_MESSAGE   = "Failed to read body:"
	DESERIALIZE_ERROR_MESSAGE = "Failed to unserialize body"
	SERIALIZE_ERROR_MESSAGE   = "Failed to serialize response"
)

type Wallet struct {
	Wallet_address   string `json:"wallet_address"`
	Currency_code    string `json:"currency_code"`
	Currency_balance int    `json:"currency_balance"`
}

type Player struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Pin_code string `json:"pin_code"`
}

type Response struct {
	User_id        int    `json:"user_id"`
	Wallet_address string `json:"wallet_address"`
}

var regexUsername = regexp.MustCompile(`^[a-z0-9_]{3,100}$`)
var regexPassword = regexp.MustCompile(`^.{6,32}$`)
var regexPin_Code = regexp.MustCompile(`^\d{6}$`)

func HandleNewPlayer(w http.ResponseWriter, r *http.Request) {
	// Try to read body

	body, err := io.ReadAll(r.Body)

	// Check read body is success
	if err != nil {
		LogError(READ_BODY_ERROR_MESSAGE, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	newPlayer := Player{}
	if err := json.Unmarshal(body, &newPlayer); err != nil {
		LogError(DESERIALIZE_ERROR_MESSAGE, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Vérification des paramètres reçus avec des regex et renvoi d'erreur http si nécessaire
	if !regexPassword.MatchString(newPlayer.Password) {
		http.Error(w, "Password must be a string of 6-32 characters.", http.StatusBadRequest)
		return
	}
	if !regexUsername.MatchString(newPlayer.Username) {
		http.Error(w, "Username must be a string of 3-100 characters, and must contain only lowercase (a-z), digits (0-9) or underscore (_)", http.StatusBadRequest)
		return
	}
	if !regexPin_Code.MatchString(newPlayer.Pin_code) {
		http.Error(w, "Pin code must be a string of exactly 6 digits (0-9).", http.StatusBadRequest)
		return
	}

	// Print passed player
	//fmt.Printf("newPlayer: %s", newPlayer.Username)

	db, _ := dbPostgres.InitializeDB()
	defer db.CloseDB()

	id := db.AddPlayer(newPlayer.Username, newPlayer.Password, newPlayer.Pin_code)
	wallet := addBlockchainWallet(newPlayer.Pin_code)
	walletAddress := db.AddWallet(id, wallet.Wallet_address, wallet.Currency_code)

	// Create response
	res := Response{User_id: id, Wallet_address: walletAddress}
	// Serialize response
	jsonResponse, err := json.Marshal(res)
	// Check serialize success
	if err != nil {
		LogError(SERIALIZE_ERROR_MESSAGE, err)
		println("L93")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func addBlockchainWallet(pin_code string) Wallet { // Méthode qui effectue une requête POST sur le serveur hébergeant la blockchain pour récupérer un wallet
	type Request struct {
		Blockchain string `json:"blockchain"`
		Pin_code   string `json:"pin_code"`
	}

	newRequest := Request{
		Blockchain: "ethereum",
		Pin_code:   pin_code,
	}

	jsonRequest, err := json.Marshal(newRequest)
	if err != nil {
		LogError(SERIALIZE_ERROR_MESSAGE, err)
		println("L115")
	}

	resp, err := http.Post("http://localhost:5001/wallets/create", "application/json", bytes.NewBuffer(jsonRequest))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	newWallet := Wallet{}
	if err := json.Unmarshal(body, &newWallet); err != nil {
		LogError(DESERIALIZE_ERROR_MESSAGE, err)
	}

	return newWallet
}

func LogError(errorMessage string, err error) {
	log.Println(errorMessage, err)
}
