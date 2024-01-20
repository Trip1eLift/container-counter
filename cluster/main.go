package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Trip1eLift/container-counter/cluster/container_counter_system"
)

const ip = "0.0.0.0"
const port = "8000"

func main() {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			container_counter_system.GetCount()
			//fmt.Printf("count: %d\n", count)
		}
	}()

	http.HandleFunc("/health", func(write http.ResponseWriter, request *http.Request) {
		// log.Println("Golang healthcheck.")
		fmt.Fprintf(write, "Healthy golang server.\n")
	})
	http.HandleFunc("/traffic", func(write http.ResponseWriter, request *http.Request) {
		container_counter_system.OnTraffic()
		fmt.Fprintf(write, "Traffic recieved.\n")
	})
	log.Println(fmt.Sprintf("Listening on %s:%s", ip, port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", ip, port), nil))
}
