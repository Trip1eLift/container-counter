package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const ip = "0.0.0.0"
const port = "8000"

func main() {
	redisOps(context.Background())

	http.HandleFunc("/health", func(write http.ResponseWriter, request *http.Request) {
		log.Println("Golang healthcheck.")
		fmt.Fprintf(write, "Healthy golang server.\n")
	})
	log.Println(fmt.Sprintf("Listening on %s:%s", ip, port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), nil))
}

func redisOps(ctx context.Context) {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       Atoi(os.Getenv("REDIS_DB")),
	})

	err := client.Ping(ctx).Err()
	if err != nil {
		log.Fatal(fmt.Sprintf("Redis connection err: %v", err))
	}
}

func Atoi(s string) int {
	num, _ := strconv.Atoi(s)
	return num
}
