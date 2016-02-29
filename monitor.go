package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	c "monitor/cpuinfo"
	tmp "monitor/cputemp"
	g "monitor/graph"
	s "monitor/stat"
	db "monitor/storage"
	t "monitor/template"
	u "monitor/util"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	width  = 800
	height = 400
)

func basicAuth(w http.ResponseWriter, r *http.Request) bool {
	denied := func() bool {
		w.Header().Set("WWW-Authenticate", `Basic realm="Beware! Protected REALM! "`)
		w.WriteHeader(401)
		w.Write([]byte("401 Unauthorized\n"))
		return false
	}
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 {
		return denied()
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return denied()
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return denied()
	}

	if !u.Login(pair[0], pair[1]) {
		log.Printf("Access denied for %s", pair[0])
		return denied()
	}

	log.Printf("User %s logged in", pair[0])

	return true
}

type Command struct {
	Please string
}

func command(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)

	var comm Command
	err := dec.Decode(&comm)
	if err != nil || comm.Please == "" {
		w.Header().Set("WWW-Authenticate", `Basic realm="Beware! Protected REALM! "`)
		w.WriteHeader(418)
		fmt.Fprintf(w, "418 Insufficiently polite")
		return
	}

	log.Printf("Command parsed: %v", comm)

	switch comm.Please {
	case "Die!":
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "By your command, master")
		os.Exit(0)
	case "Restart":
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "See you later")
		os.Exit(1)
	default:
		w.Header().Set("WWW-Authenticate", `Basic realm="Beware! Protected REALM! "`)
		w.WriteHeader(418)
		fmt.Fprintf(w, "418 Stupid command")
	}

}

func index(w http.ResponseWriter, r *http.Request) {
	if !basicAuth(w, r) {
		return
	}
	if r.Method == "POST" {
		command(w, r)
		return
	}

	title := "CPU Status REST service"

	body := []t.Link{
		t.Link{Url: "/cpu", Desc: "query cpuinfo"},
		t.Link{Url: "/cpu/0/cpu%20mhz", Desc: "query mHz of cpu 0"},
		t.Link{Url: "/stat", Desc: "query cpustat"},
		t.Link{Url: "/stat/cpu", Desc: "query jiffies aggregate over all cores"},
		t.Link{Url: "/history", Desc: "query diffies over past hour"},
		t.Link{Url: "/temperature", Desc: "query cpu core temperature"},
		t.Link{Url: "/frequency", Desc: "query cpu clock frequency"},
		t.Link{Url: "/graph", Desc: "graph diffies over past 24 hours"},
	}

	p := t.MakePage(title, body, fmt.Sprintf("%s", r))

	t.RenderTemplate(w, "index", p)
}

func reply(w http.ResponseWriter, value interface{}) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", u.JsonResponse(value))
}

func reject(w http.ResponseWriter, value interface{}) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "bad request")
	log.Printf("err: %s", value)
}

func cpu(w http.ResponseWriter, r *http.Request) {
	if !basicAuth(w, r) {
		return
	}

	log.Printf("query: %v\n", r.URL.Path)
	path := strings.Trim(r.URL.Path[len("/cpu/"):], "/")

	info := c.CpuInfo()
	if path != "" {
		query := strings.Split(path, "/")
		cpuid, err := strconv.Atoi(query[0])
		if err != nil || cpuid >= len(info) {
			reject(w, err)
			return
		}

		cpu := info[cpuid]
		if len(query) >= 2 {
			if key, ok := cpu[query[1]]; ok {
				reply(w, key)
			} else {
				reject(w, query[1])
			}
		} else {
			reply(w, cpu)
		}
	} else {
		reply(w, info)
	}

}

func stat(w http.ResponseWriter, r *http.Request) {
	if !basicAuth(w, r) {
		return
	}
	query := strings.Trim(r.URL.Path[len("/stat/"):], "/")
	log.Printf("query: %v\n", query)

	stat := s.Stat()
	if len(query) > 0 {
		if key, ok := stat[query]; ok {
			reply(w, key)
		} else {
			reject(w, query)
		}
	} else {
		reply(w, stat)
	}
}

func history(w http.ResponseWriter, r *http.Request) {
	if !basicAuth(w, r) {
		return
	}
	to := time.Now()
	from := to.Add(-1 * time.Hour)
	list := db.ReadDiffies(from, to)
	reply(w, list)
}

func temp(w http.ResponseWriter, r *http.Request) {
	if !basicAuth(w, r) {
		return
	}
	reply(w, tmp.CpuTemp())
}

func freq(w http.ResponseWriter, r *http.Request) {
	if !basicAuth(w, r) {
		return
	}
	m := make(map[string]int, 2)
	freq, max := c.CpuFreq(), c.MaxFreq()
	m["frequency"] = freq
	m["maximum"] = max

	reply(w, m)
}

func graph(w http.ResponseWriter, r *http.Request) {
	if !basicAuth(w, r) {
		return
	}
	w.WriteHeader(http.StatusOK)

	box, user, temperature, iowait := g.Graph(width, height)

	fmt.Fprintf(w, "<html><head></head><body><svg xmlns='http://www.w3.org/2000/svg' style='stroke: grey; fill: white; stroke-width: 0.5' width='%d' height='%d'>\n", width, height)
	fmt.Fprintf(w, "<polygon Points='%s' style='stroke:black; fill: none; stroke-width: 0.5'/>\n", strings.Join(g.String(box), " "))
	fmt.Fprintf(w, "<polygon Points='%s' style='stroke:#800000; fill: red; stroke-width: 0.5'/>\n", strings.Join(g.String(user), " "))
	fmt.Fprintf(w, "<polygon Points='%s' style='stroke:black; fill: none; stroke-width: 0.5'/>\n", strings.Join(g.String(temperature), " "))
	fmt.Fprintf(w, "<polygon Points='%s' style='stroke:#000800; fill: green; stroke-width: 0.5'/>\n", strings.Join(g.String(iowait), " "))
	fmt.Fprintf(w, "</svg></body></html>\n")
}

func startServer(ch chan string, port string) {
	log.Printf("Listening on port: %s with goroot: %s\n", port, runtime.GOROOT())

	http.HandleFunc("/", index)
	http.HandleFunc("/cpu/", cpu)
	http.HandleFunc("/stat/", stat)
	http.HandleFunc("/temperature/", temp)
	http.HandleFunc("/frequency/", freq)
	http.HandleFunc("/history/", history)
	http.HandleFunc("/graph", graph)

	err := http.ListenAndServeTLS(":"+port, u.Cert(), u.Key(), nil)
	if err != nil {
		panic(err)
	}

}

func main() {
	port := flag.String("port", "8080", "TCP port to listen to")
	flag.Parse()

	ch := make(chan string)

	go startServer(ch, *port)
	go db.StoreDiffies(ch)

	for s := range ch {
		log.Printf(s)
	}
}
