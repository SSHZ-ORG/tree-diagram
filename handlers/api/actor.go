package api

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/SSHZ-ORG/tree-diagram/apicache"
	"github.com/SSHZ-ORG/tree-diagram/models"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func prepareRenderActor(ctx context.Context, aids []string) (map[string][]byte, error) {
	m := make(map[string][]byte)

	var missedIDs []string

	fromCache := apicache.GetRenderActor(ctx, aids)
	for _, id := range aids {
		if data, ok := fromCache[id]; ok {
			m[id] = data
		} else {
			missedIDs = append(missedIDs, id)
		}
	}

	responses := make([][]byte, len(missedIDs))
	errs := make([]error, len(missedIDs))
	wg := sync.WaitGroup{}
	wg.Add(len(missedIDs))

	for i, id := range missedIDs {
		go func(i int, id string) {
			defer wg.Done()

			res, err := models.PrepareRenderActorResponse(ctx, id)
			if err != nil {
				errs[i] = err
				return
			}

			encoded, err := json.Marshal(res)
			if err != nil {
				errs[i] = errors.Wrap(err, "Failed to encode RenderActorResponse")
				return
			}

			responses[i] = encoded
		}(i, id)
	}

	wg.Wait()

	toCache := make(map[string][]byte)
	for i, res := range responses {
		if errs[i] != nil {
			return nil, errs[i]
		}

		id := missedIDs[i]
		toCache[id] = res
		m[id] = res
	}

	apicache.PutRenderActor(ctx, toCache)
	return m, nil
}

func renderActor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	ctx := appengine.NewContext(r)

	aid := r.FormValue("id")
	if aid == "" {
		http.Error(w, "Missing arg id", http.StatusBadRequest)
		return
	}

	responses, err := prepareRenderActor(ctx, []string{aid})
	if err != nil {
		log.Errorf(ctx, "prepareRenderActor: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writeEncodedJSON(ctx, w, responses[aid])
}

func compareActors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	ctx := appengine.NewContext(r)

	err := r.ParseForm()
	if err != nil {
		log.Errorf(ctx, "r.ParseForm: %v", err)
		http.Error(w, "Malformed Query", http.StatusBadRequest)
		return
	}

	actorIDs := r.Form["id"]

	responses, err := prepareRenderActor(ctx, actorIDs)
	if err != nil {
		log.Errorf(ctx, "prepareRenderActor: %+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	m := make(map[string]json.RawMessage)
	for k, v := range responses {
		m[k] = json.RawMessage(v)
	}

	encoded, _ := json.Marshal(m)
	writeEncodedJSON(ctx, w, encoded)
}
