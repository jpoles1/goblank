package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"
)

//const servport = "51339"
const servport string = "51339"
const rate_unit = time.Second
const rate_limit int = 4 //in the rate_unit

var last_visit time.Time = time.Now().Add(rate_unit * time.Duration(-rate_limit-1))

var time_mutex *sync.Mutex = &sync.Mutex{}

func blank_gen(f func()) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		time_mutex.Lock()
		fmt.Printf("\nVisitor here! Last visit was at: %s\n", last_visit)
		fmt.Printf("Last visit was %s ago\n", time.Since(last_visit))
		if time.Since(last_visit) > rate_unit*time.Duration(rate_limit) {
			io.WriteString(w, "Time Limit Check Passed! Proceeding...")
			f()
		} else {
			io.WriteString(w, "Time Limit Check Failed!")
		}
		last_visit = time.Now()
		time_mutex.Unlock()
	}
}
func blank() {
	_, err := exec.Command("cmd", "/C", "C:/Windows/System32/nircmd.exe", "monitor", "off").CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
}
func mute() {
	_, err := exec.Command("cmd", "/C", "C:/Windows/System32/nircmd.exe", "mutesysvolume", "1").CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
}
func main() {
	http.HandleFunc("/blank", blank_gen(func() {
		fmt.Println("Blanking...")
		blank()
	}))
	http.HandleFunc("/mute", blank_gen(func() {
		fmt.Println("Muting...")
		mute()
		blank()
	}))
	fmt.Printf("Starting server on port: %s \n", servport)
	log.Fatal(http.ListenAndServe(":"+servport, nil))
}
