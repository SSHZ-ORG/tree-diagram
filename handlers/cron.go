package handlers

import (
	"fmt"
	"net/http"
	"net/url"

	"cloud.google.com/go/civil"
	"github.com/SSHZ-ORG/tree-diagram/models"
	"github.com/SSHZ-ORG/tree-diagram/paths"
	"github.com/SSHZ-ORG/tree-diagram/scheduler"
	"github.com/SSHZ-ORG/tree-diagram/utils"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/datastore/v1"
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
	r.GET(paths.CronExportPath, exportCron)
	r.GET(paths.CronOneOff, oneOffCron)
}

func dailyCron(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(r)

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

	if err := scheduler.ScheduleCrawlActorPage(ctx, 1); err != nil {
		log.Errorf(ctx, "scheduler.ScheduleCrawlActorPage: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func exportCron(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(r)

	client, err := google.DefaultClient(ctx, datastore.DatastoreScope)
	if err != nil {
		log.Errorf(ctx, "google.DefaultClient: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	service, err := datastore.New(client)
	if err != nil {
		log.Errorf(ctx, "datastore.New: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	appID := appengine.AppID(ctx)
	bucketName := appID + ".appspot.com"

	resp, err := service.Projects.Export(appID, &datastore.GoogleDatastoreAdminV1ExportEntitiesRequest{
		OutputUrlPrefix: fmt.Sprintf("gs://%s/td-datastore/%s/", bucketName, utils.JSTToday().String()),
		EntityFilter:    &datastore.GoogleDatastoreAdminV1EntityFilter{Kinds: models.KindsToExport},
	}).Do()
	if err != nil {
		log.Errorf(ctx, "service.Projects.Export: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if resp.HTTPStatusCode != http.StatusOK {
		log.Errorf(ctx, "service.Projects.Export returned error: %v", resp.Error.Message)
		http.Error(w, resp.Error.Message, resp.HTTPStatusCode)
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
