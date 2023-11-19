package main

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

// SysInfo will record cpu and memory data
type SysInfo struct {
	CPU    float64
	Memory float64
}

var eol string

func init() {
	eol = "\n";
}
 

func formatStdOut(stdout []byte, userfulIndex int) []string {
	infoArr := strings.Split(string(stdout), eol)[userfulIndex]
	ret := strings.Fields(infoArr)
	return ret
}

func parseFloat(val string) float64 {
	floatVal, _ := strconv.ParseFloat(val, 64)
	return floatVal
}

func stat(pid int, statType string) (*SysInfo, error) {
	sysInfo := &SysInfo{}
	args := "-o pcpu,rss -p"
	stdout, _ := exec.Command("ps", args, strconv.Itoa(pid)).Output()
	ret := formatStdOut(stdout, 1)
	if len(ret) == 0 {
		return sysInfo, errors.New("Can't find process with this PID: " + strconv.Itoa(pid))
	}
	sysInfo.CPU = parseFloat(ret[0])
	sysInfo.Memory = parseFloat(ret[1]) * 1024
	return sysInfo, nil
}

func wrapper(statType string) func(pid int) (*SysInfo, error) {
	return func(pid int) (*SysInfo, error) {
		return stat(pid, statType)
	}
}


// GetStat will return current system CPU and memory data
func GetStat(pid int) (*SysInfo, error) {
	sysInfo, err := wrapper("ps")(pid)
	return sysInfo, err
}