package api

import (
	"net/http"
	"strconv"

	"github.com/SSHZ-ORG/tree-diagram/apicache"
	"github.com/SSHZ-ORG/tree-diagram/models"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func renderEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	ctx := appengine.NewContext(r)

	eid := r.FormValue("id")
	if eid == "" {
		http.Error(w, "Missing arg id", http.StatusBadRequest)
		return
	}

	if fromCache := apicache.GetRenderEvent(ctx, eid); fromCache != nil {
		writeEncodedJSON(ctx, w, fromCache)
		return
	}

	res, err := models.PrepareRenderEventResponse(ctx, eid)
	if err != nil {
		log.Errorf(ctx, "models.PrepareRenderEventResponse: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if encoded, err := encodeJSON(ctx, w, res); err == nil {
		apicache.PutRenderEvent(ctx, eid, encoded)
		writeEncodedJSON(ctx, w, encoded)
	}
}

func queryEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	ctx := appengine.NewContext(r)

	err := r.ParseForm()
	if err != nil {
		log.Errorf(ctx, "r.ParseForm: %v", err)
		http.Error(w, "Malformed Query", http.StatusBadRequest)
		return
	}

	placeID := r.Form.Get("place")
	actorIDs := r.Form["actor"]

	offset := 0
	if arg := r.Form.Get("offset"); arg != "" {
		offset, err = strconv.Atoi(arg)
		if err != nil {
			http.Error(w, "Illegal arg offset", http.StatusBadRequest)
			return
		}
	}

	if arg := r.Form.Get("page"); arg != "" && offset == 0 {
		page, err := strconv.Atoi(arg)
		if err != nil {
			http.Error(w, "Illegal arg page", http.StatusBadRequest)
			return
		}
		offset = (page - 1) * queryPageSize
	}

	events, err := models.QueryEvents(ctx, placeID, actorIDs, queryPageSize, offset)
	if err != nil {
		log.Errorf(ctx, "models.QueryEvents: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fes := make([]models.FrontendEvent, 0)
	for _, e := range events {
		fes = append(fes, e.ToFrontendEvent())
	}
	writeJSON(ctx, w, fes)
}
