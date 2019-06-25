package api

import (
	"net/http"

	"github.com/SSHZ-ORG/tree-diagram/apicache"
	"github.com/SSHZ-ORG/tree-diagram/models"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func renderActor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	ctx := appengine.NewContext(r)

	aid := r.FormValue("id")
	if aid == "" {
		http.Error(w, "Missing arg id", http.StatusBadRequest)
		return
	}

	if fromCache := apicache.GetRenderActor(ctx, aid); fromCache != nil {
		writeEncodedJSON(ctx, w, fromCache)
		return
	}

	res, err := models.PrepareRenderActorResponse(ctx, aid)
	if err != nil {
		log.Errorf(ctx, "models.PrepareRenderActorResponse: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if encoded, err := encodeJSON(ctx, w, res); err == nil {
		apicache.PutRenderActor(ctx, aid, encoded)
		writeEncodedJSON(ctx, w, encoded)
	}
}
