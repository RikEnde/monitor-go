package graph

import (
	"fmt"
	db "monitor/storage"
	"time"
)

type Point struct {
	x int
	y int
}

func (p Point) String() string {
	return fmt.Sprintf("%d,%d", p.x, p.y)
}

func String(Points []Point) []string {
	l := make([]string, len(Points))
	for i, p := range Points {
		l[i] = p.String()
	}
	return l
}

func Graph(width int, height int) ([]Point, []Point, []Point, []Point) {
	to := time.Now()
	from := to.Add(-24 * time.Hour)
	list := db.ReadDiffies(from, to)

	user := make([]Point, 1)
	iowait := make([]Point, 1)
	temperature := make([]Point, 1)

	t0 := list[0].Timestamp.Unix()
	tf := list[len(list)-1].Timestamp.Unix()
	d := int(tf - t0)

	p0 := Point{x: 0, y: height}
	user[0], iowait[0], temperature[0] = p0, p0, p0

	xold := 0
	for _, diff := range list {
		xi := width * int(diff.Timestamp.Unix()-t0) / d
		if xi-xold > 10 {
			p1 := Point{x: xold, y: height}
			user, iowait, temperature = append(user, p1), append(iowait, p1), append(temperature, p1)
			p2 := Point{x: xi, y: height}
			user, iowait, temperature = append(user, p2), append(iowait, p2), append(temperature, p2)
		}
		user = append(user, Point{x: xi, y: height - int(diff.User/100.0*float64(height))})
		iowait = append(iowait, Point{x: xi, y: height - int(diff.Iowait/100.0*float64(height))})
		temperature = append(temperature, Point{x: xi, y: height - int(diff.Temperature/100.0*float64(height))})
		xold = xi
	}

	pf := Point{x: width, y: height}
	user, iowait, temperature = append(user, pf), append(iowait, pf), append(temperature, pf)

	box := []Point{Point{x: 0, y: 0}, Point{x: width, y: 0}, Point{x: width, y: height}, Point{x: 0, y: height}}

	return box, user, temperature, iowait
}
