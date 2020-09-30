package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
		log.Fatal(http.ListenAndServe(":" + listenPort, nil))
	}()

	// Do every 10 seconds
    pushInterval := 10
    timerCh := time.Tick(time.Duration(pushInterval) * time.Second)
    for range timerCh {
		dataset := make([][]float64, 0)
		mux.Lock()
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
		if len(dataset) > 0 {
			log.Println("Sending data")
			datasetBytes, err := json.Marshal(dataset)
			if err != nil {
				log.Fatal(err)
			}
			datasetReader := bytes.NewReader(datasetBytes)
			resp, err := http.Post("http://127.0.0.1:" + proxyPort + "/cpu-usage-predictor:1.0/upload",
					"application/json", datasetReader)
			if err != nil {
				log.Fatal(err)
			}
			resp.Body.Close()
		}
	}
}
