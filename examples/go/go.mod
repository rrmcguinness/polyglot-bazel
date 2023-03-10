module example

go 1.20

replace (
	examples/go/pkg/client => ./pkg/client
	examples/go/pkg/model => ./pkg/model
	examples/go/pkg/server => ./pkg/server
	examples/go/pkg/service => ./pkg/service
)

require (
	cloud.google.com/go/bigquery v1.43.0
	cloud.google.com/go/logging v1.5.0
	github.com/stretchr/testify v1.8.1
	go.einride.tech/protobuf-bigquery v0.23.0
	google.golang.org/api v0.103.0
	google.golang.org/grpc v1.50.1
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/cenkalti/backoff/v4 v4.2.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	cloud.google.com/go v0.105.0 // indirect
	cloud.google.com/go/compute v1.12.1 // indirect
	cloud.google.com/go/compute/metadata v0.2.1 // indirect
	cloud.google.com/go/iam v0.7.0 // indirect
	cloud.google.com/go/longrunning v0.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.0 // indirect
	github.com/googleapis/gax-go/v2 v2.7.0 // indirect
	github.com/madflojo/tasks v1.0.4 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rs/xid v1.3.0 // indirect
	go.einride.tech/aip v0.58.1 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/net v0.0.0-20221014081412-f15817d10f9b // indirect
	golang.org/x/oauth2 v0.0.0-20221014153046-6fdb5e3db783 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.0.0-20220915200043-7b5979e65e41 // indirect
	golang.org/x/text v0.4.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20221027153422-115e99e71e1c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)