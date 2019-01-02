package scheduler

import (
	"net/http"
	"net/url"
	"strconv"

	"cloud.google.com/go/civil"
	"github.com/SSHZ-ORG/tree-diagram/paths"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
)

type DateQueue string

const (
	NormalDateQueue    = DateQueue("normal-date-queue")
	ThrottledDateQueue = DateQueue("throttled-date-queue")
	OnDemandDateQueue  = DateQueue("ondemand-date-queue")
)

func newCrawlDateTask(date civil.Date, page int) *taskqueue.Task {
	return taskqueue.NewPOSTTask(paths.CrawlDatePath, url.Values{
		"date": []string{date.String()},
		"page": []string{strconv.Itoa(page)},
	})
}

func (q DateQueue) Schedule(ctx context.Context, date civil.Date, page int) error {
	_, err := taskqueue.Add(ctx, newCrawlDateTask(date, page), string(q))
	if err != nil {
		log.Errorf(ctx, "Failed to enqueue: %v", err)
	}
	return err
}

func (q DateQueue) EnqueueDateRange(ctx context.Context, start, end civil.Date) error {
	var ts []*taskqueue.Task
	for cur := start; cur.Before(end); cur = cur.AddDays(1) {
		ts = append(ts, newCrawlDateTask(cur, 1))
	}
	_, err := taskqueue.AddMulti(ctx, ts, string(q))
	if err != nil {
		log.Errorf(ctx, "Failed to enqueue: %v", err)
	}
	return err
}

func (q DateQueue) CurrentTaskCount(ctx context.Context) (int, error) {
	s, err := taskqueue.QueueStats(ctx, []string{string(q)})
	if err != nil {
		return 0, err
	}

	return s[0].Tasks, nil
}

func GetCurrentQueue(r *http.Request) DateQueue {
	qn := r.Header.Get("X-AppEngine-QueueName")
	if qn == "" {
		return OnDemandDateQueue
	}
	return DateQueue(qn)
}
