package main

import (
	"fmt"
	"github.com/IlnurShafikov/wallet/services/health"
	"github.com/IlnurShafikov/wallet/services/wallet"
	"net/http"
)

func main() {
	h := &health.HelloHandler{}
	http.HandleFunc("/health", h.Live)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

	usersWallet := wallet.NewWallet()
	usersWallet.Create("01")
	usersWallet.Add("01", 10)

	balance, err := usersWallet.Get("01")
	if err != nil {
		panic(err)
	}
	fmt.Println(balance)

}
