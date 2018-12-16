package crawler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/civil"
	"github.com/PuerkitoBio/goquery"
	"github.com/SSHZ-ORG/tree-diagram/models"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

const (
	pageSize              = 30
	minNoteCountThreshold = 10

	eventListPageURLTemplate = "https://www.eventernote.com/events/search?year=%d&month=%d&day=%d&limit=%d"
)

func CrawlDateOnePage(ctx context.Context, date civil.Date) error {
	url := fmt.Sprintf(eventListPageURLTemplate, date.Year, date.Month, date.Day, pageSize)
	return crawlEventSearchPage(ctx, url)
}

func crawlEventSearchPage(ctx context.Context, url string) error {
	ts := time.Now()

	log.Infof(ctx, "Crawling event search page %v", url)

	client := urlfetch.Client(ctx)
	res, err := client.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errors.New("Status Error: " + res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	aMap := make(map[string]string)
	pMap := make(map[string]string)

	var es []*models.Event
	var eventAs [][]string
	var eventPs []string

	events := doc.Find(".gb_event_list").Children().Children()
	events.Each(func(i int, s *goquery.Selection) {
		e := &models.Event{}

		// Critical. If this is less than threshold, skip.
		e.LastNoteCount, err = strconv.Atoi(strings.TrimSpace(s.Find(".note_count").Text()))
		if err != nil {
			logError(ctx, url, i, "NoteCount", err)
			return
		}
		if e.LastNoteCount < minNoteCountThreshold {
			return
		}

		// Critical. If fails skip the event.
		e.ID, e.Name, err = parseLinkWithID(s.Find(".event h4 a"))
		if err != nil {
			logError(ctx, url, i, "ID / Name", err)
			return
		}

		// Critical. If fails skip the event.
		date, err := civil.ParseDate(strings.Split(strings.TrimSpace(s.Find(".date").Children().First().Text()), " ")[0])
		if err != nil {
			logError(ctx, url, i, "Date", err)
			return
		}
		e.Date = date.In(time.UTC)

		// Non-critical.
		pElement := s.Find(".place a")
		pID := ""
		if len(pElement.Nodes) > 0 {
			var pName string
			pID, pName, err = parseLinkWithID(pElement)
			if err != nil {
				logError(ctx, url, i, "Place", err)
				return
			}
			pMap[pID] = pName
		}

		// Non-critical.
		var aIDs []string
		s.Find(".actor a").Each(func(i int, as *goquery.Selection) {
			aID, aName, err := parseLinkWithID(as)
			if err != nil {
				logError(ctx, url, i, "Actors", err)
				return // Skips this actor only.
			}
			aMap[aID] = aName
			aIDs = append(aIDs, aID)
		})

		// Non-critical.
		if detailedTime := strings.TrimSpace(s.Find(".place .s").Text()); detailedTime != "" {
			tokens := strings.Split(detailedTime, " ")
			if len(tokens) == 6 && tokens[0] == "開場" && tokens[2] == "開演" && tokens[4] == "終演" {
				e.OpenTime, err = parseDetailedTime(tokens[1], date)
				if err != nil {
					logError(ctx, url, i, "OpenTime", err)
				}
				e.StartTime, err = parseDetailedTime(tokens[3], date)
				if err != nil {
					logError(ctx, url, i, "StartTime", err)
				}
				e.EndTime, err = parseDetailedTime(tokens[5], date)
				if err != nil {
					logError(ctx, url, i, "EndTime", err)
				}
			} else {
				log.Errorf(ctx, "URL %s item %d: Unknown detailed time: %s.", url, i, detailedTime)
			}
		}

		// Won't fail.
		e.LastUpdateTime = ts

		eventPs = append(eventPs, pID)
		eventAs = append(eventAs, aIDs)
		es = append(es, e)
	})

	aKeys, err := models.EnsureActors(ctx, aMap)
	if err != nil {
		log.Errorf(ctx, "EnsureActors: %v", err)
		return err
	}

	pKeys, err := models.EnsurePlaces(ctx, pMap)
	if err != nil {
		log.Errorf(ctx, "EnsurePlaces: %v", err)
		return err
	}

	for i, e := range es {
		if eventPs[i] != "" {
			e.Place = pKeys[eventPs[i]]
		}

		for _, aID := range eventAs[i] {
			e.Actors = append(e.Actors, aKeys[aID])
		}
	}

	if err := models.InsertOrUpdateEvents(ctx, es); err != nil {
		return err
	}

	log.Infof(ctx, "Updated %d events.", len(es))
	return nil
}

func logError(ctx context.Context, url string, i int, field string, err error) {
	log.Errorf(ctx, "Error on URL %s item %d: Cannot parse %s: %v.", url, i, field, err)
}

func parseLinkWithID(s *goquery.Selection) (string, string, error) {
	if href, ok := s.Attr("href"); ok {
		tmp := strings.Split(href, "/")
		return tmp[len(tmp)-1], strings.TrimSpace(s.Text()), nil
	}

	return "", "", errors.New(fmt.Sprintf("Error trying to parse link with ID: %v", s))
}

func parseDetailedTime(s string, date civil.Date) (time.Time, error) {
	if s == "-" {
		// If unknown, EventerNote returns this.
		return time.Time{}, nil
	}

	t, err := time.Parse("15:04", s)
	if err != nil {
		return time.Time{}, err
	}

	civilTime := civil.TimeOf(t)

	l, _ := time.LoadLocation("Asia/Tokyo")
	return civil.DateTime{Date: date, Time: civilTime}.In(l), nil
}
