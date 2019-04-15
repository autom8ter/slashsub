package functions

import (
	"github.com/autom8ter/slashsub"
	"net/http"
)

func SlashFunction(w http.ResponseWriter, r *http.Request) {
	s, err := slashsub.New()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.ServeHTTP(w, r)
}
