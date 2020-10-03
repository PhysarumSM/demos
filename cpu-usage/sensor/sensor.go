package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) != 3 {
        fmt.Println("Usage:", os.Args[0], "<proxy-port> <listen-port>")
        os.Exit(1)
    }
    proxyPort := os.Args[1]
    listenPort := os.Args[2]

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request from:", r.RemoteAddr)
		fmt.Fprintf(w, "OK")
	})

	go func(){
		log.Fatal(http.ListenAndServe(":" + listenPort, nil))
	}()

	initResp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s/cpu-usage-aggregator:1.0", proxyPort))
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(initResp.Body)
	if err != nil {
		log.Fatal(err)
	}
	initResp.Body.Close()
	log.Println("Initial request for aggregator response:", body)
	log.Println("Wait 10 seconds before start sending CPU data")
	time.Sleep(time.Second * 10)

	rand.Seed(time.Now().UnixNano())
	id := rand.Uint32()
	for {
		cmd := exec.Command("bash", "-c",
				`top -bn2 -d 0.5 | fgrep 'Cpu(s)' | tail -1 | awk  -F'id,' '{ n=split($1, vals, ","); v=vals[n]; sub("%", "", v); printf "%f", 100 - v }'`)
		output, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		cpuUtilization, err := strconv.ParseFloat(string(output), 32)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("CPU:", cpuUtilization)

		go func() {
			resp, err := http.Get(fmt.Sprintf(
					"http://127.0.0.1:%s/cpu-usage-aggregator:1.0/upload/%d/%f", proxyPort, id, cpuUtilization))
			if err != nil {
				log.Fatal(err)
			}
			resp.Body.Close()
		}()
	}
}


