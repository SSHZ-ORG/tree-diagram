package scheduler

import (
	"math/rand"
	"net/http"
	"net/url"
	"strconv"

	"cloud.google.com/go/civil"
	"github.com/SSHZ-ORG/tree-diagram/paths"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/taskqueue"
)

type DateQueue string

const (
	NormalDateQueue    = DateQueue("normal-date-queue")
	ThrottledDateQueue = DateQueue("throttled-date-queue")
	DeadSlowDateQueue  = DateQueue("deadslow-date-queue")
	OnDemandDateQueue  = DateQueue("ondemand-date-queue")

	maxTasksPerAddMulti = 100
)

func newCrawlDateTask(date civil.Date, page int) *taskqueue.Task {
	return taskqueue.NewPOSTTask(paths.CrawlDatePath, url.Values{
		"date": []string{date.String()},
		"page": []string{strconv.Itoa(page)},
	})
}

// Errors wrapped.
func (q DateQueue) Schedule(ctx context.Context, date civil.Date, page int) error {
	_, err := taskqueue.Add(ctx, newCrawlDateTask(date, page), string(q))
	return errors.Wrap(err, "taskqueue.Add failed")
}

// Errors wrapped.
func (q DateQueue) EnqueueDateRange(ctx context.Context, start, end civil.Date, shuffle bool) error {
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
func (q DateQueue) CurrentTaskCount(ctx context.Context) (int, error) {
	s, err := taskqueue.QueueStats(ctx, []string{string(q)})
	if err != nil {
		return 0, errors.Wrap(err, "taskqueue.QueueStats failed")
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
