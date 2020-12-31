package api

import (
	"net/http"

	"github.com/SSHZ-ORG/tree-diagram/models"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func renderPlace(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
