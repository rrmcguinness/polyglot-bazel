module example

go 1.20

replace (
	examples/go/pkg/client => ./pkg/client
	examples/go/pkg/model => ./pkg/model
	examples/go/pkg/server => ./pkg/server
	examples/go/pkg/service => ./pkg/service
)

require (
	github.com/google/uuid v1.1.2
	github.com/stretchr/testify v1.8.2
	google.golang.org/grpc v1.49.0
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
