package scheduler

import (
	"context"
	"net/url"
	"strconv"

	"github.com/SSHZ-ORG/tree-diagram/paths"
	"github.com/pkg/errors"
	"google.golang.org/appengine/v2/taskqueue"
)

const actorQueueName = "actor-queue"

// Remember that offset is 1-based.
// Errors wrapped.
func ScheduleCrawlActorPage(ctx context.Context, offset int) error {
	task := taskqueue.NewPOSTTask(paths.CrawlActorPath, url.Values{
		"offset": []string{strconv.Itoa(offset)},
	})
	_, err := taskqueue.Add(ctx, task, actorQueueName)
	return errors.Wrap(err, "taskqueue.Add failed")
}
