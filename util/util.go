package util

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
#include <unistd.h>
#include <sys/types.h>
#include <pwd.h>
#include <stdlib.h>
*/
import "C"

func UserHZ() int {
	var sc_clk_tck C.long
	sc_clk_tck = C.sysconf(C._SC_CLK_TCK)
	return int(sc_clk_tck)
}

func Timestamp() int {
	return int(time.Now().Unix())
}

func Trim(str string) string {
	return strings.Trim(str, " \t\r\n")
}

func ReadFile(name string) []string {
	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	lines := make([]string, 10)
	in := bufio.NewScanner(file)
	for in.Scan() {
		lines = append(lines, in.Text())
	}
	return lines
}

func Atoi(str []string) []int {
	arr := make([]int, 0, len(str))
	for _, s := range str {
		i, err := strconv.Atoi(s)
		if err != nil {
			/* Don't second guess the kernel, just replace with 0 */
			i = 0
		}
		arr = append(arr, i)
	}

	return arr
}
