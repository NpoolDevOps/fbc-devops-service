module github.com/NpoolDevOps/fbc-devops-service

go 1.16

require (
	github.com/EntropyPool/entropy-logger v0.0.0-20210320022718-3091537e035f
	github.com/NpoolDevOps/fbc-auth-service v0.0.0-20210411082937-64d2220ba926
	github.com/NpoolDevOps/fbc-license-service v0.0.0-20210726032917-16eea84e1d8c
	github.com/NpoolRD/http-daemon v0.0.0-20210324100344-82fee56de8ac
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/google/uuid v1.2.0
	github.com/jinzhu/gorm v1.9.16
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
