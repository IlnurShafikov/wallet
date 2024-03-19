package main

import (
	"fmt"
	"github.com/IlnurShafikov/wallet/services/wallet"
	"log"
	"net/http"
)

func main() {
	userWallet := wallet.NewWallet()
	handler := wallet.NewHandler(userWallet)

	http.HandleFunc("/create_wallet", handler.CreateWallet)
	fmt.Println("Server work")
	http.HandleFunc("/balance", handler.GetWallet)
	http.HandleFunc("/update", handler.UpdateBalance)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error start server", err)
	}
}
