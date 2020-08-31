package main

import (
	"net/http"
	"os"

	Hot "github.com/18211167516/hotstart"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world233333!!!!"))
}

func main() {
	http.HandleFunc("/hello", handler)
	pid := os.Getpid()
	address := ":9999"
	err := Hot.ListenAndServe(address, nil)
}
