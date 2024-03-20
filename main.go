package main

import (
	"fmt"
	"github.com/IlnurShafikov/wallet/configs"
	"github.com/IlnurShafikov/wallet/services/wallet"
	"log"
	"net/http"
	"strconv"
)

func main() {
	cfg, err := configs.Parse()
	if err != nil {
		panic(err)
	}

	userWallet := wallet.NewWallet()
	handler := wallet.NewHandler(userWallet)

	http.HandleFunc("/create_wallet", handler.CreateWallet)
	fmt.Println("Server work")
	http.HandleFunc("/balance", handler.GetWallet)
	http.HandleFunc("/update", handler.UpdateBalance)

	if err := http.ListenAndServe(":"+strconv.Itoa(cfg.Port), nil); err != nil {
		log.Fatal("Error start server", err)
	}
}
