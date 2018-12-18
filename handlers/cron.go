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
	r.HandleFunc(paths.CronRevivePath, reviveCron)
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

func reviveCron(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	tc, err := scheduler.ThrottledDateQueue.CurrentTaskCount(ctx)
	if err != nil {
		log.Errorf(ctx, "DateQueue.CurrentTaskCount: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if tc > 0 {
		log.Infof(ctx, "ThrottledDateQueue not empty. Skipping reviveCron.")
		return
	}

	l, _ := time.LoadLocation("Asia/Tokyo")
	now := civil.DateOf(time.Now().In(l))

	if err := scheduler.ThrottledDateQueue.EnqueueDateRange(ctx, now.AddDays(-1800), now.AddDays(-30)); err != nil {
		log.Errorf(ctx, "DateQueue.EnqueueDateRange: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
