package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/SSHZ-ORG/tree-diagram/apicache"
	"github.com/SSHZ-ORG/tree-diagram/models"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

const (
	actorPageSize = 500

	actorSearchAPITemplate = "https://www.eventernote.com/api/actors/search?limit=%d&offset=%d&simple=1"
)

// Crawls one page at the given offset. offset is 1-index.
// Returns the next offset to use, or 0 if it's the last page.
// Errors wrapped.
func CrawlActorOnePage(ctx context.Context, offset int) (int, error) {
	url := fmt.Sprintf(actorSearchAPITemplate, actorPageSize, offset)

	ts := time.Now()

	log.Infof(ctx, "Crawling actor API page %v", url)

	ctxWithTimeout, _ := context.WithTimeout(ctx, crawlHTTPTimeout)
	client := urlfetch.Client(ctxWithTimeout)
	res, err := client.Get(url)
	if err != nil {
		return 0, errors.Wrap(err, "URL fetch failed")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return 0, errors.New("URL Fetch returned unexpected status: " + res.Status)
	}

	jsonBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, errors.Wrap(err, "ioutil.ReadAll failed to read response")
	}

	json := gjson.ParseBytes(jsonBytes)

	var actorIDs []string
	json.Get("results.#.id").ForEach(func(key, value gjson.Result) bool {
		actorIDs = append(actorIDs, strconv.FormatInt(value.Int(), 10))
		return true
	})

	actorMap, err := models.GetActorMap(ctx, actorIDs)
	if err != nil {
		return 0, err
	}

	var updatedActorIDs []string
	var oas, nas []*models.Actor
	json.Get("results").ForEach(func(key, value gjson.Result) bool {
		id := strconv.FormatInt(value.Get("id").Int(), 10)
		if oa, ok := actorMap[id]; ok {
			na := &models.Actor{
				ID:                id,
				Name:              value.Get("name").String(),
				LastFavoriteCount: int(value.Get("favorite_count").Int()),
				LastUpdateTime:    ts,
			}
			if oa.LastUpdateTime.IsZero() || !oa.Equal(na) {
				updatedActorIDs = append(updatedActorIDs, id)
				oas = append(oas, oa)
				nas = append(nas, na)
			}
		}
		return true
	})

	err = models.UpdateActors(ctx, oas, nas)
	if err != nil {
		return 0, err
	}

	if err := apicache.ClearRenderActors(ctx, updatedActorIDs); err != nil {
		return 0, err
	}

	log.Infof(ctx, "Updated %d actors.", len(nas))

	total := int(json.Get("info.total").Int())
	current := offset + int(json.Get("info.return_count").Int()) - 1 // Because offset is 1-based.
	if current >= total {
		return 0, nil
	}
	return current + 1, nil
}
