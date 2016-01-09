package stat

import (
	"fmt"
	u "monitor/util"
	"strings"
)

type KV map[string][]int

type CpuStat struct {
	User        int
	Nice        int
	System      int
	Idle        int
	Iowait      int
	Irq         int
	Softirq     int
	Steal       int
	Freq        int
	Temperature int
	Timestamp   int
	Cores       int
}

func Cpu() CpuStat {
	s := Stat()
	cores := -1
	for key, _ := range s {
		if strings.HasPrefix(key, "cpu") {
			cores++
		}
	}
	cpu := s["cpu"]
	return CpuStat{
		User:      cpu[0],
		Nice:      cpu[1],
		System:    cpu[2],
		Idle:      cpu[3],
		Iowait:    cpu[4],
		Irq:       cpu[5],
		Softirq:   cpu[6],
		Timestamp: u.Timestamp(),
		Cores:     cores,
	}
}

func Stat() KV {
	stat := u.ReadFile("/proc/stat")
	stat = append(stat, fmt.Sprintf("timestamp %d", u.Timestamp()))
	m := make(KV, len(stat)+1)
	for _, str := range stat {
		fields := strings.Fields(str)
		if len(fields) >= 2 {
			m[fields[0]] = u.Atoi(fields[1:])
		}
	}
	return m
}
