package main

import (
    "bufio"
    "fmt"
    "net/http"
    "os"
	"time"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage:", os.Args[0], "PORT")
        os.Exit(1)
    }

    port := os.Args[1]
	pollInterval := 4
	timerCh := time.Tick(time.Duration(pollInterval) * time.Second)
	// do every 2 seconds
	for range timerCh {
		// make request
        resp, err := http.Get("http://127.0.0.1:" + port + "/hello-world-server")
        if err != nil {
            // panic(err)
            fmt.Println(err.Error())
            continue
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
