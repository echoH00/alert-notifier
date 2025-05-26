package config

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

const (
	MinWorkers     = 5
	MaxWorkers     = 10
	MaxRetry       = 3
	//IdleTimeout    = 5 * time.Minute
        IdleTimeout = 2 * time.Hour
	JobQueueSize   = 100
	DeadLetterSize = 100
)

type Contact struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type AlertConfig struct {
	Contacts []Contact `json:"contacts"`
}

var CurrentCfg AlertConfig

func Loadconfig(path string) (err error) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed read config.json, Err: %v", err)
		return err
	}
	err = json.Unmarshal(file, &CurrentCfg)
	if err != nil {
		log.Printf("解析配置文件失败, Err: %v", err)
		return err
	}
        for _, users := range CurrentCfg.Contacts {
		log.Printf("sendTo %v,%v",users.Name, users.Phone)
	}
	return nil
}

