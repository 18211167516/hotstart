// +build !windows

package main

import (
	"net/http"
	"time"

	Hot "github.com/18211167516/hotstart"
)

func initServer(address string, Router http.Handler) server {

	s := Hot.NewServer(address, Router)
	s.ReadHeaderTimeout = 10 * time.Millisecond
	s.WriteTimeout = 10 * time.Second
	s.MaxHeaderBytes = 1 << 20
	return s
}
