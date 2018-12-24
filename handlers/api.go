package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SSHZ-ORG/tree-diagram/models"
	"github.com/SSHZ-ORG/tree-diagram/paths"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func RegisterAPI(r *mux.Router) {
	r.HandleFunc(paths.APIGetNoteCountHistoryPath, getNoteCountHistory).Methods("GET", "OPTIONS")
	r.HandleFunc(paths.APIQueryEventsPath, queryEvents)
}

func getNoteCountHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	ctx := appengine.NewContext(r)

	eid := r.FormValue("id")
	if eid == "" {
		http.Error(w, "Missing arg id", http.StatusBadRequest)
		return
	}

	snapshots, err := models.GetSnapshotsForEvent(ctx, eid)
	if err != nil {
		log.Errorf(ctx, "models.GetSnapshotsForEvent: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writeJSON(ctx, w, snapshots)
}

func queryEvents(w http.ResponseWriter, r *http.Request) {
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

	events, err := models.QueryEvents(ctx, placeID, actorIDs, page)
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
	encoded, err := json.Marshal(v)
	if err != nil {
		log.Errorf(ctx, "%v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(encoded)
}
