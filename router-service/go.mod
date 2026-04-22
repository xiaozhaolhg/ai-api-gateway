module github.com/ai-api-gateway/router-service

go 1.21

require (
	google.golang.org/grpc v1.60.1
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/driver/sqlite v1.5.4
	gorm.io/gorm v1.25.5
)

replace github.com/ai-api-gateway/api => ./api
