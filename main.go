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
	http.ListenAndServe(":8080", nil)

	usersWallet := wallet.NewWallet()
	usersWallet.Create("01")
	usersWallet.Add("01", 10)

	balance, err := usersWallet.Get("01")
	if err != nil {
		panic(err)
	}
	fmt.Println(balance)

}
