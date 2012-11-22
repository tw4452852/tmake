package main

import (
	"testing"
)

var tests = []struct {
	cmd string
	err	string
}{
	{"echo hello",	""},
	{"ls /",		""},
	{"",			"cmdline is null"},
}

func TestEntry(t *testing.T) {
	for _, test := range tests {
		exec , err := NewEntry(test.cmd)
		if err != nil {
			if errString(err) != test.err {
				t.Errorf("%q test failed: result=%q, expected=%q\n", test.cmd, err.Error(), test.err)
			}
			continue
		}
		err = exec.Run()
		if err != nil {
			t.Log(exec.ErrString())
		}
		if errString(err) != test.err {
			t.Errorf("%q test failed: result=%q, expected=%q\n", test.cmd, err.Error(), test.err)
		}
	}
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
