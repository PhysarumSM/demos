package main

import (
	"fmt"
	"io"
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

	logFile, err := os.Create("sensor.log")
	if err != nil {
		log.Fatalln(err)
	}
	defer logFile.Close()
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request from:", r.RemoteAddr)
		fmt.Fprintf(w, "OK")
	})

	go func(){
		log.Fatalln(http.ListenAndServe(":" + listenPort, nil))
	}()

	success := false
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
		log.Println("Initial prefetch request")
		initResp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s/cpu-usage-aggregator:1.0", proxyPort))
		if err != nil {
			log.Println(err)
			continue
		}
		body, err := ioutil.ReadAll(initResp.Body)
		initResp.Body.Close()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("Initial request for aggregator response:", string(body))
		success = true
		break
	}
	if !success {
		log.Fatalln("Could not make inital prefetch request")
	}
	log.Println("Wait 10 seconds before start sending CPU data")
	time.Sleep(time.Second * 10)

	rand.Seed(time.Now().UnixNano())
	id := rand.Uint32()
	for {
		cmd := exec.Command("bash", "-c",
				`top -bn2 -d 0.5 | fgrep 'Cpu(s)' | tail -1 | awk  -F'id,' '{ n=split($1, vals, ","); v=vals[n]; sub("%", "", v); printf "%f", 100 - v }'`)
		output, err := cmd.Output()
		if err != nil {
			log.Fatalln(err)
		}
		cpuUtilization, err := strconv.ParseFloat(string(output), 32)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("CPU:", cpuUtilization)

		go func() {
			resp, err := http.Get(fmt.Sprintf(
					"http://127.0.0.1:%s/cpu-usage-aggregator:1.0/upload/%d/%f", proxyPort, id, cpuUtilization))
			if err != nil {
				log.Fatalln(err)
			}
			resp.Body.Close()
		}()
	}
}


