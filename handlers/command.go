package handlers

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/civil"
	"github.com/SSHZ-ORG/tree-diagram/models"
	"github.com/SSHZ-ORG/tree-diagram/paths"
	"github.com/SSHZ-ORG/tree-diagram/scheduler"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func RegisterCommand(r *mux.Router) {
	r.HandleFunc(paths.CommandEnqueueDateRangePath, enqueueDateRange)
	r.HandleFunc(paths.CommandCompressEventSnapshots, compressEventSnapshots)

	r.HandleFunc(paths.CommandBackfill, backfillCES)
}

func enqueueDateRange(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	startArg := r.FormValue("start")
	if startArg == "" {
		http.Error(w, "Missing arg start", http.StatusBadRequest)
		return
	}
	start, err := civil.ParseDate(startArg)
	if err != nil {
		http.Error(w, "Illegal arg start", http.StatusBadRequest)
		return
	}

	endArg := r.FormValue("end")
	if endArg == "" {
		http.Error(w, "Missing arg end", http.StatusBadRequest)
		return
	}
	end, err := civil.ParseDate(endArg)
	if err != nil {
		http.Error(w, "Illegal arg end", http.StatusBadRequest)
		return
	}

	if !start.Before(end) {
		http.Error(w, "start not before end", http.StatusBadRequest)
		return
	}

	if err := scheduler.NormalDateQueue.EnqueueDateRange(ctx, start, end, false); err != nil {
		log.Errorf(ctx, "EventDateQueue.EnqueueDateRange: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte(fmt.Sprintf("Enqueued %s to %s.", start.String(), end.String())))
}

func compressEventSnapshots(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	eid := r.FormValue("id")
	if eid == "" {
		http.Error(w, "Missing arg id", http.StatusBadRequest)
		return
	}

	if err := models.CompressSnapshots(ctx, eid); err != nil {
		log.Errorf(ctx, "models.CompressSnapshots: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte("OK"))
}

func backfillCES(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	ids, err := models.GetSomeEventIDsWithNonCompressedSnapshots(ctx)
	if err != nil {
		log.Errorf(ctx, "models.GetSomeEventIDsWithNonCompressedSnapshots: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, id := range ids {
		scheduler.ScheduleCompressEventSnapshots(ctx, id)
	}

	_, _ = w.Write([]byte("OK"))
}
