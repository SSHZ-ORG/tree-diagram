package handlers

import (
	"net/http"
	"time"

	"cloud.google.com/go/civil"
	"github.com/SSHZ-ORG/tree-diagram/paths"
	"github.com/SSHZ-ORG/tree-diagram/scheduler"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func RegisterCron(r *mux.Router) {
	r.HandleFunc(paths.CronDailyPath, dailyCron)
}

func dailyCron(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	l, _ := time.LoadLocation("Asia/Tokyo")
	now := civil.DateOf(time.Now().In(l))

	if err := scheduler.NormalDateQueue.EnqueueDateRange(ctx, now.AddDays(-30), now.AddDays(180)); err != nil {
		log.Errorf(ctx, "DateQueue.EnqueueDateRange: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
