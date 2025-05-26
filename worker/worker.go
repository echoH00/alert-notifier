package worker

import (
	"github.com/echoH00/alert-notifier/alertstore"
	"github.com/echoH00/alert-notifier/config"
	"github.com/echoH00/alert-notifier/sms"
	"log"
	"strings"
	"time"
)

func startWorker(id int) {
	go func() {
		log.Printf("Worker %d started", id)

		for {
			select {
			case job := <-JobQueue():
				TrackActivity(id)
				processJob(job, id)

			default:
				// 如果满足退出条件，则退出
				if CanExitWorker(id) {
					log.Printf("Worker %d exiting due to inactivity", id)
					return
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()
}

func GenSendmsg(job alertstore.AlertJob) string {
	desp := job.Alert.Annotations["summary"]
	desp = strings.ReplaceAll(desp, "\n", "")
	loc, _ := time.LoadLocation("Asia/Shanghai")
	start := job.Alert.StartsAt.In(loc).Format("0102/15:04")
	var end string
	if job.Alert.Status == "resolved" {
		end = job.Alert.EndsAt.In(loc).Format("0102/15:04")
		return "[恢复]" + desp + " " + end
	}
	return start + " " + desp
}

func processJob(job alertstore.AlertJob, workerID int) {
	content := GenSendmsg(job)
	for _, contact := range config.CurrentCfg.Contacts {
		//for _, p := range config.NotifyPhones {
		var err error
		success := false

		for attempt := 1; attempt <= config.MaxRetry; attempt++ {
			err = sms.SendMsg(contact.Phone, content)
			if err == nil {
				log.Printf("[Worker%v] send [%v] to [%s] successfully", workerID, content, contact.Name)
				success = true
				break
			}
			log.Printf("Worker %d failed to send SMS to %s (attempt %d): %v", workerID, contact.Name, attempt, err)
			time.Sleep(3 * time.Second)
		}
		if !success {
			log.Printf("[Worker%v] failed send to [%v] After 3 retries Move %v to DeadlettterQueue", workerID, contact.Name, job)
			select {
			case DeadLetterQueue() <- job:
			default:
				log.Println("Deadletter queue full, dropping job")
			}
		}
	}
}

