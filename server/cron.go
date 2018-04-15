package main

import (
	"os/exec"

	"github.com/robfig/cron"

	log "github.com/sirupsen/logrus"
)

var tfGenCmd = []string{"./tfgen", "--host", "127.0.0.1:50051"}

func NewCron(schedule cron.Schedule) *cron.Cron {
	job := cron.New()
	job.Schedule(schedule, cron.FuncJob(func() {
		log.Infof("Running command %v", tfGenCmd)
		cmd := exec.Command(tfGenCmd[0], tfGenCmd[1:]...)
		msg, err := cmd.Output()
		if err != nil {
			log.Errorf("Command %v error: %v", tfGenCmd, err)
		} else {
			log.Infof("Command output: %s", string(msg))
		}
	}))
	return job
}
