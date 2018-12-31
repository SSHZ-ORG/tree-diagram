package reporter

import (
	"bytes"
	"encoding/gob"

	"cloud.google.com/go/civil"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
)

const (
	memcacheKeyPrefex = "TDREP:D1:"
	datastoreKind     = "Report"
)

type reportDatastoreWrapper struct {
	Data []byte
}

func getMemcacheKey(date civil.Date) string {
	return memcacheKeyPrefex + date.String()
}

func getDatastoreKey(ctx context.Context, date civil.Date) *datastore.Key {
	return datastore.NewKey(ctx, datastoreKind, date.String(), 0, nil)
}

func putReportIntoMemory(ctx context.Context, date civil.Date, r report) error {
	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(r); err != nil {
		return err
	}

	data := b.Bytes()

	err := memcache.Set(ctx, &memcache.Item{
		Key:   getMemcacheKey(date),
		Value: data,
	})
	if err != nil {
		return err
	}

	_, err = datastore.Put(ctx, getDatastoreKey(ctx, date), &reportDatastoreWrapper{Data: data})
	return err
}

func getReportFromMemory(ctx context.Context, date civil.Date) (report, error) {
	var data []byte

	if item, err := memcache.Get(ctx, getMemcacheKey(date)); err == nil {
		data = item.Value
	} else {
		w := &reportDatastoreWrapper{}
		err := datastore.Get(ctx, getDatastoreKey(ctx, date), w)
		if err != nil && err != datastore.ErrNoSuchEntity {
			return report{}, err
		}
		data = w.Data
	}

	if len(data) == 0 {
		return report{}, nil
	}

	r := report{}
	err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(&r)
	return r, err
}
