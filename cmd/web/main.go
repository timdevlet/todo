package main

import (
	"runtime/debug"

	log "github.com/sirupsen/logrus"
	"github.com/timdevlet/todo/internal/configs"
	"github.com/timdevlet/todo/internal/web"
)

func initLog(format string) {
	log.SetLevel(log.DebugLevel)
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("exception: %s", string(debug.Stack()))
		}
	}()

	opt := configs.NewConfigsFromEnv()

	initLog(opt.LOG_FORMAT)

	// Loop optDebug and print all fields.
	for key, value := range opt.GetFieldsWithValues() {
		log.Debug(key + ": " + value)
	}

	//

	w := web.NewWeb(opt)
	w.Init()
	w.Run()
}
