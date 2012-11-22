//	entry is a minimal work cell

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

type entry struct {
	procName  string        //program name
	procArgs  []string      //program args
	exitCode  error         //exit error if has
	errBuffer *bytes.Buffer //save error info
}

func NewEntry(cmd string) (*entry, error) {
	fields := strings.Fields(strings.TrimSpace(cmd))
	if len(fields) == 0 {
		return nil, fmt.Errorf("cmdline is null")
	}
	procName := fields[0]
	var procArgs []string
	if len(fields) > 1 {
		procArgs = fields[1:]
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
	return e.errBuffer.String()
}

func (e *entry) Run() error {
	//sanity check
	if e.procName == "" {
		return fmt.Errorf("cmdline is null")
	}

	cmd := exec.Command(e.procName, e.procArgs...)
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

func getErrorInfo(errInfo *bytes.Buffer, all *bytes.Buffer) {
	for {
		line, err := all.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("parse line err:", err)
			continue
		}
		//TODO:filter errinfo
		errInfo.WriteString(line)
	}
}
