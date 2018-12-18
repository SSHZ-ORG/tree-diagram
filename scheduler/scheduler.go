package scheduler

import (
	"net/url"
	"strconv"

	"cloud.google.com/go/civil"
	"github.com/SSHZ-ORG/tree-diagram/paths"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
)

const (
	normalQueueName = "normal-date-queue"
)

func ScheduleNormalQueue(ctx context.Context, date civil.Date, page int) error {
	t := taskqueue.NewPOSTTask(paths.CrawlDatePath, url.Values{
		"date": []string{date.String()},
		"page": []string{strconv.Itoa(page)},
	})

	_, err := taskqueue.Add(ctx, t, normalQueueName)
	if err != nil {
		log.Errorf(ctx, "Failed to enqueue: %v", err)
	}
	return err
}

func EnqueueCrawlDateRange(ctx context.Context, begin, end civil.Date) error {
	for cur := begin; cur.Before(end); cur = cur.AddDays(1) {
		if err := ScheduleNormalQueue(ctx, cur, 1); err != nil {
			return err
		}
	}
	return nil
}
