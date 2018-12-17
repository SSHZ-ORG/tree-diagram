package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/civil"
	"github.com/SSHZ-ORG/tree-diagram/crawler"
	"github.com/SSHZ-ORG/tree-diagram/utils"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/admin/enqueue", enqueue)
	r.HandleFunc("/admin/daily", daily)
	r.HandleFunc("/admin/crawl/date", crawlDate)

	http.Handle("/", r)
	appengine.Main()
}

func daily(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	l, _ := time.LoadLocation("Asia/Tokyo")
	now := civil.DateOf(time.Now().In(l))

	if err := enqueueCrawlDate(ctx, now.AddDays(-30), now.AddDays(180)); err != nil {
		log.Errorf(ctx, "enqueueCrawlDate: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func enqueue(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	beginArg := r.FormValue("begin")
	if beginArg == "" {
		http.Error(w, "Missing arg begin", http.StatusBadRequest)
		return
	}
	begin, err := civil.ParseDate(beginArg)
	if err != nil {
		http.Error(w, "Illegal arg begin", http.StatusBadRequest)
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

	if !begin.Before(end) {
		http.Error(w, "begin not before end", http.StatusBadRequest)
		return
	}

	if err := enqueueCrawlDate(ctx, begin, end); err != nil {
		log.Errorf(ctx, "enqueueCrawlDate: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte(fmt.Sprintf("Enqueued %s to %s.", begin.String(), end.String())))
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
		if err := utils.ScheduleNormalQueue(ctx, date, page+1); err != nil {
			log.Errorf(ctx, "enqueueCrawlDate: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	_, _ = w.Write([]byte("OK"))
}

func enqueueCrawlDate(ctx context.Context, begin, end civil.Date) error {
	for cur := begin; cur.Before(end); cur = cur.AddDays(1) {
		if err := utils.ScheduleNormalQueue(ctx, cur, 1); err != nil {
			return err
		}
	}
	return nil
}
