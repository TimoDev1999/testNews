package db

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-pg/pg/v10"
	"github.com/spf13/viper"
)

var DB *pg.DB

func Init() {
	viper.AutomaticEnv()

	poolSizeStr := os.Getenv("DATABASE_POOLSIZE")
	if poolSizeStr == "" {
		poolSizeStr = "10"
	}

	poolSize, err := strconv.Atoi(poolSizeStr)
	if err != nil {
		log.Fatalf("Invalid pool size: %s", err)
	}

	DB = pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_PORT")),
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Database: os.Getenv("DATABASE_DBNAME"),
		PoolSize: poolSize,
	})
}
