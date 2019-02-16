package utils

import (
	"time"

	"cloud.google.com/go/civil"
)

func JST() *time.Location {
	l, _ := time.LoadLocation("Asia/Tokyo")
	return l
}

func JSTToday() civil.Date {
	return civil.DateOf(time.Now().In(JST()))
}
