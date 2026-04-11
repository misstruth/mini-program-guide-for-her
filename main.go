package main

import (
	"log"
	"net/http"
	"wxcloudrun-golang/service"
)

func main() {
	http.HandleFunc("/", service.HomeHandler)
	http.HandleFunc("/api/guide", service.GuideHandler)

	log.Fatal(http.ListenAndServe(":80", nil))
}
