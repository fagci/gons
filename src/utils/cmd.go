package utils

import (
	"bytes"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

func RunCommand(command string, wg *sync.WaitGroup, timeout time.Duration) {
	defer wg.Done()
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
		EPrintln("[W] Cmd '" + command + "' timeout")
	case err := <-done:
		if err != nil {
			EPrintln("[W] Cmd '"+command+"' run failed:", err)
		}
		if stdout.Len() != 0 {
			EPrint("[i] Out: " + stdout.String())
		}
		if stderr.Len() != 0 {
			EPrint("[i] Err: " + stderr.String())
		}
	}
}
