module github.com/ai-api-gateway/auth-service

go 1.21

require (
	github.com/ai-api-gateway/api v0.0.0
	google.golang.org/grpc v1.60.1
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/driver/sqlite v1.5.4
	gorm.io/gorm v1.25.5
)

require google.golang.org/protobuf v1.32.0 // indirect

replace github.com/ai-api-gateway/api => ./api
