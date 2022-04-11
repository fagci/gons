package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"sync"
	"time"
)

type Flags uint8

const (
	ERR Flags = 1 << iota
	WARN
	INFO
)

var forbiddenInQuotesCharsRegexp = regexp.MustCompile("[$\"'`]")

func (b Flags) Set(flag Flags) Flags    { return b | flag }
func (b Flags) Clear(flag Flags) Flags  { return b &^ flag }
func (b Flags) Toggle(flag Flags) Flags { return b ^ flag }
func (b Flags) Has(flag Flags) bool     { return b&flag != 0 }

func RunCommand(command string, wg *sync.WaitGroup, timeout time.Duration, flags Flags) {
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

	select {
	case <-time.After(timeout):
		cmd.Process.Kill()
		if flags.Has(WARN) {
			EPrintln("[W:CB:timeout]", "'"+command+"'")
		}
	case err := <-done:
		if err != nil && flags.Has(ERR) {
			EPrintln("[W:CB:E]", "'"+command+"'", err)
			if stderr.Len() != 0 {
				EPrintf("[i:CB:err] %s", stderr.String())
			}
		}
		if flags.Has(INFO) && stdout.Len() != 0 {
			fmt.Print(stdout.String())
		}
	}
}

func FilterValueInQuotes(v string) string {
	return forbiddenInQuotesCharsRegexp.ReplaceAllString(v, "")
}
