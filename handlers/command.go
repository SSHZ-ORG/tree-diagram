package handlers

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/civil"
	"github.com/SSHZ-ORG/tree-diagram/paths"
	"github.com/SSHZ-ORG/tree-diagram/scheduler"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine/v2"
	"google.golang.org/appengine/v2/log"
)

func RegisterCommand(r *httprouter.Router) {
	r.GET(paths.CommandEnqueueDateRangePath, enqueueDateRange)
}

func enqueueDateRange(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
