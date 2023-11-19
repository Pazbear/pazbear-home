package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main(){
	myPid, _ := strconv.Atoi(os.Args[1])
 
    for {
        sysInfo, _ := GetStat(myPid)
        fmt.Println("CPU Usage     :", sysInfo.CPU)
        fmt.Println("Mem Usage(RSS):", sysInfo.Memory)
        time.Sleep(5 * time.Second)
    }
}