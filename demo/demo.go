package main

import (
	"log"
	"net/http"
	"os"
	"time"

	Hot "github.com/18211167516/hotstart"
)

func hello(w http.ResponseWriter, r *http.Request) {
	time.Sleep(20 * time.Second)
	w.Write([]byte("hello world233333!!!!"))
}

func test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/test", test)
	pid := os.Getpid()
	address := ":9999"
	err := Hot.ListenAndServe(address, nil)
	log.Printf("process with pid %d stoped, error: %s.\n", pid, err)
}
