package scheduler

import (
	"net/url"

	"github.com/SSHZ-ORG/tree-diagram/paths"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/taskqueue"
)

const oneOffQueueName = "one-off-queue"

func ScheduleOneOff(ctx context.Context, params url.Values) error {
	task := taskqueue.NewPOSTTask(paths.CronOneOff, params)
	_, err := taskqueue.Add(ctx, task, oneOffQueueName)
	return errors.Wrap(err, "taskqueue.Add failed")
}
