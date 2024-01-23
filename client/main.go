package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const ip = "0.0.0.0"
const port = "7001"

var mode int = 0
var interval int = 3

func init() {
	addresses := []string{
		os.Getenv("CONTAINER_1"),
		os.Getenv("CONTAINER_2"),
		os.Getenv("CONTAINER_3"),
		os.Getenv("CONTAINER_4"),
		os.Getenv("CONTAINER_5"),
		os.Getenv("CONTAINER_6"),
		os.Getenv("CONTAINER_7"),
		os.Getenv("CONTAINER_8"),
		os.Getenv("CONTAINER_9"),
		os.Getenv("CONTAINER_10"),
	}

	go func() {
		for {
			time.Sleep(time.Duration(interval) * time.Second)
			for i := 0; i < mode; i++ {
				http.Get(fmt.Sprintf("http://%s/traffic", addresses[i]))
				fmt.Printf("Sending request to container %d: %s\n", i+1, addresses[i])
				// resp, err := http.Get(fmt.Sprintf("http://%s/traffic", addresses[i]))
				// fmt.Printf("resp: %v err: %v\n", resp, err)
			}
		}
	}()
}

func main() {
	http.HandleFunc("/health", func(write http.ResponseWriter, request *http.Request) {
		// log.Println("Golang healthcheck.")
		fmt.Fprintf(write, "Healthy golang server.\n")
	})
	http.HandleFunc("/toggle", func(write http.ResponseWriter, request *http.Request) {
		data := make(map[string]string)
		err := json.NewDecoder(request.Body).Decode(&data)
		if err != nil {
			fmt.Fprintf(write, "Failed to toggle due to decoding error: %v\n", err)
			return
		}
		//fmt.Printf("req body: %v\n", data)

		mode = atoi(data["mode"])
		interval = atoi(data["interval"])
		default_case()

		fmt.Printf("Toggled to mode: %d with interval: %d\n", mode, interval)
		fmt.Fprintf(write, "Toggled to mode: %d with interval: %d\n", mode, interval)
		return
	})
	log.Println(fmt.Sprintf("Listening on %s:%s", ip, port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), nil))
}

func atoi(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return num
}

func default_case() {
	if mode < 0 {
		mode = 0
	}
	if mode > 10 {
		mode = 10
	}
	if interval <= 0 {
		interval = 3
	}
}
