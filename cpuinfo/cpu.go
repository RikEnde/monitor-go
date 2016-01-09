package cpuinfo

import (
	"fmt"
	"io/ioutil"
	u "monitor/util"
	"strconv"
	"strings"
)

type KV map[string]string

func readAsInt(path string) (int, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	value, err := strconv.ParseFloat(u.Trim(string(file)), 64)
	if err != nil {
		return 0, err
	}
	return int(0.5 + value/1000.0), nil
}

func MaxFreq() int {
	/* Assuming symmetric multiprocessing, all cpus will run at the same frequency, and cpu0 will always be present */
	max, err := readAsInt("/sys/devices/system/cpu/cpu0/cpufreq/scaling_max_freq")
	if err != nil {
		panic(err)
	}
	return max
}

func CpuFreq() int {
	/* Assuming symmetric multiprocessing, all cpus will run at the same frequency, and cpu0 will always be present */
	freq, err := readAsInt("/sys/devices/system/cpu/cpu0/cpufreq/scaling_cur_freq")
	if err != nil {
		panic(err)
	}
	return freq
}

func CpuInfo() []KV {
	file, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		panic(err)
	}
	cpus := strings.Split(u.Trim(string(file)), "\n\n")

	var maps = make([]KV, len(cpus))
	for id, cpu := range cpus {

		cpuinfo := append(strings.Split(cpu, "\n"), fmt.Sprintf("timestamp:%d", u.Timestamp()))

		var cpuMap = make(KV, len(cpuinfo))
		for _, str := range cpuinfo {
			chunks := strings.Split(str, ":")
			if len(chunks) == 2 {
				cpuMap[u.Trim(strings.ToLower(chunks[0]))] = u.Trim(chunks[1])
			}
		}
		maps[id] = cpuMap
	}

	return maps
}
