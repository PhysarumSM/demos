package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage:", os.Args[0], "<proxy-port> <listen-port>")
		os.Exit(1)
	}
	proxyPort := os.Args[1]
	listenPort := os.Args[2]

	data := make(map[uint64][]float64)
	var mux sync.Mutex

	logFile, err := os.Create("aggregator.log")
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
	http.HandleFunc("/upload/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request from:", r.RemoteAddr, ", for:", r.URL.Path)
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 2 {
			msg := "Error: expected /id/value"
			log.Println(msg)
			fmt.Fprintf(w, "%v", msg)
			return
		}
		id, err := strconv.ParseUint(parts[len(parts)-2], 10, 32)
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "%v", err)
			return
		}
		value, err := strconv.ParseFloat(parts[len(parts)-1], 32)
		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "%v", err)
			return
		}
		mux.Lock()
		data[id] = append(data[id], value)
		mux.Unlock()
		fmt.Fprintf(w, "OK")
	})

	go func(){
		log.Fatalln(http.ListenAndServe(":" + listenPort, nil))
	}()

	success := false
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
		log.Println("Initial prefetch request")
		initResp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s/cpu-usage-predictor:1.0", proxyPort))
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
		log.Println("Initial request for predictor response:", string(body))
		success = true
		break
	}
	if !success {
		log.Fatalln("Could not make inital prefetch request")
	}

	// Do every 15 seconds
    pushInterval := 15
    timerCh := time.Tick(time.Duration(pushInterval) * time.Second)
    for range timerCh {
		dataset := make([][]float64, 0)
		mux.Lock()
		log.Println("Data:", data)
		for id, dataPoints := range data {
			for len(dataPoints) >= 5 {
				row := make([]float64, 5)
				copy(row, dataPoints[:5])
				dataset = append(dataset, row)
				dataPoints = dataPoints[5:]
			}
			data[id] = make([]float64, len(dataPoints))
			copy(data[id], dataPoints)
		}
		mux.Unlock()
		log.Println("Formatted dataset:", dataset)
		if len(dataset) > 0 {
			log.Println("Sending data")
			datasetBytes, err := json.Marshal(dataset)
			if err != nil {
				log.Fatalln(err)
			}
			datasetReader := bytes.NewReader(datasetBytes)
			resp, err := http.Post("http://127.0.0.1:" + proxyPort + "/cpu-usage-predictor:1.0/upload",
					"application/json", datasetReader)
			if err != nil {
				log.Fatalln(err)
			}
			resp.Body.Close()
		}
	}
}
