package main

import (
    "bufio"
    "fmt"
    "net/http"
	"time"
)

func main() {
	pollInterval := 10
	timerCh := time.Tick(time.Duration(pollInterval) * time.Second)
	// do every 10 seconds
	for range timerCh {
		// make request
        resp, err := http.Get("http://127.0.0.1:4201/hello-world-server")
        if err != nil {
            panic(err)
        }

        fmt.Println("Response status:", resp.Status)

		// print out result
        scanner := bufio.NewScanner(resp.Body)
        for i := 0; scanner.Scan() && i < 5; i++ {
            fmt.Println(scanner.Text())
        }

        if err := scanner.Err(); err != nil {
            panic(err)
        }

		resp.Body.Close()
	}
}
