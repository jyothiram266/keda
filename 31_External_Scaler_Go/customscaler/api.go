package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func setValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	number := vars["number"]
	CustomQueueLength, _ = strconv.ParseInt(number, 10, 64)
	log.Printf("new value: %d\n", CustomQueueLength)
}

func RunManagementApi() {
	r := mux.NewRouter()
	r.HandleFunc("/api/queue/{number:[0-9]+}", setValue).Methods("POST")
	http.Handle("/", r)
	fmt.Printf("Running http management server on port: %d\n", 9090)
	http.ListenAndServe(":9090", nil)
}

var CustomQueueLength int64 = 0

func reduceCustomQueueLength() {
	for {
		if CustomQueueLength > 0 {
			CustomQueueLength--
			log.Printf("Reduced queue length value: %d\n", CustomQueueLength)
			time.Sleep(1 * time.Minute)
		}
	}
}
