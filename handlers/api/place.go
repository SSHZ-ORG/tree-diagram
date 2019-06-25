package api

import (
	"net/http"

	"github.com/SSHZ-ORG/tree-diagram/models"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func renderPlace(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	ctx := appengine.NewContext(r)

	pid := r.FormValue("id")
	if pid == "" {
		http.Error(w, "Missing arg id", http.StatusBadRequest)
		return
	}

	res, err := models.PrepareRenderPlaceResponse(ctx, pid)
	if err != nil {
		log.Errorf(ctx, "models.PrepareRenderPlaceResponse: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writeJSON(ctx, w, res)
}
