package scheduler

import (
	"context"
	"net/url"

	"github.com/SSHZ-ORG/tree-diagram/paths"
	"github.com/pkg/errors"
	"google.golang.org/appengine/v2/taskqueue"
)

const oneOffQueueName = "one-off-queue"

func ScheduleOneOff(ctx context.Context, params url.Values) error {
	task := taskqueue.NewPOSTTask(paths.CronOneOff, params)
	_, err := taskqueue.Add(ctx, task, oneOffQueueName)
	return errors.Wrap(err, "taskqueue.Add failed")
}
