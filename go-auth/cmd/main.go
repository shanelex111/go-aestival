package main

import (
	"github.com/shanelex111/go-common/pkg/cache/redis"
	"github.com/shanelex111/go-common/pkg/config"
	"github.com/shanelex111/go-common/pkg/db/mysql"
	"github.com/shanelex111/go-common/pkg/engine"
	"github.com/shanelex111/go-common/pkg/server"
)

func main() {
	v, err := config.Load(".", "config")
	if err != nil {
		panic(err)
	}

	// server init
	server.Init(v, engine.Init, mysql.Init, redis.Init)
}
