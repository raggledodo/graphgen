package main

import (
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

func runCmd(cmd []string) {
	log.Infof("Running command %v", cmd)
	command := exec.Command(cmd[0], cmd[1:]...)
	msg, err := command.Output()
	if err != nil {
		log.Errorf("Command %v error: %v", cmd, err)
		msg, _ = exec.Command("cat", cmd[1]).Output()
		fmt.Println(string(msg))
	} else {
		msglines := strings.FieldsFunc(string(msg),
			func(c rune) bool { return c == '\n' })
		log.Info("Command output:")
		for _, line := range msglines {
			log.Info("    ", line)
		}
	}
}
