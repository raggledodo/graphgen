package main

import (
	"time"
)

type CronMap interface {
	Run()
	AddCmd(string, func() []string, time.Duration)
	Stop(string)
	StopAll()
}

type CmdMap struct {
	jobs    map[string]*Job
	running bool
}

type Job struct {
	term chan struct{}
	run  func()
}

func NewCronMap() CronMap {
	return &CmdMap{jobs: make(map[string]*Job)}
}

func (cmds *CmdMap) Run() {
	cmds.running = true
	for _, job := range cmds.jobs {
		go job.run()
	}
}

func (cmds *CmdMap) AddCmd(key string, getCmds func() []string, freq time.Duration) {
	term := make(chan struct{})
	run := func() {
		tock := time.Tick(freq)
		for {
			select {
			case <-tock:
				cmd := getCmds()
				runCmd(cmd)
			case <-term:
				return
			}
		}
	}
	if cmds.running {
		go run()
	}
	cmds.jobs[key] = &Job{
		term: term,
		run:  run,
	}
}

func (cmds *CmdMap) Stop(key string) {
	cmds.jobs[key].term <- struct{}{}
	delete(cmds.jobs, key)
}

func (cmds *CmdMap) StopAll() {
	for _, job := range cmds.jobs {
		job.term <- struct{}{}
	}
	cmds.jobs = make(map[string]*Job)
}
