package main

import (
	"fmt"

	"github.com/IlnurShafikov/wallet/services/wallet"
)


func main() {
	usersWallet := wallet.NewWallet()
	usersWallet.Create("01")
	usersWallet.Add("01", 10)
	  
	balance, err := usersWallet.Get("01")
	if err != nil {
		panic(err)
	}

	

	fmt.Println(balance)
}