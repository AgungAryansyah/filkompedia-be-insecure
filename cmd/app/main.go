package main

import "github.com/AgungAryansyah/filkompedia-be-unsecure/pkg/config"

func main() {
	config.LoadEnv()

	app := config.StartFiber()
	db := config.StartSqlx()
	redis := config.StartRedis()

	config.StartUp(&config.Config{
		App:   app,
		DB:    db,
		Redis: redis,
	})
}
