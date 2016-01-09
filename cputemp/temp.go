package cputemp

import (
	"io/ioutil"
	"log"
	u "monitor/util"
	"strconv"
)

func CpuTemp() float64 {
	file, err := ioutil.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		panic(err)
	}
	temp, err := strconv.Atoi(u.Trim(string(file)))
	if err != nil {
		log.Fatal(err)
	}

	return float64(temp) / 1000.0
}
