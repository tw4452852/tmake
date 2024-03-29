package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	isX86  bool
	ishelp bool
)

func init() {
	flag.BoolVar(&isX86, "86", false, "build for x86")
	flag.BoolVar(&ishelp, "h", false, "show help")
}

func main() {
	flag.Parse()
	if ishelp {
		flag.Usage()
		os.Exit(1)
	}
	//get the stuff according to the args
	cmdlines := getCmdlines(flag.Args(), isX86)
	doMake(cmdlines)
}

func doMake(cmdlines []string) {
	for _, cmdline := range cmdlines {
		exec, err := NewEntry(cmdline)
		if err != nil {
			fmt.Printf("NewEntry %q failed: %s", cmdline, err)
			break
		}
		err = exec.Run()
		if err != nil {
			fmt.Printf("exec %q failed: %s: %s\n", cmdline, err, exec.ErrString())
			break
		}
		log.Printf("%s done\n", cmdline)
	}
}

type alias struct {
	normal string
	short  string
}

//order is required
var allStuffSequence = []alias{
	{"clean", "c"},   //make clean
	{"prepare", "p"}, //make prepare
	{"tw_mk ", ""},   //make
	{"strip", "s"},   //make strip
	{"install", "i"}, //make install
	{"os", "o"},      //make os
}

const (
	makecmd  = "make "
	buildX86 = "BUILD_FOR=TSERIES_X86_"
)

func getCmdlines(args []string, x86 bool) (cmdlines []string) {
	//default: make
	if len(args) == 0 {
		args = append(args, "")
	}
	cmdlines = make([]string, 0, len(args))
	for _, stuff := range allStuffSequence {
		for _, arg := range args {
			if arg == stuff.normal || arg == stuff.short {
				var cmdline string
				if arg == "" {
					//for myself makefile script for parallel
					cmdline = stuff.normal
				} else {
					cmdline = makecmd + stuff.normal
				}
				//add build x86 flag
				if x86 {
					cmdline += buildX86
				}
				cmdlines = append(cmdlines, cmdline)
			}
		}
	}
	return cmdlines
}
