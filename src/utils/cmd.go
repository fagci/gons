package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

func RunCommand(command string, wg *sync.WaitGroup, timeout time.Duration) {
	cmd := exec.Command("sh", "-c", command)

	var stderr bytes.Buffer
	var stdout bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	cmd.Start()

	done := make(chan error)

	go func() {
		done <- cmd.Wait()
	}()

	t := time.After(timeout)
	select {
	case <-t:
		cmd.Process.Kill()
		os.Stderr.WriteString(fmt.Sprintln("Cmd '" + command + "' timeout"))
	case err := <-done:
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintln("Cmd '"+command+"' run failed:", err))
			os.Stderr.WriteString(fmt.Sprintln("Out: '" + stdout.String()))
			os.Stderr.WriteString(fmt.Sprintln("Err: '" + stderr.String()))
		}
	}

	wg.Done()
}
