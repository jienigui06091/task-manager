package main

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"task-manager/internal/cache"
	"task-manager/internal/middleware"
	"task-manager/internal/model"
	"task-manager/internal/router"
)

func main() {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("load .env: %v", err)
	}

	if err := model.InitDB(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := model.CloseDB(); err != nil {
			log.Printf("close database: %v", err)
		}
	}()

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("REDIS_DB must be an integer: %v", err)
	}

	cache.InitRedis(
		os.Getenv("REDIS_ADDR"),
		os.Getenv("REDIS_USERNAME"),
		os.Getenv("REDIS_PASSWORD"),
		redisDB,
	)
	defer func() {
		if err := cache.RedisClient.Close(); err != nil {
			log.Printf("close redis: %v", err)
		}
	}()

	if err := cache.RedisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}

	r := gin.New()
	r.Use(middleware.ErrorHandler())
	router.SetRouter(r)

	log.Println("server running on :8089")
	if err := r.Run(":8089"); err != nil {
		log.Fatal(err)
	}
}
