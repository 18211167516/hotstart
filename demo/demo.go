package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func hello(w http.ResponseWriter, r *http.Request) {
	time.Sleep(20 * time.Second)
	w.Write([]byte("hello world233333!!!!"))
}

func test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test 22222"))
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/test", test)
	pid := os.Getpid()
	address := ":9999"
	s := &http.Server{
		Addr:    address,
		Handler: nil,
	}
	err := ListenAndServer(s)
	log.Printf("process with pid %d stoped, error: %s.\n", pid, err)
}
