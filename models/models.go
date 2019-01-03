package models

import (
	"github.com/scylladb/go-set/strset"
	"google.golang.org/appengine/datastore"
)

var AllKinds = []string{actorKind, placeKind, eventKind, eventSnapshotKind}

func areKeysSetsEqual(a, b []*datastore.Key) bool {
	as := strset.New()
	bs := strset.New()

	for _, k := range a {
		as.Add(k.String())
	}
	for _, k := range b {
		bs.Add(k.String())
	}

	return as.IsEqual(bs)
}
