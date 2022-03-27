package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

func RunCommand(command string, wg *sync.WaitGroup, timeout time.Duration) {
    shell := "sh"
    opt := "-c"

    if runtime.GOOS == "windows" {
        shell = "cmd.exe"
        opt = "/C"
    }

	cmd := exec.Command(shell, opt, command)

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
		os.Stderr.WriteString(fmt.Sprintln("[W] Cmd '" + command + "' timeout"))
	case err := <-done:
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintln("[W] Cmd '"+command+"' run failed:", err))
			os.Stderr.WriteString(fmt.Sprintln("[W] Out: " + stdout.String()))
			os.Stderr.WriteString(fmt.Sprintln("[W] Err: " + stderr.String()))
		}
	}

	wg.Done()
}
