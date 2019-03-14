package scheduler

import (
	"net/url"

	"github.com/SSHZ-ORG/tree-diagram/paths"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
)

const compressEventSnapshotQueueName = "compress-event-queue"

// This just attempts to schedule it. Don't care if it fails.
func ScheduleCompressEventSnapshots(ctx context.Context, eventID string) {
	task := taskqueue.NewPOSTTask(paths.CommandCompressEventSnapshots, url.Values{
		"id": []string{eventID},
	})
	_, err := taskqueue.Add(ctx, task, compressEventSnapshotQueueName)
	if err != nil {
		log.Errorf(ctx, "taskqueue.Add: %+v", err)
	}
	log.Infof(ctx, "Scheduled compression for event %s", eventID)
}
