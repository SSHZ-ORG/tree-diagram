package handlers

import (
	"net/http"
	"strconv"

	"cloud.google.com/go/civil"
	"github.com/SSHZ-ORG/tree-diagram/crawler"
	"github.com/SSHZ-ORG/tree-diagram/paths"
	"github.com/SSHZ-ORG/tree-diagram/scheduler"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func RegisterCrawl(r *mux.Router) {
	r.HandleFunc(paths.CrawlDatePath, crawlDate)
}

func crawlDate(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	dateArg := r.FormValue("date")
	if dateArg == "" {
		http.Error(w, "Missing arg date", http.StatusBadRequest)
		return
	}
	date, err := civil.ParseDate(dateArg)
	if err != nil {
		http.Error(w, "Illegal arg date", http.StatusBadRequest)
		return
	}

	pageArg := r.FormValue("page")
	if pageArg == "" {
		http.Error(w, "Missing arg page", http.StatusBadRequest)
		return
	}
	page, err := strconv.Atoi(pageArg)
	if err != nil {
		http.Error(w, "Illegal arg page", http.StatusBadRequest)
		return
	}

	shouldContinue, err := crawler.CrawlDateOnePage(ctx, date, page)
	if err != nil {
		log.Errorf(ctx, "crawler.CrawlDateOnePage: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if shouldContinue {
		if err := scheduler.GetCurrentQueue(r).Schedule(ctx, date, page+1); err != nil {
			log.Errorf(ctx, "DateQueue.Schedule: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	_, _ = w.Write([]byte("OK"))
}
