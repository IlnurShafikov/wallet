package health

import (
	"net/http"
)

type HelloHandler struct{}

func (h *HelloHandler) Live(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
