package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	i "monitor/cpuinfo"
	t "monitor/cputemp"
	s "monitor/stat"
	u "monitor/util"
	"time"
)

const SECONDS = 30

type CpuInfo struct {
	User        float64
	Nice        float64
	System      float64
	Idle        float64
	Iowait      float64
	Irq         float64
	Softirq     float64
	Steal       float64
	Freq        float64
	Temperature float64
	Timestamp   time.Time
}

func DbInfo() string {
	cred := u.DbCredentials()
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", cred["user"], cred["password"], cred["database"])
}

func diffies(s s.CpuStat, l s.CpuStat) (float64, float64, float64, float64, float64, float64, float64) {
	delta := s.Timestamp - l.Timestamp
	HZ := float64(u.UserHZ())
	// convert to percent
	f := func(i int) float64 {
		return 100.0 * float64(i) / (float64(delta) * float64(s.Cores) * HZ)
	}

	return f(s.User - l.User), f(s.Nice - l.Nice), f(s.System - l.System), f(s.Idle - l.Idle), f(s.Iowait - l.Iowait), f(s.Irq - l.Irq), f(s.Softirq - l.Softirq)
}

func StoreDiffies(ch chan string) {
	db, err := sql.Open("postgres", DbInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO cpuinfo(time, userland, nice, system, idle, iowait, irq, softirq, steal, freq, temperature) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)")
	defer stmt.Close()

	ticker := time.NewTicker(time.Millisecond * 1000 * SECONDS)
	last := s.Cpu()
	for range ticker.C {
		stat := s.Cpu()
		temp := t.CpuTemp()
		freq := 100.0 * float64(i.CpuFreq()) / float64(i.MaxFreq())
		user, nice, system, idle, iowait, irq, softirq := diffies(stat, last)
		res, err := stmt.Exec(time.Unix(int64(stat.Timestamp), 0), user, nice, system, idle, iowait, irq, softirq, 0, freq, temp)
		if err != nil || res == nil {
			log.Fatal(err)
			break
		}
		ch <- fmt.Sprintf("time %d diff %2d user %2.1f nice %2.1f system %2.1f idle %2.1f iowait %2.1f irq %2.1f softirq %2.1f temperature %2.1f frequency %4.2f\n",
			stat.Timestamp, stat.Timestamp-last.Timestamp, user, nice, system, idle, iowait, irq, softirq, temp, freq)
		last = stat
	}
}

func ReadDiffies(from time.Time, to time.Time) []CpuInfo {
	db, err := sql.Open("postgres", DbInfo())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT time, userland, nice, system, idle, iowait, irq, softirq, steal, freq, temperature FROM cpuinfo WHERE time BETWEEN $1 AND $2 ORDER BY time ASC", from, to)
	if err != nil {
		panic(err)
	}

	list := make([]CpuInfo, 0)

	for rows.Next() {
		stat := new(CpuInfo)
		err = rows.Scan(&stat.Timestamp, &stat.User, &stat.Nice, &stat.System, &stat.Idle, &stat.Iowait, &stat.Irq, &stat.Softirq, &stat.Steal, &stat.Freq, &stat.Temperature)
		if err != nil {
			panic(err)
		}
		list = append(list, *stat)
	}
	return list
}
