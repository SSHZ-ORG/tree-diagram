package models

import (
	"fmt"
	"testing"
	"time"

	"github.com/qedus/nds"
	"golang.org/x/net/context"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
)

const testEventID = "0"

func newContext(t *testing.T) (context.Context, func()) {
	c, closeFunc, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	return c, closeFunc
}

func createSnapshot(timestamp int64, noteCount int, actors ...*datastore.Key) *EventSnapshot {
	return &EventSnapshot{
		EventID:   "0",
		Timestamp: time.Unix(timestamp, 0).UTC(),
		NoteCount: noteCount,
		Actors:    actors,
	}
}

func putSnapshots(ctx context.Context, s []*EventSnapshot) error {
	for _, snapshot := range s {
		if _, err := nds.Put(ctx, datastore.NewIncompleteKey(ctx, eventSnapshotKind, getEventKey(ctx, testEventID)), snapshot); err != nil {
			return err
		}
	}
	return nil
}

func fakeActorID(ctx context.Context, id string) *datastore.Key {
	return datastore.NewKey(ctx, "Unused", id, 0, nil)
}

func TestCompressSnapshots_NoPreviousCES(t *testing.T) {
	ctx, cancel := newContext(t)
	defer cancel()

	actor1 := fakeActorID(ctx, "1")
	actor2 := fakeActorID(ctx, "2")

	expected := []*EventSnapshot{
		createSnapshot(10, 1, actor1),
		createSnapshot(20, 1),
		createSnapshot(30, 2),
		createSnapshot(40, 2, actor1, actor2),
		createSnapshot(50, 2),
		createSnapshot(60, 1, actor2),
		createSnapshot(70, 1),
	}

	if err := putSnapshots(ctx, expected); err != nil {
		t.Fatal(err)
	}

	if err := CompressSnapshots(ctx, testEventID); err != nil {
		t.Fatal(err)
	}

	snapshots, err := getSnapshotsForEvent(ctx, getEventKey(ctx, testEventID))
	if err != nil {
		t.Fatal(err)
	}

	checkEquals(t, expected, snapshots)
}

func TestCompressSnapshots_HasPreviousCES_FirstCanCompress(t *testing.T) {
	ctx, cancel := newContext(t)
	defer cancel()

	actor1 := fakeActorID(ctx, "1")
	actor2 := fakeActorID(ctx, "2")

	expected := []*EventSnapshot{
		createSnapshot(10, 1, actor1),
	}

	if err := putSnapshots(ctx, expected); err != nil {
		t.Fatal(err)
	}

	if err := CompressSnapshots(ctx, testEventID); err != nil {
		t.Fatal(err)
	}

	moreSnapshots := []*EventSnapshot{
		createSnapshot(20, 1),
		createSnapshot(30, 2),
		createSnapshot(40, 2, actor1, actor2),
		createSnapshot(50, 2),
		createSnapshot(60, 1, actor2),
		createSnapshot(70, 1),
	}

	if err := putSnapshots(ctx, moreSnapshots); err != nil {
		t.Fatal(err)
	}

	expected = append(expected, moreSnapshots...)

	if err := CompressSnapshots(ctx, testEventID); err != nil {
		t.Fatal(err)
	}

	snapshots, err := getSnapshotsForEvent(ctx, getEventKey(ctx, testEventID))
	if err != nil {
		t.Fatal(err)
	}

	checkEquals(t, expected, snapshots)
}

func TestCompressSnapshots_HasPreviousCES_FirstCannotCompress(t *testing.T) {
	ctx, cancel := newContext(t)
	defer cancel()

	actor1 := fakeActorID(ctx, "1")
	actor2 := fakeActorID(ctx, "2")

	expected := []*EventSnapshot{
		createSnapshot(10, 1, actor1),
		createSnapshot(20, 1),
	}

	if err := putSnapshots(ctx, expected); err != nil {
		t.Fatal(err)
	}

	if err := CompressSnapshots(ctx, testEventID); err != nil {
		t.Fatal(err)
	}

	moreSnapshots := []*EventSnapshot{
		createSnapshot(30, 2),
		createSnapshot(40, 2, actor1, actor2),
		createSnapshot(50, 2),
		createSnapshot(60, 1, actor2),
		createSnapshot(70, 1),
	}

	if err := putSnapshots(ctx, moreSnapshots); err != nil {
		t.Fatal(err)
	}

	expected = append(expected, moreSnapshots...)

	if err := CompressSnapshots(ctx, testEventID); err != nil {
		t.Fatal(err)
	}

	snapshots, err := getSnapshotsForEvent(ctx, getEventKey(ctx, testEventID))
	if err != nil {
		t.Fatal(err)
	}

	checkEquals(t, expected, snapshots)
}

func checkEquals(t *testing.T, expected, actual []*EventSnapshot) {
	if len(expected) != len(actual) {
		t.Fatalf("len(expected) != len(actual): %d -> %d", len(expected), len(actual))
	}

	for i, s := range expected {
		es := fmt.Sprintf("%v+", s)
		as := fmt.Sprintf("%v+", actual[i])
		if es != as {
			t.Errorf("Found diff in [%d] expected -> actual: %s -> %s", i, es, as)
		}
	}
}
