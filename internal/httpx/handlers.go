package httpx

import (
	"net/http"

	"github.com/moukoublen/goboilerplate/build"
)

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	RespondJSON(r.Context(), w, http.StatusOK, build.GetInfo())
}
