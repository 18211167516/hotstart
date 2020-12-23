// +build windows

package main

import (
	"net/http"
	"time"
)

func initServer(address string, Router http.Handler) server {

	return &http.Server{
		Addr:           address,
		Handler:        Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
