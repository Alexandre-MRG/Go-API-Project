package main

import (
	"mockblockchain/server"
	"net/http"
)

func main() {
	http.HandleFunc("/wallets/create", server.HandleNewWallet)
	http.ListenAndServe(":5001", nil)
}
