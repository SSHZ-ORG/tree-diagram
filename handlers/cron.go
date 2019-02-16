package handlers

import (
	"net/http"

	"github.com/SSHZ-ORG/tree-diagram/models"
	"github.com/SSHZ-ORG/tree-diagram/paths"
	"github.com/SSHZ-ORG/tree-diagram/scheduler"
	"github.com/SSHZ-ORG/tree-diagram/utils"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/datastore/v1"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func RegisterCron(r *mux.Router) {
	r.HandleFunc(paths.CronDailyPath, dailyCron)
	r.HandleFunc(paths.CronRevivePath, reviveCron)
	r.HandleFunc(paths.CronCleanupPath, cleanupCron)
	r.HandleFunc(paths.CronExportPath, exportCron)
}

func dailyCron(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	today := utils.JSTToday()

	if err := scheduler.NormalDateQueue.EnqueueDateRange(ctx, today.AddDays(-30), today.AddDays(360)); err != nil {
		log.Errorf(ctx, "DateQueue.EnqueueDateRange: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func reviveCron(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	tc, err := scheduler.ThrottledDateQueue.CurrentTaskCount(ctx)
	if err != nil {
		log.Errorf(ctx, "DateQueue.CurrentTaskCount: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if tc > 0 {
		log.Infof(ctx, "ThrottledDateQueue not empty. Skipping reviveCron.")
		return
	}

	today := utils.JSTToday()

	if err := scheduler.ThrottledDateQueue.EnqueueDateRange(ctx, today.AddDays(-1800), today.AddDays(-30)); err != nil {
		log.Errorf(ctx, "DateQueue.EnqueueDateRange: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func cleanupCron(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if err := models.CleanupFinishedEvents(ctx, utils.JSTToday()); err != nil {
		log.Errorf(ctx, "models.CleanupFinishedEvents: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func exportCron(w http.ResponseWriter, r *http.Request) {
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
		OutputUrlPrefix: "gs://" + bucketName + "/td-datastore/latest/",
		EntityFilter:    &datastore.GoogleDatastoreAdminV1EntityFilter{Kinds: models.AllKinds},
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
