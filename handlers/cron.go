package handlers

import (
	"net/http"
	"net/url"

	"cloud.google.com/go/civil"
	"github.com/SSHZ-ORG/tree-diagram/models"
	"github.com/SSHZ-ORG/tree-diagram/paths"
	"github.com/SSHZ-ORG/tree-diagram/scheduler"
	"github.com/SSHZ-ORG/tree-diagram/utils"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine/v2"
	"google.golang.org/appengine/v2/log"
)

const (
	dailyEndDate       = 360
	dailyToReviveDate  = -30
	reviveToUndeadDate = -1800

	undeadStartDate = "2000-01-01"
)

func RegisterCron(r *httprouter.Router) {
	r.GET(paths.CronDailyPath, dailyCron)
	r.GET(paths.CronRevivePath, reviveCron)
	r.GET(paths.CronUndeadPath, undeadCron)
	r.GET(paths.CronCleanupPath, cleanupCron)
	r.GET(paths.CronDailyActorPath, dailyActorCron)
	r.GET(paths.CronOneOff, oneOffCron)
}

func dailyCron(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(r)

	tc, err := scheduler.NormalDateQueue.CurrentTaskCount(ctx)
	if err != nil {
		log.Errorf(ctx, "EventDateQueue.CurrentTaskCount: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if tc > 0 {
		log.Warningf(ctx, "NormalDateQueue not empty. Skipping dailyCron.")
		return
	}

	today := utils.JSTToday()

	if err := scheduler.NormalDateQueue.EnqueueDateRange(ctx, today.AddDays(dailyToReviveDate), today.AddDays(dailyEndDate), false); err != nil {
		log.Errorf(ctx, "EventDateQueue.EnqueueDateRange: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func reviveCron(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(r)

	tc, err := scheduler.ThrottledDateQueue.CurrentTaskCount(ctx)
	if err != nil {
		log.Errorf(ctx, "EventDateQueue.CurrentTaskCount: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if tc > 0 {
		log.Infof(ctx, "ThrottledDateQueue not empty. Skipping reviveCron.")
		return
	}

	today := utils.JSTToday()

	if err := scheduler.ThrottledDateQueue.EnqueueDateRange(ctx, today.AddDays(reviveToUndeadDate), today.AddDays(dailyToReviveDate), false); err != nil {
		log.Errorf(ctx, "EventDateQueue.EnqueueDateRange: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func undeadCron(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(r)

	tc, err := scheduler.DeadSlowDateQueue.CurrentTaskCount(ctx)
	if err != nil {
		log.Errorf(ctx, "EventDateQueue.CurrentTaskCount: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if tc > 0 {
		log.Infof(ctx, "DeadSlowDateQueue not empty. Skipping undeadCron.")
		return
	}

	today := utils.JSTToday()

	startDate, _ := civil.ParseDate(undeadStartDate)
	if err := scheduler.DeadSlowDateQueue.EnqueueDateRange(ctx, startDate, today.AddDays(reviveToUndeadDate), true); err != nil {
		log.Errorf(ctx, "EventDateQueue.EnqueueDateRange: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func cleanupCron(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(r)

	if err := models.CleanupFinishedEvents(ctx, utils.JSTToday()); err != nil {
		log.Errorf(ctx, "models.CleanupFinishedEvents: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func dailyActorCron(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(r)

	ht, err := scheduler.ActorQueueHasTask(ctx)
	if err != nil {
		log.Errorf(ctx, "scheduler.ActorQueueHasTask: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if ht {
		log.Warningf(ctx, "ActorQueue has task. Skipping dailyActorCron.")
		return
	}

	if err := scheduler.ScheduleCrawlActorPage(ctx, 1); err != nil {
		log.Errorf(ctx, "scheduler.ScheduleCrawlActorPage: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func oneOffCron(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(r)

	cursor := r.FormValue("cursor")
	log.Debugf(ctx, "Received cursor %s", cursor)

	newCursor, err := models.OneoffBackfillModelVersion(ctx, cursor)
	if err != nil {
		log.Errorf(ctx, "models.OneoffBackfillModelVersion: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if newCursor != "" {
		log.Debugf(ctx, "Scheduling cursor %s", newCursor)
		err = scheduler.ScheduleOneOff(ctx, url.Values{
			"cursor": []string{newCursor},
		})
		if err != nil {
			log.Errorf(ctx, "scheduler.ScheduleOneOff: %+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
