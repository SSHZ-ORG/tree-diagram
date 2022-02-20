package scheduler

import (
	"context"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"

	"cloud.google.com/go/civil"
	"github.com/SSHZ-ORG/tree-diagram/paths"
	"github.com/pkg/errors"
	"google.golang.org/appengine/v2/taskqueue"
)

type EventDateQueue string

const (
	NormalDateQueue    = EventDateQueue("normal-date-queue")
	ThrottledDateQueue = EventDateQueue("throttled-date-queue")
	DeadSlowDateQueue  = EventDateQueue("deadslow-date-queue")
	OnDemandDateQueue  = EventDateQueue("ondemand-date-queue")
)

func newCrawlDateTask(date civil.Date, page int) *taskqueue.Task {
	return taskqueue.NewPOSTTask(paths.CrawlDatePath, url.Values{
		"date": []string{date.String()},
		"page": []string{strconv.Itoa(page)},
	})
}

// Errors wrapped.
func (q EventDateQueue) Schedule(ctx context.Context, date civil.Date, page int) error {
	_, err := taskqueue.Add(ctx, newCrawlDateTask(date, page), string(q))
	return errors.Wrap(err, "taskqueue.Add failed")
}

// Errors wrapped.
func (q EventDateQueue) EnqueueDateRange(ctx context.Context, start, end civil.Date, shuffle bool) error {
	var ts []*taskqueue.Task
	for cur := start; cur.Before(end); cur = cur.AddDays(1) {
		ts = append(ts, newCrawlDateTask(cur, 1))
	}

	var batches [][]*taskqueue.Task
	for maxTasksPerAddMulti < len(ts) {
		ts, batches = ts[maxTasksPerAddMulti:], append(batches, ts[0:maxTasksPerAddMulti:maxTasksPerAddMulti])
	}
	batches = append(batches, ts)

	if shuffle {
		rand.Shuffle(len(batches), func(i, j int) {
			batches[i], batches[j] = batches[j], batches[i]
		})
	}

	for _, batch := range batches {
		_, err := taskqueue.AddMulti(ctx, batch, string(q))
		if err != nil {
			return errors.Wrap(err, "taskqueue.AddMulti failed")
		}
	}
	return nil
}

// Errors wrapped.
func (q EventDateQueue) CurrentTaskCount(ctx context.Context) (int, error) {
	s, err := taskqueue.QueueStats(ctx, []string{string(q)})
	if err != nil {
		return 0, errors.Wrap(err, "taskqueue.QueueStats failed")
	}

	return s[0].Tasks, nil
}

func GetCurrentEventDateQueue(r *http.Request) EventDateQueue {
	qn := r.Header.Get("X-AppEngine-QueueName")
	if qn == "" {
		return OnDemandDateQueue
	}
	return EventDateQueue(qn)
}
