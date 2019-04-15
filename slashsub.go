package slashsub

import (
	"github.com/autom8ter/slashsub/internal"
	"net/http"
)

func SlashFunction(w http.ResponseWriter, r *http.Request) {
	s, err := internal.New("SlashCmdService", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.ServeHTTP(w, r)
}
