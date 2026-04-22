module github.com/ai-api-gateway/billing-service

go 1.21

require (
	github.com/ai-api-gateway/api v0.0.0
	github.com/google/uuid v1.3.1
	google.golang.org/grpc v1.64.0
	gopkg.in/yaml.v3 v3.0.1
)

require google.golang.org/protobuf v1.33.0 // indirect

replace github.com/ai-api-gateway/api => ../api
