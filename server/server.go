package server

import (
	"encoding/json"
	"github.com/echoH00/alert-notifier/alertstore"
	"github.com/echoH00/alert-notifier/worker"
	"github.com/prometheus/alertmanager/template"
	"log"
	"net/http"
	"time"
)

func AlertHandler(w http.ResponseWriter, r *http.Request) {
	var data template.Data
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("解析告信息警失败,Err: %v", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	log.Printf("Receive: %s from %s", r.URL.Path, r.RemoteAddr)

	hasResolved := false
	for _, alert := range data.Alerts {
		if alert.Status == "resolved" {
			hasResolved = true
			break
		}
	}

	for _, a := range data.Alerts {
		// 包含恢复 只发送恢复
		if hasResolved {
			if a.Status != "resolved" {
				continue
			}
		}

		// 全部是firing
		job := alertstore.AlertJob{
			Alert:     a,
			Retry:     0,
			Timestamp: time.Now().Unix(),
		}
		worker.EnqueueJob(job)
	}
	resp, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
	log.Println("request done===")
}

