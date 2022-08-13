package utils

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/SSHZ-ORG/tree-diagram/pb"
	"google.golang.org/protobuf/proto"
)

var JST = func() *time.Location {
	l, _ := time.LoadLocation("Asia/Tokyo")
	return l
}()

func JSTToday() civil.Date {
	return civil.DateOf(time.Now().In(JST))
}

func ToProtoDate(d civil.Date) *pb.Date {
	return &pb.Date{
		Year:  proto.Int32(int32(d.Year)),
		Month: proto.Int32(int32(d.Month)),
		Day:   proto.Int32(int32(d.Day)),
	}
}
