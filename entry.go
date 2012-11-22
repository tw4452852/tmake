//	entry is a minimal work cell

package main

import (
	"os/exec"
	"bytes"
	"fmt"
	"strings"
	"log"
	"regexp"
	"io"
)

type entry struct {
	procName	string			//program name
	procArgs	string			//program args
	exitCode	error			//exit error if has
	errBuffer	*bytes.Buffer	//save error info
}

func NewEntry(cmd string) (*entry, error) {
	fields := strings.Fields(strings.TrimSpace(cmd))
	if len(fields) == 0 {
		return nil, fmt.Errorf("cmdline is null")
	}
	procName := fields[0]
	var procArgs string
	if len(fields) > 1 {
		procArgs = strings.Join(fields[1:], " ")
	}
	return &entry{
		procName: procName,
		procArgs: procArgs,
	}, nil
}

func (e *entry) ExitCode() error {
	return e.exitCode
}

func (e *entry) ErrString() string {
	if e.errBuffer == nil {
		return ""
	}
	return e.errBuffer.String()
}

func (e *entry) Run() error {
	//sanity check
	if e.procName == "" {
		return fmt.Errorf("cmdline is null")
	}

	cmd := exec.Command(e.procName, e.procArgs)
	allOut := new(bytes.Buffer)
	//let all output to one buffer
	cmd.Stdout = allOut
	cmd.Stderr = allOut

	e.exitCode = cmd.Run()
	if e.exitCode != nil {
		if e.errBuffer == nil {
			e.errBuffer = new(bytes.Buffer)
		}
		getErrorInfo(e.errBuffer, allOut)
	}
	return e.exitCode
}

var errPattern = regexp.MustCompile(`(?i:err|fail)`)

func getErrorInfo(errInfo *bytes.Buffer, all *bytes.Buffer) {
	var prevLine string
	for {
		line, err := all.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("parse line err:", err)
			continue
		}
		if errPattern.MatchString(line) {
			errInfo.WriteString(prevLine)
			errInfo.WriteString(line)
		}
		prevLine = line
	}
}
