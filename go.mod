module github.com/NpoolDevOps/fbc-devops-service

go 1.15

require (
	github.com/EntropyPool/entropy-logger v0.0.0-20210210082337-af230fd03ce7
	github.com/NpoolDevOps/fbc-auth-service v0.0.0-20210323131841-28695bb4b6a8
	github.com/NpoolDevOps/fbc-license-service v0.0.0-20210328062839-d1527bc31f7e
	github.com/NpoolRD/http-daemon v0.0.0-20220506133728-7943c2cae9a7 // indirect
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/google/uuid v1.2.0
	github.com/jinzhu/gorm v1.9.16
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
