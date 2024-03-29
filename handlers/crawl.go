package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"cloud.google.com/go/civil"
	"github.com/SSHZ-ORG/tree-diagram/crawler"
	"github.com/SSHZ-ORG/tree-diagram/paths"
	"github.com/SSHZ-ORG/tree-diagram/scheduler"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine/v2"
	"google.golang.org/appengine/v2/log"
)

func RegisterCrawl(r *httprouter.Router) {
	r.POST(paths.CrawlDatePath, crawlDate)
	r.POST(paths.CrawlActorPath, crawlActor)
}

func crawlDate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
		log.Errorf(ctx, "crawler.CrawlDateOnePage: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if shouldContinue {
		if err := scheduler.GetCurrentEventDateQueue(r).Schedule(ctx, date, page+1); err != nil {
			log.Errorf(ctx, "EventDateQueue.Schedule: %+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	_, _ = w.Write([]byte("OK"))
}

func crawlActor(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(r)

	offsetArg := r.FormValue("offset")
	if offsetArg == "" {
		http.Error(w, "Missing arg offset", http.StatusBadRequest)
		return
	}
	offset, err := strconv.Atoi(offsetArg)
	if err != nil {
		http.Error(w, "Illegal arg offset", http.StatusBadRequest)
		return
	}

	nextOffset, err := crawler.CrawlActorOnePage(ctx, offset)
	if err != nil {
		log.Errorf(ctx, "crawler.CrawlActorOnePage: %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if nextOffset > 0 {
		if err := scheduler.ScheduleCrawlActorPage(ctx, nextOffset); err != nil {
			log.Errorf(ctx, "ScheduleCrawlActorPage: %+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	_, _ = w.Write([]byte(fmt.Sprintf("Next offset: %d", nextOffset)))
}
