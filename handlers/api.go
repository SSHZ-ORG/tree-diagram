package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SSHZ-ORG/tree-diagram/handlers/apicache"
	"github.com/SSHZ-ORG/tree-diagram/models"
	"github.com/SSHZ-ORG/tree-diagram/paths"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

const (
	queryPageSize = 10
)

func RegisterAPI(r *mux.Router) {
	r.HandleFunc(paths.APIRenderEventPath, renderEvent).Methods("GET", "OPTIONS")
	r.HandleFunc(paths.APIRenderPlacePath, renderPlace).Methods("GET", "OPTIONS")

	r.HandleFunc(paths.APIQueryEventsPath, queryEvents).Methods("GET", "OPTIONS")
}

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
		log.Errorf(ctx, "models.PrepareRenderEventResponse: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if encoded, err := encodeJSON(ctx, w, res); err == nil {
		apicache.PutRenderEvent(ctx, eid, encoded)
		writeEncodedJSON(ctx, w, encoded)
	}
}

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
		log.Errorf(ctx, "models.PrepareRenderPlaceResponse: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writeJSON(ctx, w, res)
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

	page := 1
	if arg := r.Form.Get("page"); arg != "" {
		page, err = strconv.Atoi(arg)
		if err != nil {
			http.Error(w, "Illegal arg page", http.StatusBadRequest)
			return
		}
	}

	events, err := models.QueryEvents(ctx, placeID, actorIDs, queryPageSize, (page-1)*queryPageSize)
	if err != nil {
		log.Errorf(ctx, "models.QueryEvents: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var fes []models.FrontendEvent
	for _, e := range events {
		fes = append(fes, e.ToFrontendEvent())
	}
	writeJSON(ctx, w, fes)
}

func writeJSON(ctx context.Context, w http.ResponseWriter, v interface{}) {
	if encoded, err := encodeJSON(ctx, w, v); err == nil {
		writeEncodedJSON(ctx, w, encoded)
	}
}

func encodeJSON(ctx context.Context, w http.ResponseWriter, v interface{}) ([]byte, error) {
	encoded, err := json.Marshal(v)
	if err != nil {
		log.Errorf(ctx, "%v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil, err
	}

	return encoded, nil
}

func writeEncodedJSON(ctx context.Context, w http.ResponseWriter, encoded []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(encoded)
}
