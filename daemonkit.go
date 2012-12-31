/*
 *	Daemon Kit
 *	----------
 * 	Copyright (c) 2013, Scott Cagno, All rights reserved. 
 *	Use of this source code is governed by a BSD-style
 *	license that can be found in the LICENSE file.
 */

package daemonkit

import (
 	"io/ioutil"
	"strings"
	"strconv"
	"syscall"
	"time"
	"fmt"
	"os"
)

// daemonizer
type Daemonizer struct {
	tempf	string
	psdef	os.ProcAttr
}

// return a new daemonizer instance
func NewDaemonizer(tempf string) *Daemonizer {
	if !strings.HasSuffix(tempf, "/") {
		tempf = tempf+"/"
	}
	return &Daemonizer{tempf: tempf}
}

// watch command line arguments
func (d *Daemonizer) WatchCli(args []string) {
	if len(args) < 3 {
		fmt.Printf("usage: %s {start|sample|stop|restart} prog ...args\n", args[0])
		os.Exit(1)
	}
	actn := args[1]
	prgm := args[2]
	othr := args[3:]
	switch actn {
	case "start":
		d.Start(prgm, othr)
	case "stop":
		d.Stop(prgm, othr)
	case "restart":
		d.Stop(prgm, othr)
		defer d.Start(prgm, othr)
	case "sample":
		d.Sample(prgm)
	default:
		fmt.Printf("usage: %s {start|sample|stop|restart} prog ...args\n", args[0])
		os.Exit(1)
	}
}

// start daemon process
func (d *Daemonizer) Start(prog string, args []string) {
	dn, _ := os.Open(os.DevNull)
	d.psdef.Files = []*os.File{dn, dn, dn}
	ps, err := os.StartProcess(prog, args, &d.psdef)
	if err != nil {
		vomit(err)
	}
	d.WPidFile(prog, ps.Pid)
	ps.Release()
	fmt.Printf("[START] '%s' 	-- OK!\n", prog)
}

// stop daemon process
func (d *Daemonizer) Stop(prog string, args []string) {
	fileData := d.RPidFile(prog)
	pid, _ := strconv.Atoi(fileData[0])
	ps, err := os.FindProcess(pid)
	if err != nil {
		vomit(err)
	}
	err = ps.Signal(syscall.SIGINT)
	if err != nil {
		vomit(err)
	}
	fileName := fmt.Sprintf("%s%s.pid", d.tempf, prog)
	err = os.Remove(fileName)
	if err != nil {
		vomit(err)
	}
	fmt.Printf("[STOP] '%s' 	-- OK!\n", prog)
}

// sample daemon process
func (d *Daemonizer) Sample(prog string) {
	fileData := d.RPidFile(prog)
	pid, _ := strconv.Atoi(fileData[0])
	ps, err := os.FindProcess(pid)
	if err != nil {
		vomit(err)
	}
	ps.Release()
	if len(fileData) == 3 {
		t, _ := strconv.ParseInt(fileData[2], 0, 64)
		elapsed := timeElapsed(t)
		fmt.Printf("\n  ----------\n  PRG: %s\n  ----------\n  PID: %s\n  TMP: %s\n  UPT: %s\n\n", prog, fileData[0], fileData[1], elapsed)
	}
}

// write tmp file containing program pid and program info
func (d *Daemonizer) WPidFile(prog string, pid int) {
	fileName := fmt.Sprintf("%s%s.pid", d.tempf, prog)
	fileData := fmt.Sprintf("%d,%s,%d", pid, fileName, time.Now().Unix())
	err := ioutil.WriteFile(fileName, []byte(fileData), 0644)
	if err != nil {
		vomit(err)
	}
}

// read tmp file containing program pid and program info
func (d *Daemonizer) RPidFile(prog string) []string {
	fileName := fmt.Sprintf("%s%s.pid", d.tempf, prog)
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		vomit(err)
	}
	fileData := strings.Split(string(data), ",")
	if len(fileData) == 2 {
		fileData = append(fileData, "unknown")
	}
	return fileData
}

// error handler
func vomit(err error) {
	fmt.Printf("**%s\n", err)
	os.Exit(1)
}

// return rounded time elapsed, string
func timeElapsed(t1 int64) string {
	t := time.Now().Unix()-t1
	if t < 1 {
		return fmt.Sprintf("0 seconds")
	}
	if t == 1 {
		return fmt.Sprintf("1 second")
	}
	if t > 1 && t < 60 {
		return fmt.Sprintf("%d seconds", t)
	}
	if t == 60 {
		return fmt.Sprintf("1 minute")
	}
	if t > 60 && t < 3600 {
		return fmt.Sprintf("%d minutes", t/60)
	}
	if t == 3600 {
		return fmt.Sprintf("1 hour")
	}
	if t > 3600 && t < 86400 {
		return fmt.Sprintf("%d hours | %d minutes", t/3600, t/60)
	}
	if t == 86400 {
		return fmt.Sprintf("1 day")
	}
	if t > 86400 {
		return fmt.Sprintf("%d days | %d hours", t/86400, t/3600)
	}
	return ""
}
