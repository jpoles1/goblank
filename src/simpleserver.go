package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

//const servport = "51339"
const servport string = "8000"
const time_spacer int = 4 //in seconds

var last_visit time.Time = time.Now().Add(-time.Minute)

var time_mutex *sync.Mutex = &sync.Mutex{}

func hello(w http.ResponseWriter, r *http.Request) {
	time_mutex.Lock()
	fmt.Println("Visitor here!")
	io.WriteString(w, "Hello world!")
	last_visit = time.Now()
	time_mutex.Unlock()
}

func main() {
	http.HandleFunc("/", hello)
	fmt.Println("Starting server on port: " + servport)
	log.Fatal(http.ListenAndServe(":"+servport, nil))
}
