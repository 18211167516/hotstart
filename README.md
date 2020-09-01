# HotStart

[![godoc](https://camo.githubusercontent.com/a613b85cc4087229d2655adedfcd8fb2580dc56a/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f726f796c6565303730342f67726f6e3f7374617475732e737667)](https://godoc.org/github.com/18211167516/hotstart)

# Installation

```golang
go get github.com/18211167516/hotstart
```

# Usage

```golang
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
	err := Hot.ListenAndServer(s)
	log.Printf("process with pid %d stoped, error: %s.\n", pid, err)
}

```
