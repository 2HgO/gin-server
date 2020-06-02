package main

import (
	"log"

	gocron "github.com/robfig/cron/v3"
)

func init() {
	job := gocron.New()
	job.AddFunc("@midnight", func() {
		log.Println("I am a crom Job")
	})
	job.Start()
}
