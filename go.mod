module github.com/SSHZ-ORG/tree-diagram

go 1.15

require (
	cloud.google.com/go v0.100.2
	github.com/PuerkitoBio/goquery v1.8.0
	github.com/improbable-eng/grpc-web v0.15.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/pkg/errors v0.9.1
	github.com/qedus/nds v1.0.0
	github.com/scylladb/go-set v1.0.2
	github.com/tidwall/gjson v1.14.0
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
	google.golang.org/api v0.69.0
	google.golang.org/appengine/v2 v2.0.1
	google.golang.org/grpc v1.44.0
	google.golang.org/protobuf v1.27.1
)

replace github.com/qedus/nds => github.com/SSHZ-ORG/nds v1.0.1-0.20220220041449-5427bae4887c
