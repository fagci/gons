package utils

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

func RunCommand(cmd string, wg *sync.WaitGroup) {
	if out, err := exec.Command("sh", "-c", cmd).CombinedOutput(); err != nil {
		os.Stderr.WriteString(fmt.Sprintln("Cmd '"+cmd+"' run failed:", err))
		os.Stderr.WriteString(string(out))
	}
	wg.Done()
}
