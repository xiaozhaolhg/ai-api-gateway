module github.com/ai-api-gateway/monitor-service

go 1.21

require (
	github.com/ai-api-gateway/api v0.0.0
	google.golang.org/grpc v1.60.1
	google.golang.org/protobuf v1.32.0
	gorm.io/driver/sqlite v1.5.4
	gorm.io/gorm v1.25.5
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/ai-api-gateway/api => ../api
