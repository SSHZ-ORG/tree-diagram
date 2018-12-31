package reporter

import (
	"bytes"
	"encoding/gob"
	"sort"
	"time"

	"cloud.google.com/go/civil"
	"github.com/SSHZ-ORG/tree-diagram/utils"
	"golang.org/x/net/context"
)

const (
	increaseMaxItems      = 30
	increaseDiffThreshold = 30

	decreaseMaxItems      = 10
	decreaseDiffThreshold = -1
)

type ReportItem struct {
	EventID string
	Name    string

	Date string

	OldNoteCount  int
	NewNoteCount  int
	DiffNoteCount int
}

type reportItemList []ReportItem

func (v reportItemList) Len() int {
	return len(v)
}
func (v reportItemList) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}
func (v reportItemList) Less(i, j int) bool {
	return v[i].DiffNoteCount < v[j].DiffNoteCount
}

type report struct {
	Increase reportItemList
	Decrease reportItemList
}

func AddToTodayReport(ctx context.Context, ris []ReportItem) error {
	return addToReport(ctx, ris, civil.DateOf(time.Now().In(time.UTC)))
}

func addToReport(ctx context.Context, ris []ReportItem, d civil.Date) error {
	r, err := getReportFromMemory(ctx, d)
	if err != nil {
		return err
	}

	var knownIDs []string
	for _, i := range r.Increase {
		knownIDs = append(knownIDs, i.EventID)
	}
	for _, i := range r.Decrease {
		knownIDs = append(knownIDs, i.EventID)
	}
	knownIDsSet := utils.NewStringSetFromSlice(knownIDs)

	for _, item := range ris {
		if knownIDsSet.Contains(item.EventID) {
			continue
		}

		if item.DiffNoteCount >= increaseDiffThreshold {
			r.Increase = append(r.Increase, item)
			knownIDsSet.Add(item.EventID)
		}
		if item.DiffNoteCount <= decreaseDiffThreshold {
			r.Decrease = append(r.Decrease, item)
			knownIDsSet.Add(item.EventID)
		}
	}

	sort.Sort(sort.Reverse(r.Increase))
	if len(r.Increase) > increaseMaxItems {
		r.Increase = r.Increase[:increaseMaxItems]
	}

	sort.Sort(r.Decrease)
	if len(r.Decrease) > decreaseMaxItems {
		r.Decrease = r.Decrease[:decreaseMaxItems]
	}

	return putReportIntoMemory(ctx, d, r)
}

func RenderYesterdayReport(ctx context.Context) (string, error) {
	return RenderReport(ctx, civil.DateOf(time.Now().In(time.UTC)).AddDays(-1))
}

func RenderReport(ctx context.Context, d civil.Date) (string, error) {
	report, err := getReportFromMemory(ctx, d)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	err = reportTemplate.Execute(&b, report)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func init() {
	gob.Register(report{})
}
