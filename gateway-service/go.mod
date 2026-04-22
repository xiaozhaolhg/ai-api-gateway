module github.com/ai-api-gateway/gateway-service

go 1.21

require (
	github.com/ai-api-gateway/api v0.0.0
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.33.0
	gopkg.in/yaml.v3 v3.0.1
	github.com/gin-gonic/gin v1.9.1
)

replace github.com/ai-api-gateway/api => ../api
