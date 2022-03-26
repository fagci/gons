package utils

import (
	"fmt"
	"os/exec"
	"sync"
)

func RunCommand(cmd string, wg *sync.WaitGroup) {
	if out, err := exec.Command("sh", "-c", cmd).CombinedOutput(); err != nil {
		fmt.Println("Cmd '", cmd, "' run failed:", err)
		fmt.Println(string(out))
	}
	wg.Done()
}
