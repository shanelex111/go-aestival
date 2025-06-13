package main

import (
	"go-auth/cmd/router"
	"go-auth/internal/metadata/account"
	"go-auth/internal/metadata/account_platform"
	"go-auth/internal/metadata/aestival_member"
	"go-auth/internal/metadata/device"
	"go-auth/internal/metadata/verification_code"

	"github.com/shanelex111/go-common/pkg/config"
	"github.com/shanelex111/go-common/pkg/db/mysql"
	"github.com/shanelex111/go-common/pkg/engine"
	"github.com/shanelex111/go-common/pkg/server"
)

func main() {
	// 1. load config - 加载配置文件
	v, err := config.Load(".", "config")
	if err != nil {
		panic(err)
	}

	// 2. init server components - 初始化组件
	server.Init(v,
		engine.Init,
		mysql.Init,
		//redis.Init,
	)

	// 3. load service configs - 加载业务配置
	server.Load(v,
		account.Load,
		account_platform.Load,
		aestival_member.Load,
		device.Load,
		verification_code.Load,
	)

	// 4. services run - 业务启动
	server.Run(router.Run)
}
