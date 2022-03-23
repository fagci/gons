package utils

import (
	"fmt"
	"os/exec"
	"sync"
)

func RunCommand(cmd string, wg *sync.WaitGroup) {
	if _, err := exec.Command("bash", "-c", cmd).Output(); err != nil {
		fmt.Println("Cmd '", cmd, "' run failed:", err)
	}
	wg.Done()
}
