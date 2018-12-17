package utils

import (
	"net/url"
	"strconv"

	"cloud.google.com/go/civil"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
)

const (
	normalQueueName = "normal-date-queue"
)

func ScheduleNormalQueue(ctx context.Context, date civil.Date, page int) error {
	t := taskqueue.NewPOSTTask("/admin/crawl/date", url.Values{
		"date": []string{date.String()},
		"page": []string{strconv.Itoa(page)},
	})

	_, err := taskqueue.Add(ctx, t, normalQueueName)
	if err != nil {
		log.Errorf(ctx, "Failed to enqueue: %v", err)
	}
	return err
}
