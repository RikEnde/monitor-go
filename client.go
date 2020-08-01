package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	u "monitor/util"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

func readInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", prompt)
	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	return u.Trim(text)
}

func main() {
	port := flag.String("port", "8080", "TCP port to listen to")
	host := flag.String("host", "localhost", "Host to connect to")
	user := flag.String("user", "user", "Username")
	password := flag.String("password", "", "Password")
	flag.Parse()
	if "" == *password {
		*password = readInput("password")
	}
	tail := strings.Join(flag.Args(), "/")

	var url string
	if tail != "" {
		url = fmt.Sprintf("https://%s:%s/%s/", *host, *port, tail)
	} else {
		url = fmt.Sprintf("https://%s:%s/stat/cpu/", *host, *port)
	}

	if *port != "8080" || *host != "localhost" {
		log.Printf("%s\nConnecting to %s\n", runtime.GOROOT(), url)
	}

	// We're going to accept self signed certificates
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr, Timeout: 10 * time.Second}

	request, err := http.NewRequest("GET", url, nil)
	request.SetBasicAuth(*user, *password)

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("can't read response body: %s", err)
	}
	fmt.Printf("%s", string(body))
}
