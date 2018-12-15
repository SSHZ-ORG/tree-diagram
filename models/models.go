package models

import (
	"github.com/SSHZ-ORG/tree-diagram/utils"
	"google.golang.org/appengine/datastore"
)

func areKeysSetsEqual(a, b []*datastore.Key) bool {
	as := utils.NewStringSet()
	bs := utils.NewStringSet()

	for _, k := range a {
		as.Add(k.String())
	}
	for _, k := range b {
		bs.Add(k.String())
	}

	return as.Equals(bs)
}
