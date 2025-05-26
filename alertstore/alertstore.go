package alertstore

import "github.com/prometheus/alertmanager/template"

type AlertJob struct {
	Alert     template.Alert
	Retry     int
	Timestamp int64
}

