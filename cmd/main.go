package main

import (
	"flag"
	"github.com/echoH00/alert-notifier/config"
	"github.com/echoH00/alert-notifier/server"
	"github.com/echoH00/alert-notifier/worker"
	"log"
	"net/http"
)

func main() {
	configPath := flag.String("config", "config/config.json", "Path to config file")
	flag.Parse()
	err := config.Loadconfig(*configPath)
	if err != nil {
		log.Printf("Err: %v", err)
		return
	}
	log.Printf("Used config: %v", *configPath)
	worker.InitPool()
	http.HandleFunc("/webhook", server.AlertHandler)
	log.Println("Server started on :5001")
	log.Fatal(http.ListenAndServe(":5001", nil))
}

