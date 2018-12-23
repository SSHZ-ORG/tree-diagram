package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/SSHZ-ORG/tree-diagram/models"
	"github.com/SSHZ-ORG/tree-diagram/paths"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func RegisterAPI(r *mux.Router) {
	r.HandleFunc(paths.APIGetNoteCountHistoryPath, getNoteCountHistory).Methods("GET", "OPTIONS")
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
