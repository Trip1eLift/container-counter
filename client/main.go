package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const ip = "0.0.0.0"
const port = "7000"

var mode int = 0

type toggle_mode struct {
	mode int
}

func main() {
	http.HandleFunc("/health", func(write http.ResponseWriter, request *http.Request) {
		log.Println("Golang healthcheck.")
		fmt.Fprintf(write, "Healthy golang server.\n")
	})
	http.HandleFunc("/toggle", func(write http.ResponseWriter, request *http.Request) {
		fmt.Printf("req body: %v\n", request.Body)

		var tm toggle_mode
		body, _ := ioutil.ReadAll(request.Body)
		err := json.Unmarshal(body, &tm)
		// decoder := json.NewDecoder(request.Body)
		// var tm toggle_mode
		// err := decoder.Decode(&tm)
		if err != nil {
			fmt.Fprintf(write, "Failed to toggle due to decoding error: %v\n", err)
			return
		} else {
			mode = tm.mode
			fmt.Printf("Toggled to mode: %d\n", tm.mode)
			fmt.Fprintf(write, "Toggled to %d.\n", tm.mode)
			return
		}
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
