package health

import (
	"fmt"
	"net/http"
)

type HelloHandler struct{}

func (h *HelloHandler) Live(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprint(w, "Ilnur")
	if err != nil {
		fmt.Println(err)
	}
}
