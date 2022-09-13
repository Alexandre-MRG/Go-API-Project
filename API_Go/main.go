package main

import (
	"blockchain_backend/server"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	http.HandleFunc("/players/create", server.HandleNewPlayer)
	err := http.ListenAndServe(":5000", nil)

	if err != nil {
		println("erreur ListenAndSErve")
		println(err)
	}
}
