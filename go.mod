module github.com/SSHZ-ORG/tree-diagram

go 1.15

require (
	cloud.google.com/go v0.101.0
	github.com/PuerkitoBio/goquery v1.8.0
	github.com/improbable-eng/grpc-web v0.15.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/pkg/errors v0.9.1
	github.com/qedus/nds v1.0.0
	github.com/scylladb/go-set v1.0.2
	github.com/tidwall/gjson v1.14.1
	golang.org/x/oauth2 v0.0.0-20220411215720-9780585627b5
	google.golang.org/api v0.77.0
	google.golang.org/appengine/v2 v2.0.1
	google.golang.org/grpc v1.46.0
	google.golang.org/protobuf v1.28.0
)

replace github.com/qedus/nds => github.com/SSHZ-ORG/nds v1.0.1-0.20220220041449-5427bae4887c
