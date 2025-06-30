package main

import (
	"go-auth/cmd/router"
	"go-auth/internal/metadata/account"
	"go-auth/internal/metadata/account_platform"
	"go-auth/internal/metadata/device"
	"go-auth/internal/metadata/verification_code"
	"go-auth/internal/token"

	"github.com/shanelex111/go-common/pkg/cache/redis"
	"github.com/shanelex111/go-common/pkg/config"
	"github.com/shanelex111/go-common/pkg/db/mysql"
	"github.com/shanelex111/go-common/pkg/engine"
	"github.com/shanelex111/go-common/pkg/log"
	"github.com/shanelex111/go-common/pkg/server"
	"github.com/shanelex111/go-common/third_party/geo"
)

func main() {
	// 1. load config - 加载配置文件
	v, err := config.Load(".", "config")
	if err != nil {
		panic(err)
	}

	// 2. init server components - 初始化组件
	server.Init(v,
		log.Init,
		engine.Init,

		mysql.Init,
		redis.Init,

		geo.Init,
	)

	// 3. load service configs - 加载业务配置
	server.Load(v,
		account.Load,
		account_platform.Load,
		device.Load,
		verification_code.Load,

		token.Load,
	)

	// 4. services run - 业务启动
	server.Run(router.Run)
}
