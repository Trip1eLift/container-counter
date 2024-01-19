package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Trip1eLift/container-counter/cluster/container_counter_system"
)

const ip = "0.0.0.0"
const port = "8000"

func main() {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			count := container_counter_system.GetCount()
			fmt.Printf("container %s with count: %d\n", os.Getenv("CONTAINER_NUM"), count)
		}
	}()

	http.HandleFunc("/health", func(write http.ResponseWriter, request *http.Request) {
		log.Println("Golang healthcheck.")
		fmt.Fprintf(write, "Healthy golang server.\n")
	})
	http.HandleFunc("/traffic", func(write http.ResponseWriter, request *http.Request) {
		container_counter_system.OnTraffic()
	})
	log.Println(fmt.Sprintf("Listening on %s:%s", ip, port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), nil))
}
