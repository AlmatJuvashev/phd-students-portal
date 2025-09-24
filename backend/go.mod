module phd-portal/backend

go 1.22.5

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/jmoiron/sqlx v1.4.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
)

require github.com/aws/aws-sdk-go-v2 v1.30.0

require (
	github.com/aws/aws-sdk-go-v2/config v1.27.33
	github.com/aws/aws-sdk-go-v2/credentials v1.17.43
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.16.20
	github.com/aws/aws-sdk-go-v2/service/s3 v1.63.0
)

require github.com/redis/go-redis/v9 v9.5.1
