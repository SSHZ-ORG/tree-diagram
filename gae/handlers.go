package main

import (
	"net/http"
	"net/url"
	"time"

	"cloud.google.com/go/civil"
	"github.com/SSHZ-ORG/tree-diagram/crawler"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
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

	enqueueCrawlDate(ctx, now.AddDays(-30), now.AddDays(180))
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
	enqueueCrawlDate(ctx, begin, end)

	w.Write([]byte("OK"))
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

	if err := crawler.CrawlDateOnePage(ctx, date); err != nil {
		log.Errorf(ctx, "crawler.CrawlDateOnePage: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.Write([]byte("OK"))
	}
}

func enqueueCrawlDate(ctx context.Context, begin, end civil.Date) {
	for cur := begin; cur.Before(end); cur = cur.AddDays(1) {
		t := taskqueue.NewPOSTTask("/admin/crawl/date", url.Values{
			"date": []string{cur.String()},
		})

		_, err := taskqueue.Add(ctx, t, "normal-date-queue")
		if err != nil {
			log.Errorf(ctx, "Failed to enqueue: %v", err)
		}
	}
}
