package logs

import (
	"encoding/json"
	"time"
)

type Log struct {
	Message    string                 `json:"message"`
	IP         string                 `json:"ip"`
	Host       string                 `json:"host"`
	Method     string                 `json:"method"`
	RequestURI string                 `json:"request_uri"`
	Status     int                    `json:"status"`
	Request    map[string]interface{} `json:"request"`
	Agent      string                 `json:"agent"`
	Referer    string                 `json:"referer"`
	Start      time.Time              `json:"start"`
	Stop       time.Time              `json:"stop"`
	ID         string                 `json:"id"`
}

func (l *Log) ToJson() string {
	e, err := json.Marshal(l)
	if err != nil {
		return ""
	}

	return string(e)
}
