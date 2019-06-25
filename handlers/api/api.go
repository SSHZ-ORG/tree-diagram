package api

import (
	"encoding/json"
	"net/http"

	"github.com/SSHZ-ORG/tree-diagram/paths"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

const (
	queryPageSize = 10
)

func RegisterAPI(r *mux.Router) {
	r.HandleFunc(paths.APICompareActorsPath, compareActors).Methods("GET", "OPTIONS")

	r.HandleFunc(paths.APIRenderEventPath, renderEvent).Methods("GET", "OPTIONS")
	r.HandleFunc(paths.APIRenderPlacePath, renderPlace).Methods("GET", "OPTIONS")
	r.HandleFunc(paths.APIRenderActorPath, renderActor).Methods("GET", "OPTIONS")

	r.HandleFunc(paths.APIQueryEventsPath, queryEvents).Methods("GET", "OPTIONS")
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
